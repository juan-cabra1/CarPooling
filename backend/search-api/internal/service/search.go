package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"search-api/internal/cache"
	"search-api/internal/clients"
	"search-api/internal/domain"
	"search-api/internal/repository"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchService defines the interface for search operations
type SearchService interface {
	// SearchTrips performs a full search with filters using waterfall strategy:
	// Cache → Solr → MongoDB fallback
	SearchTrips(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResponse, error)

	// SearchByLocation performs geospatial search using MongoDB 2dsphere
	// Always uses MongoDB for location-based queries
	SearchByLocation(ctx context.Context, lat, lng float64, radiusKm int, filters map[string]interface{}) (*domain.SearchResponse, error)

	// GetTrip retrieves a single trip with caching
	GetTrip(ctx context.Context, tripID string) (*domain.SearchTrip, error)

	// GetAutocomplete returns city name suggestions
	GetAutocomplete(ctx context.Context, prefix string, limit int) ([]string, error)

	// GetPopularRoutes returns the most searched routes
	GetPopularRoutes(ctx context.Context, limit int) ([]*domain.PopularRoute, error)

	// DenormalizeTrip handles trip.created event: fetches data from external APIs,
	// builds search_text, calculates popularity_score, and stores in MongoDB + Solr
	DenormalizeTrip(ctx context.Context, tripID string) error

	// InvalidateCache removes cached data for a specific trip
	InvalidateCache(ctx context.Context, tripID string) error
}

// searchService implements SearchService
type searchService struct {
	tripRepo         repository.TripRepository
	popularRouteRepo repository.PopularRouteRepository
	cache            cache.Cache
	solrClient       *clients.SolrClient
	tripsClient      clients.TripsClient
	usersClient      clients.UsersClient
	cacheTTL         time.Duration
}

// NewSearchService creates a new SearchService instance
func NewSearchService(
	tripRepo repository.TripRepository,
	popularRouteRepo repository.PopularRouteRepository,
	cache cache.Cache,
	solrClient *clients.SolrClient,
	tripsClient clients.TripsClient,
	usersClient clients.UsersClient,
) SearchService {
	return &searchService{
		tripRepo:         tripRepo,
		popularRouteRepo: popularRouteRepo,
		cache:            cache,
		solrClient:       solrClient,
		tripsClient:      tripsClient,
		usersClient:      usersClient,
		cacheTTL:         15 * time.Minute, // Default cache TTL
	}
}

// SearchTrips implements the waterfall search strategy
func (s *searchService) SearchTrips(ctx context.Context, query *domain.SearchQuery) (*domain.SearchResponse, error) {
	startTime := time.Now()

	// Validate and set defaults
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	query.SetDefaults()

	// Generate cache key
	cacheKey := s.buildSearchCacheKey(query)

	// Step 1: Try Cache
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		log.Debug().
			Str("cache_key", cacheKey).
			Dur("duration_ms", time.Since(startTime)).
			Msg("Cache hit for search query")

		// Track popular routes asynchronously (fire-and-forget)
		go s.trackPopularRoute(context.Background(), query.OriginCity, query.DestinationCity)

		return cached, nil
	}

	var trips []*domain.SearchTrip
	var total int64
	var err error
	var source string

	// Step 2: Try Solr (for non-geospatial queries)
	if !query.IsGeospatial() && s.solrClient != nil {
		trips, total, err = s.searchWithSolr(ctx, query)
		if err == nil {
			source = "solr"
		} else {
			log.Warn().
				Err(err).
				Msg("Solr search failed, falling back to MongoDB")
		}
	}

	// Step 3: Fallback to MongoDB
	if trips == nil {
		trips, total, err = s.searchWithMongoDB(ctx, query)
		if err != nil {
			log.Error().
				Err(err).
				Interface("query", query).
				Msg("MongoDB search failed")
			return nil, fmt.Errorf("search failed: %w", err)
		}
		source = "mongodb"
	}

	// Build response
	response := s.buildSearchResponse(trips, total, query.Page, query.Limit)

	// Cache the result
	if err := s.cacheSearchResult(ctx, cacheKey, response); err != nil {
		log.Warn().Err(err).Msg("Failed to cache search result")
	}

	// Track popular routes asynchronously (fire-and-forget)
	go s.trackPopularRoute(context.Background(), query.OriginCity, query.DestinationCity)

	log.Info().
		Str("source", source).
		Int("results", len(trips)).
		Int64("total", total).
		Dur("duration_ms", time.Since(startTime)).
		Msg("Search completed")

	return response, nil
}

// SearchByLocation performs geospatial search using MongoDB
func (s *searchService) SearchByLocation(ctx context.Context, lat, lng float64, radiusKm int, filters map[string]interface{}) (*domain.SearchResponse, error) {
	startTime := time.Now()

	// Validate coordinates
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}
	if radiusKm <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}

	// Generate cache key for geospatial query
	cacheKey := s.buildLocationCacheKey(lat, lng, radiusKm, filters)

	// Try cache first
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		log.Debug().
			Str("cache_key", cacheKey).
			Dur("duration_ms", time.Since(startTime)).
			Msg("Cache hit for location search")
		return cached, nil
	}

	// MongoDB geospatial search
	trips, err := s.tripRepo.SearchByLocation(ctx, lat, lng, radiusKm, filters)
	if err != nil {
		log.Error().
			Err(err).
			Float64("lat", lat).
			Float64("lng", lng).
			Int("radius_km", radiusKm).
			Msg("Geospatial search failed")
		return nil, fmt.Errorf("geospatial search failed: %w", err)
	}

	// Build response (geospatial doesn't have total count from repo)
	response := &domain.SearchResponse{
		Trips:      trips,
		Total:      int64(len(trips)),
		Page:       1,
		Limit:      len(trips),
		TotalPages: 1,
	}

	// Cache the result
	if err := s.cacheSearchResult(ctx, cacheKey, response); err != nil {
		log.Warn().Err(err).Msg("Failed to cache location search result")
	}

	log.Info().
		Int("results", len(trips)).
		Dur("duration_ms", time.Since(startTime)).
		Msg("Geospatial search completed")

	return response, nil
}

// GetTrip retrieves a single trip with caching
func (s *searchService) GetTrip(ctx context.Context, tripID string) (*domain.SearchTrip, error) {
	startTime := time.Now()

	// Try cache first
	cacheKey := s.buildTripCacheKey(tripID)
	if cached, err := s.getTripFromCache(ctx, cacheKey); err == nil && cached != nil {
		log.Debug().
			Str("trip_id", tripID).
			Dur("duration_ms", time.Since(startTime)).
			Msg("Cache hit for trip")
		return cached, nil
	}

	// Fetch from MongoDB
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to fetch trip")
		return nil, fmt.Errorf("failed to fetch trip: %w", err)
	}

	if trip == nil {
		return nil, fmt.Errorf("trip not found")
	}

	// Cache the trip
	if err := s.cacheTripData(ctx, cacheKey, trip); err != nil {
		log.Warn().Err(err).Msg("Failed to cache trip")
	}

	log.Debug().
		Str("trip_id", tripID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("Trip fetched successfully")

	return trip, nil
}

// GetAutocomplete returns city name suggestions
func (s *searchService) GetAutocomplete(ctx context.Context, prefix string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Normalize prefix
	prefix = strings.TrimSpace(strings.ToLower(prefix))
	if len(prefix) < 2 {
		return []string{}, nil
	}

	// TODO: Implement autocomplete using Solr facets or a dedicated collection
	// For now, return empty array
	// In production, you would:
	// 1. Use Solr facets on origin.city and destination.city fields
	// 2. Or maintain a separate "cities" collection in MongoDB with pre-computed popular cities
	log.Warn().Msg("GetAutocomplete not yet fully implemented")

	return []string{}, nil
}

// GetPopularRoutes returns the most searched routes
func (s *searchService) GetPopularRoutes(ctx context.Context, limit int) ([]*domain.PopularRoute, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	routes, err := s.popularRouteRepo.GetTopRoutes(ctx, limit)
	if err != nil {
		log.Error().
			Err(err).
			Int("limit", limit).
			Msg("Failed to fetch popular routes")
		return nil, fmt.Errorf("failed to fetch popular routes: %w", err)
	}

	// Convert to pointer slice
	routePtrs := make([]*domain.PopularRoute, len(routes))
	for i := range routes {
		routePtrs[i] = &routes[i]
	}

	return routePtrs, nil
}

// DenormalizeTrip handles trip.created event
func (s *searchService) DenormalizeTrip(ctx context.Context, tripID string) error {
	startTime := time.Now()

	log.Info().
		Str("trip_id", tripID).
		Msg("Starting trip denormalization")

	// Step 1: Fetch full trip from trips-api
	trip, err := s.tripsClient.GetTrip(ctx, tripID)
	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to fetch trip from trips-api")
		return fmt.Errorf("failed to fetch trip: %w", err)
	}

	// Step 2: Fetch driver info from users-api
	driver, err := s.usersClient.GetUser(ctx, trip.DriverID)
	if err != nil {
		log.Error().
			Err(err).
			Int64("driver_id", trip.DriverID).
			Msg("Failed to fetch driver from users-api")
		return fmt.Errorf("failed to fetch driver: %w", err)
	}

	// Step 3: Build SearchTrip with denormalized data
	searchTrip := &domain.SearchTrip{
		ID:                       primitive.NewObjectID(),
		TripID:                   trip.ID.Hex(),
		DriverID:                 trip.DriverID,
		Driver:                   s.mapUserToDriver(driver),
		Origin:                   trip.Origin.ToLocation(),
		Destination:              trip.Destination.ToLocation(),
		DepartureDatetime:        trip.DepartureDatetime,
		EstimatedArrivalDatetime: trip.EstimatedArrivalDatetime,
		PricePerSeat:             trip.PricePerSeat,
		TotalSeats:               trip.TotalSeats,
		AvailableSeats:           trip.AvailableSeats,
		Car:                      trip.Car,
		Preferences:              trip.Preferences,
		Status:                   trip.Status,
		Description:              trip.Description,
		CreatedAt:                trip.CreatedAt,
		UpdatedAt:                trip.UpdatedAt,
	}

	// Step 4: Build search_text from city names + description
	searchTrip.SearchText = BuildSearchText(searchTrip)

	// Step 5: Calculate popularity_score
	searchTrip.PopularityScore = CalculatePopularityScore(searchTrip)

	// Step 6: Store in MongoDB
	if err := s.tripRepo.Create(ctx, searchTrip); err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to store trip in MongoDB")
		return fmt.Errorf("failed to store trip in MongoDB: %w", err)
	}

	// Step 7: Index in Solr (graceful degradation if Solr is down)
	if s.solrClient != nil {
		if err := s.solrClient.Index(ctx, searchTrip); err != nil {
			// Log warning but don't fail the operation
			log.Warn().
				Err(err).
				Str("trip_id", tripID).
				Msg("Failed to index trip in Solr (non-critical)")
		}
	}

	log.Info().
		Str("trip_id", tripID).
		Dur("duration_ms", time.Since(startTime)).
		Msg("Trip denormalization completed successfully")

	return nil
}

// InvalidateCache removes cached data for a specific trip
func (s *searchService) InvalidateCache(ctx context.Context, tripID string) error {
	cacheKey := s.buildTripCacheKey(tripID)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		log.Warn().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to invalidate trip cache")
		return err
	}

	log.Debug().
		Str("trip_id", tripID).
		Msg("Trip cache invalidated")

	return nil
}

// ===== Helper Methods =====

// buildSearchCacheKey generates a deterministic cache key for a search query
func (s *searchService) buildSearchCacheKey(query *domain.SearchQuery) string {
	return fmt.Sprintf("search:query:%s", query.Hash())
}

// buildLocationCacheKey generates cache key for geospatial queries
func (s *searchService) buildLocationCacheKey(lat, lng float64, radiusKm int, filters map[string]interface{}) string {
	// Simple hash for location queries
	filtersJSON, _ := json.Marshal(filters)
	return fmt.Sprintf("search:location:%.6f:%.6f:%d:%s", lat, lng, radiusKm, string(filtersJSON))
}

// buildTripCacheKey generates cache key for a single trip
func (s *searchService) buildTripCacheKey(tripID string) string {
	return fmt.Sprintf("trip:%s", tripID)
}

// getFromCache retrieves search results from cache
func (s *searchService) getFromCache(ctx context.Context, cacheKey string) (*domain.SearchResponse, error) {
	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var response domain.SearchResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// getTripFromCache retrieves a single trip from cache
func (s *searchService) getTripFromCache(ctx context.Context, cacheKey string) (*domain.SearchTrip, error) {
	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var trip domain.SearchTrip
	if err := json.Unmarshal([]byte(data), &trip); err != nil {
		return nil, err
	}

	return &trip, nil
}

// cacheSearchResult stores search results in cache
func (s *searchService) cacheSearchResult(ctx context.Context, cacheKey string, response *domain.SearchResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, cacheKey, string(data), s.cacheTTL)
}

// cacheTripData stores a single trip in cache
func (s *searchService) cacheTripData(ctx context.Context, cacheKey string, trip *domain.SearchTrip) error {
	data, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, cacheKey, string(data), s.cacheTTL)
}

// searchWithSolr performs search using Apache Solr
func (s *searchService) searchWithSolr(ctx context.Context, query *domain.SearchQuery) ([]*domain.SearchTrip, int64, error) {
	// Build simple query string
	queryStr := "*:*"
	if query.SearchText != "" {
		queryStr = fmt.Sprintf("search_text:%s", query.SearchText)
	}

	// Build filters map
	filters := make(map[string]interface{})
	filters["status"] = "published"

	if query.OriginCity != "" {
		filters["origin_city"] = query.OriginCity
	}
	if query.DestinationCity != "" {
		filters["destination_city"] = query.DestinationCity
	}
	if query.MinSeats > 0 {
		filters["available_seats"] = fmt.Sprintf("[%d TO *]", query.MinSeats)
	}
	if query.MaxPrice > 0 {
		filters["price_per_seat"] = fmt.Sprintf("[* TO %f]", query.MaxPrice)
	}

	// Execute Solr search with new simple client
	docs, total, err := s.solrClient.Search(ctx, queryStr, filters, query.Page, query.Limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert Solr documents to SearchTrip entities
	// Extract trip IDs and fetch full data from MongoDB for complete details
	tripIDs := make([]string, 0, len(docs))
	for _, doc := range docs {
		if id, ok := doc["id"].(string); ok {
			tripIDs = append(tripIDs, id)
		}
	}

	// If no results from Solr, return empty
	if len(tripIDs) == 0 {
		return []*domain.SearchTrip{}, 0, nil
	}

	// Fetch trips from MongoDB one by one (TODO: implement batch fetch)
	trips := make([]*domain.SearchTrip, 0, len(tripIDs))
	for _, tripID := range tripIDs {
		trip, err := s.tripRepo.FindByTripID(ctx, tripID)
		if err != nil {
			log.Warn().Err(err).Str("trip_id", tripID).Msg("Failed to fetch trip from MongoDB, skipping")
			continue
		}
		trips = append(trips, trip)
	}

	return trips, int64(total), nil
}

// searchWithMongoDB performs search using MongoDB
func (s *searchService) searchWithMongoDB(ctx context.Context, query *domain.SearchQuery) ([]*domain.SearchTrip, int64, error) {
	// Build MongoDB filters
	filters := s.buildMongoFilters(query)

	// Execute MongoDB search
	trips, total, err := s.tripRepo.Search(ctx, filters, query.Page, query.Limit)
	if err != nil {
		return nil, 0, err
	}

	return trips, total, nil
}

// buildMongoFilters converts SearchQuery to MongoDB filters
func (s *searchService) buildMongoFilters(query *domain.SearchQuery) map[string]interface{} {
	filters := make(map[string]interface{})

	// Always filter by published status and available seats
	filters["status"] = "published"
	filters["available_seats"] = map[string]interface{}{"$gte": 1}

	// City filters
	if query.OriginCity != "" {
		filters["origin.city"] = query.OriginCity
	}
	if query.DestinationCity != "" {
		filters["destination.city"] = query.DestinationCity
	}

	// Seats filter
	if query.MinSeats > 0 {
		filters["available_seats"] = map[string]interface{}{"$gte": query.MinSeats}
	}

	// Price filter
	if query.MaxPrice > 0 {
		filters["price_per_seat"] = map[string]interface{}{"$lte": query.MaxPrice}
	}

	// Preference filters
	if query.PetsAllowed != nil {
		filters["preferences.pets_allowed"] = *query.PetsAllowed
	}
	if query.SmokingAllowed != nil {
		filters["preferences.smoking_allowed"] = *query.SmokingAllowed
	}
	if query.MusicAllowed != nil {
		filters["preferences.music_allowed"] = *query.MusicAllowed
	}

	// Driver rating filter
	if query.MinDriverRating > 0 {
		filters["driver.rating"] = map[string]interface{}{"$gte": query.MinDriverRating}
	}

	// Date range filter
	if !query.DateFrom.IsZero() || !query.DateTo.IsZero() {
		dateFilter := make(map[string]interface{})
		if !query.DateFrom.IsZero() {
			dateFilter["$gte"] = query.DateFrom
		}
		if !query.DateTo.IsZero() {
			dateFilter["$lte"] = query.DateTo
		}
		filters["departure_datetime"] = dateFilter
	}

	return filters
}

// buildSearchResponse builds a SearchResponse from results
func (s *searchService) buildSearchResponse(trips []*domain.SearchTrip, total int64, page, limit int) *domain.SearchResponse {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &domain.SearchResponse{
		Trips:      trips,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// trackPopularRoute increments search count for a route asynchronously
func (s *searchService) trackPopularRoute(ctx context.Context, originCity, destinationCity string) {
	if originCity == "" || destinationCity == "" {
		return
	}

	if err := s.popularRouteRepo.IncrementSearchCount(ctx, originCity, destinationCity); err != nil {
		log.Warn().
			Err(err).
			Str("origin", originCity).
			Str("destination", destinationCity).
			Msg("Failed to track popular route")
	}
}

// mapUserToDriver converts User DTO to Driver domain model
func (s *searchService) mapUserToDriver(user *domain.User) domain.Driver {
	return domain.Driver{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		PhotoURL:   user.PhotoURL,
		Rating:     user.AverageRatingAsDriver,
		TotalTrips: user.TotalTripsAsDriver,
	}
}

// ===== Exported Helper Functions =====

// BuildSearchText concatenates city names and description for full-text search
func BuildSearchText(trip *domain.SearchTrip) string {
	parts := []string{
		trip.Origin.City,
		trip.Origin.Province,
		trip.Destination.City,
		trip.Destination.Province,
		trip.Description,
		trip.Driver.Name,
	}

	// Filter empty strings and join with spaces
	var nonEmpty []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			nonEmpty = append(nonEmpty, trimmed)
		}
	}

	return strings.Join(nonEmpty, " ")
}

// CalculatePopularityScore calculates a popularity score based on:
// - Number of bookings (approximated by occupied seats)
// - Driver rating
// - Trip recency (newer trips get slight boost)
func CalculatePopularityScore(trip *domain.SearchTrip) float64 {
	score := 0.0

	// Component 1: Booking popularity (occupied seats)
	// Weight: 40%
	occupiedSeats := trip.TotalSeats - trip.AvailableSeats
	if trip.TotalSeats > 0 {
		occupancyRate := float64(occupiedSeats) / float64(trip.TotalSeats)
		score += occupancyRate * 40.0
	}

	// Component 2: Driver rating (0-5 scale)
	// Weight: 40%
	driverScore := (trip.Driver.Rating / 5.0) * 40.0
	score += driverScore

	// Component 3: Driver experience (total trips)
	// Weight: 10%
	// Logarithmic scale: ln(trips + 1) capped at 100 trips = max score
	experienceScore := 0.0
	if trip.Driver.TotalTrips > 0 {
		// ln(101) ≈ 4.615, so we normalize by dividing by 4.615
		experienceScore = (float64(trip.Driver.TotalTrips) / 100.0) * 10.0
		if experienceScore > 10.0 {
			experienceScore = 10.0
		}
	}
	score += experienceScore

	// Component 4: Recency bonus (trips created recently)
	// Weight: 10%
	// Trips created in last 7 days get full points, older trips decay linearly
	daysSinceCreation := time.Since(trip.CreatedAt).Hours() / 24
	recencyScore := 10.0
	if daysSinceCreation > 7 {
		recencyScore = 10.0 * (1.0 - (daysSinceCreation-7.0)/30.0)
		if recencyScore < 0 {
			recencyScore = 0
		}
	}
	score += recencyScore

	// Ensure score is within 0-100 range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
