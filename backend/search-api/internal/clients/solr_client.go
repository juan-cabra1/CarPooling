package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"search-api/internal/domain"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// SolrClient provides a simple HTTP client for Apache Solr
type SolrClient struct {
	baseURL string
	core    string
	client  *http.Client
}

// SolrDocument represents a trip document in Solr format
type SolrDocument struct {
	ID string `json:"id"`

	// Driver information
	DriverID         []int64   `json:"driver_id"`
	DriverName       []string  `json:"driver_name"`
	DriverRating     []float64 `json:"driver_rating"`
	DriverTotalTrips []int     `json:"driver_total_trips"`

	// Location information
	OriginCity          []string  `json:"origin_city"`
	OriginProvince      []string  `json:"origin_province"`
	OriginLat           []float64 `json:"origin_lat"`
	OriginLng           []float64 `json:"origin_lng"`
	DestinationCity     []string  `json:"destination_city"`
	DestinationProvince []string  `json:"destination_province"`
	DestinationLat      []float64 `json:"destination_lat"`
	DestinationLng      []float64 `json:"destination_lng"`

	// Trip timing
	DepartureDatetime        []string `json:"departure_datetime"`
	EstimatedArrivalDatetime []string `json:"estimated_arrival_datetime"`

	// Pricing and availability
	PricePerSeat   []float64 `json:"price_per_seat"`
	TotalSeats     []int     `json:"total_seats"`
	AvailableSeats []int     `json:"available_seats"`

	// Vehicle information
	CarBrand []string `json:"car_brand"`
	CarModel []string `json:"car_model"`
	CarYear  []int    `json:"car_year"`
	CarColor []string `json:"car_color"`

	// Preferences
	PetsAllowed    []bool `json:"pets_allowed"`
	SmokingAllowed []bool `json:"smoking_allowed"`
	MusicAllowed   []bool `json:"music_allowed"`

	// Trip details
	Status      []string `json:"status"`
	Description []string `json:"description"`

	// Search-specific fields
	SearchText      []string  `json:"search_text"`
	PopularityScore []float64 `json:"popularity_score"`

	// Timestamps
	CreatedAt []string `json:"created_at"`
	UpdatedAt []string `json:"updated_at"`
}

// SolrResponse represents a search response from Solr
type SolrResponse struct {
	Response struct {
		NumFound int            `json:"numFound"`
		Start    int            `json:"start"`
		Docs     []SolrDocument `json:"docs"`
	} `json:"response"`
}

// SolrUpdateResponse represents an update/delete response from Solr
type SolrUpdateResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
}

// NewSolrClient creates a new Solr client
func NewSolrClient(baseURL, core string) *SolrClient {
	// baseURL format: http://localhost:8983/solr
	fullURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(baseURL, "/"), core)

	log.Info().
		Str("solr_url", fullURL).
		Str("core", core).
		Msg("Initializing Solr client")

	return &SolrClient{
		baseURL: fullURL,
		core:    core,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Index adds or updates a trip document in Solr
func (s *SolrClient) Index(ctx context.Context, trip *domain.SearchTrip) error {
	if trip == nil {
		return fmt.Errorf("trip cannot be nil")
	}

	// Convert SearchTrip to SolrDocument
	doc := s.mapTripToSolrDocument(trip)

	// Marshal document array
	data, err := json.Marshal([]SolrDocument{doc})
	if err != nil {
		log.Error().Err(err).Str("trip_id", trip.TripID).Msg("Failed to marshal Solr document")
		return fmt.Errorf("error marshalling document: %w", err)
	}

	// Send update request
	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("trip_id", trip.TripID).Msg("Failed to execute Solr index request")
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Str("trip_id", trip.TripID).Msg("Solr returned non-OK status")
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		log.Error().Int("solr_status", updateResp.ResponseHeader.Status).Str("trip_id", trip.TripID).Msg("Solr update failed")
		return fmt.Errorf("solr update failed with status %d", updateResp.ResponseHeader.Status)
	}

	log.Debug().Str("trip_id", trip.TripID).Msg("Trip indexed successfully in Solr")
	return nil
}

// Search performs a search query in Solr with filters, using two-phase strategy:
// 1. Try exact match first
// 2. If no results and city filters are present, try partial match
func (s *SolrClient) Search(ctx context.Context, query string, filters map[string]interface{}, page int, limit int, sortBy string, sortOrder string) ([]map[string]interface{}, int, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Check if we have city filters
	hasCityFilters := false
	if originCity, ok := filters["origin_city"].(string); ok && originCity != "" {
		hasCityFilters = true
	}
	if destCity, ok := filters["destination_city"].(string); ok && destCity != "" {
		hasCityFilters = true
	}

	// Phase 1: Try exact match first
	docs, total, err := s.searchWithFilters(ctx, query, filters, page, limit, false, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	// If we have results or no city filters, return immediately
	if total > 0 || !hasCityFilters {
		return docs, total, nil
	}

	// Phase 2: No results with exact match, try partial match on cities
	log.Debug().Msg("No exact match found in Solr, trying partial match on city names")
	docs, total, err = s.searchWithFilters(ctx, query, filters, page, limit, true, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// searchWithFilters performs the actual Solr search with specified match type
func (s *SolrClient) searchWithFilters(ctx context.Context, query string, filters map[string]interface{}, page int, limit int, usePartialMatch bool, sortBy string, sortOrder string) ([]map[string]interface{}, int, error) {
	// Calculate offset
	start := (page - 1) * limit

	// Build query parameters
	params := url.Values{}

	// Main query
	if query == "" {
		params.Set("q", "*:*")
	} else {
		params.Set("q", query)
	}

	params.Set("wt", "json")
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("rows", fmt.Sprintf("%d", limit))

	// Add sorting
	if sortField, sortDir := s.buildSortParam(sortBy, sortOrder); sortField != "" {
		params.Set("sort", fmt.Sprintf("%s %s", sortField, sortDir))
	}

	// Add filters
	if len(filters) > 0 {
		filterQueries := s.buildFilterQueries(filters, usePartialMatch)
		for _, fq := range filterQueries {
			params.Add("fq", fq)
		}
	}

	// Execute search
	searchURL := fmt.Sprintf("%s/select?%s", s.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute Solr search")
		return nil, 0, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Msg("Solr search returned non-OK status")
		return nil, 0, fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var solrResp SolrResponse
	if err := json.NewDecoder(resp.Body).Decode(&solrResp); err != nil {
		log.Error().Err(err).Msg("Failed to decode Solr response")
		return nil, 0, fmt.Errorf("error decoding response: %w", err)
	}

	// Convert SolrDocuments to generic maps
	docs := make([]map[string]interface{}, len(solrResp.Response.Docs))
	for i, doc := range solrResp.Response.Docs {
		docs[i] = s.solrDocumentToMap(doc)
	}

	log.Debug().
		Int("num_found", solrResp.Response.NumFound).
		Int("returned", len(docs)).
		Bool("partial_match", usePartialMatch).
		Str("sort", fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Msg("Solr search completed successfully")

	return docs, solrResp.Response.NumFound, nil
}

// Delete removes a trip document from Solr
func (s *SolrClient) Delete(ctx context.Context, tripID string) error {
	if tripID == "" {
		return fmt.Errorf("tripID cannot be empty")
	}

	data := fmt.Sprintf(`{"delete":{"id":"%s"}}`, tripID)
	url := fmt.Sprintf("%s/update?commit=true", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to execute Solr delete")
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status", resp.StatusCode).Str("trip_id", tripID).Msg("Solr delete returned non-OK status")
		return fmt.Errorf("solr returned status %d", resp.StatusCode)
	}

	var updateResp SolrUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if updateResp.ResponseHeader.Status != 0 {
		log.Error().Int("solr_status", updateResp.ResponseHeader.Status).Str("trip_id", tripID).Msg("Solr delete failed")
		return fmt.Errorf("solr delete failed with status %d", updateResp.ResponseHeader.Status)
	}

	log.Debug().Str("trip_id", tripID).Msg("Trip deleted successfully from Solr")
	return nil
}

// Ping checks if Solr is reachable
func (s *SolrClient) Ping(ctx context.Context) error {
	params := url.Values{}
	params.Set("q", "*:*")
	params.Set("rows", "0")
	params.Set("wt", "json")

	url := fmt.Sprintf("%s/select?%s", s.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating ping request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("core", s.core).Msg("Solr ping failed")
		return fmt.Errorf("solr ping failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr ping returned status %d", resp.StatusCode)
	}

	log.Debug().Str("core", s.core).Msg("Solr ping successful")
	return nil
}

// Helper: mapTripToSolrDocument converts SearchTrip to SolrDocument
// Only indexes non-empty fields to prevent Solr index pollution
func (s *SolrClient) mapTripToSolrDocument(trip *domain.SearchTrip) SolrDocument {
	doc := SolrDocument{
		ID: trip.TripID,
	}

	// Driver information (with validation)
	if trip.DriverID > 0 {
		doc.DriverID = []int64{trip.DriverID}
	}
	if trip.Driver.Name != "" {
		doc.DriverName = []string{trip.Driver.Name}
	}
	if trip.Driver.Rating > 0 {
		doc.DriverRating = []float64{trip.Driver.Rating}
	}
	if trip.Driver.TotalTrips > 0 {
		doc.DriverTotalTrips = []int{trip.Driver.TotalTrips}
	}

	// Origin location (validate non-empty)
	if trip.Origin.City != "" {
		doc.OriginCity = []string{trip.Origin.City}
	}
	if trip.Origin.Province != "" {
		doc.OriginProvince = []string{trip.Origin.Province}
	}
	// Origin coordinates (for display only)
	if len(trip.Origin.Coordinates.Coordinates) == 2 {
		doc.OriginLat = []float64{trip.Origin.Coordinates.Lat()}
		doc.OriginLng = []float64{trip.Origin.Coordinates.Lng()}
	}

	// Destination location (validate non-empty)
	if trip.Destination.City != "" {
		doc.DestinationCity = []string{trip.Destination.City}
	}
	if trip.Destination.Province != "" {
		doc.DestinationProvince = []string{trip.Destination.Province}
	}
	// Destination coordinates (for display only)
	if len(trip.Destination.Coordinates.Coordinates) == 2 {
		doc.DestinationLat = []float64{trip.Destination.Coordinates.Lat()}
		doc.DestinationLng = []float64{trip.Destination.Coordinates.Lng()}
	}

	// Trip timing (check for non-zero times)
	if !trip.DepartureDatetime.IsZero() {
		doc.DepartureDatetime = []string{s.formatSolrDate(trip.DepartureDatetime)}
	}
	if !trip.EstimatedArrivalDatetime.IsZero() {
		doc.EstimatedArrivalDatetime = []string{s.formatSolrDate(trip.EstimatedArrivalDatetime)}
	}

	// Pricing and availability (always include as they have defaults)
	if trip.PricePerSeat > 0 {
		doc.PricePerSeat = []float64{trip.PricePerSeat}
	}
	doc.TotalSeats = []int{trip.TotalSeats}
	doc.AvailableSeats = []int{trip.AvailableSeats}

	// Car information (validate non-empty)
	if trip.Car.Brand != "" {
		doc.CarBrand = []string{trip.Car.Brand}
	}
	if trip.Car.Model != "" {
		doc.CarModel = []string{trip.Car.Model}
	}
	if trip.Car.Year > 0 {
		doc.CarYear = []int{trip.Car.Year}
	}
	if trip.Car.Color != "" {
		doc.CarColor = []string{trip.Car.Color}
	}

	// Preferences (always include as they're boolean)
	doc.PetsAllowed = []bool{trip.Preferences.PetsAllowed}
	doc.SmokingAllowed = []bool{trip.Preferences.SmokingAllowed}
	doc.MusicAllowed = []bool{trip.Preferences.MusicAllowed}

	// Trip details
	if trip.Status != "" {
		doc.Status = []string{trip.Status}
	}
	if trip.Description != "" {
		doc.Description = []string{trip.Description}
	}

	// Search-specific fields
	if trip.SearchText != "" {
		doc.SearchText = []string{trip.SearchText}
	}
	if trip.PopularityScore > 0 {
		doc.PopularityScore = []float64{trip.PopularityScore}
	}

	// Timestamps (check for non-zero)
	if !trip.CreatedAt.IsZero() {
		doc.CreatedAt = []string{s.formatSolrDate(trip.CreatedAt)}
	}
	if !trip.UpdatedAt.IsZero() {
		doc.UpdatedAt = []string{s.formatSolrDate(trip.UpdatedAt)}
	}

	return doc
}

// Helper: solrDocumentToMap converts SolrDocument to generic map
func (s *SolrClient) solrDocumentToMap(doc SolrDocument) map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = doc.ID

	// Extract first element from arrays (Solr uses multi-valued fields)
	if len(doc.DriverID) > 0 {
		m["driver_id"] = doc.DriverID[0]
	}
	if len(doc.DriverName) > 0 {
		m["driver_name"] = doc.DriverName[0]
	}
	if len(doc.DriverRating) > 0 {
		m["driver_rating"] = doc.DriverRating[0]
	}
	if len(doc.DriverTotalTrips) > 0 {
		m["driver_total_trips"] = doc.DriverTotalTrips[0]
	}
	if len(doc.OriginCity) > 0 {
		m["origin_city"] = doc.OriginCity[0]
	}
	if len(doc.OriginProvince) > 0 {
		m["origin_province"] = doc.OriginProvince[0]
	}
	if len(doc.OriginLat) > 0 {
		m["origin_lat"] = doc.OriginLat[0]
	}
	if len(doc.OriginLng) > 0 {
		m["origin_lng"] = doc.OriginLng[0]
	}
	if len(doc.DestinationCity) > 0 {
		m["destination_city"] = doc.DestinationCity[0]
	}
	if len(doc.DestinationProvince) > 0 {
		m["destination_province"] = doc.DestinationProvince[0]
	}
	if len(doc.DestinationLat) > 0 {
		m["destination_lat"] = doc.DestinationLat[0]
	}
	if len(doc.DestinationLng) > 0 {
		m["destination_lng"] = doc.DestinationLng[0]
	}
	if len(doc.DepartureDatetime) > 0 {
		m["departure_datetime"] = doc.DepartureDatetime[0]
	}
	if len(doc.EstimatedArrivalDatetime) > 0 {
		m["estimated_arrival_datetime"] = doc.EstimatedArrivalDatetime[0]
	}
	if len(doc.PricePerSeat) > 0 {
		m["price_per_seat"] = doc.PricePerSeat[0]
	}
	if len(doc.TotalSeats) > 0 {
		m["total_seats"] = doc.TotalSeats[0]
	}
	if len(doc.AvailableSeats) > 0 {
		m["available_seats"] = doc.AvailableSeats[0]
	}
	if len(doc.CarBrand) > 0 {
		m["car_brand"] = doc.CarBrand[0]
	}
	if len(doc.CarModel) > 0 {
		m["car_model"] = doc.CarModel[0]
	}
	if len(doc.CarYear) > 0 {
		m["car_year"] = doc.CarYear[0]
	}
	if len(doc.CarColor) > 0 {
		m["car_color"] = doc.CarColor[0]
	}
	if len(doc.PetsAllowed) > 0 {
		m["pets_allowed"] = doc.PetsAllowed[0]
	}
	if len(doc.SmokingAllowed) > 0 {
		m["smoking_allowed"] = doc.SmokingAllowed[0]
	}
	if len(doc.MusicAllowed) > 0 {
		m["music_allowed"] = doc.MusicAllowed[0]
	}
	if len(doc.Status) > 0 {
		m["status"] = doc.Status[0]
	}
	if len(doc.Description) > 0 {
		m["description"] = doc.Description[0]
	}
	if len(doc.PopularityScore) > 0 {
		m["popularity_score"] = doc.PopularityScore[0]
	}
	if len(doc.CreatedAt) > 0 {
		m["created_at"] = doc.CreatedAt[0]
	}
	if len(doc.UpdatedAt) > 0 {
		m["updated_at"] = doc.UpdatedAt[0]
	}

	return m
}

// Helper: buildFilterQueries converts filter map to Solr filter queries
// usePartialMatch: if true, city fields will use wildcard matching instead of exact match
func (s *SolrClient) buildFilterQueries(filters map[string]interface{}, usePartialMatch bool) []string {
	var fqs []string

	for key, value := range filters {
		switch v := value.(type) {
		case string:
			if v != "" {
				// For city fields, support partial matching with wildcard
				if usePartialMatch && (key == "origin_city" || key == "destination_city") {
					// Use wildcard for prefix search (case-insensitive by default in Solr)
					fqs = append(fqs, fmt.Sprintf(`%s:%s*`, key, strings.ToLower(v)))
				} else {
					// Wrap string values in quotes to handle spaces and special characters
					fqs = append(fqs, fmt.Sprintf(`%s:"%s"`, key, v))
				}
			}
		case int:
			fqs = append(fqs, fmt.Sprintf("%s:%d", key, v))
		case float64:
			fqs = append(fqs, fmt.Sprintf("%s:%f", key, v))
		case bool:
			fqs = append(fqs, fmt.Sprintf("%s:%t", key, v))
		}
	}

	return fqs
}

// Helper: formatSolrDate formats time.Time to ISO 8601 for Solr
func (s *SolrClient) formatSolrDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func (s *SolrClient) buildSortParam(sortBy string, sortOrder string) (string, string) {
	// Default sort direction
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	var solrField string

	switch sortBy {
	case "earliest":
		solrField = "departure_datetime"
		sortOrder = "asc"
	case "cheapest":
		solrField = "price_per_seat"
		sortOrder = "asc"
	case "best_rated":
		solrField = "driver_rating"
		sortOrder = "desc"
	case "popularity":
		solrField = "popularity_score"
		sortOrder = "desc"
	case "":
		// No sorting
		return "", ""
	default:
		// Unknown sort option, default to departure date
		solrField = "departure_datetime"
		sortOrder = "asc"
	}

	return solrField, sortOrder
}
