package repository

import (
	"context"
	"fmt"
	"time"

	"search-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TripRepository handles trip search operations
type TripRepository interface {
	Create(ctx context.Context, trip *domain.SearchTrip) error
	FindByID(ctx context.Context, id string) (*domain.SearchTrip, error)
	FindByTripID(ctx context.Context, tripID string) (*domain.SearchTrip, error)
	Update(ctx context.Context, trip *domain.SearchTrip) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateStatusByTripID(ctx context.Context, tripID string, status string) error
	UpdateAvailability(ctx context.Context, id string, availableSeats int) error
	UpdateAvailabilityByTripID(ctx context.Context, tripID string, availableSeats int, reservedSeats int, status string) error
	DeleteByTripID(ctx context.Context, tripID string) error
	Search(ctx context.Context, filters map[string]interface{}, page, limit int, sortBy string, sortOrder string) ([]*domain.SearchTrip, int64, error)
	SearchByLocation(ctx context.Context, lat, lng float64, radiusKm int, additionalFilters map[string]interface{}) ([]*domain.SearchTrip, error)
	SearchByRoute(ctx context.Context, originCity, destinationCity string, filters map[string]interface{}) ([]*domain.SearchTrip, error)
}

type tripRepository struct {
	collection *mongo.Collection
}

// NewTripRepository creates a new trip repository instance
func NewTripRepository(db *mongo.Database) TripRepository {
	return &tripRepository{
		collection: db.Collection("trips"),
	}
}

// Create inserts a new trip into the database
func (r *tripRepository) Create(ctx context.Context, trip *domain.SearchTrip) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if trip.ID.IsZero() {
		trip.ID = primitive.NewObjectID()
	}

	if trip.CreatedAt.IsZero() {
		trip.CreatedAt = time.Now()
	}
	trip.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

// FindByID retrieves a trip by its ID
func (r *tripRepository) FindByID(ctx context.Context, id string) (*domain.SearchTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid trip ID: %w", err)
	}

	var trip domain.SearchTrip
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&trip)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find trip: %w", err)
	}

	return &trip, nil
}

// Update updates an existing trip
func (r *tripRepository) Update(ctx context.Context, trip *domain.SearchTrip) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	trip.UpdatedAt = time.Now()

	filter := bson.M{"_id": trip.ID}
	update := bson.M{"$set": trip}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("trip not found")
	}

	return nil
}

// UpdateStatus updates only the status of a trip
func (r *tripRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid trip ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("trip not found")
	}

	return nil
}

// UpdateAvailability updates the available seats of a trip
func (r *tripRepository) UpdateAvailability(ctx context.Context, id string, availableSeats int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid trip ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"available_seats": availableSeats,
			"updated_at":      time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip availability: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("trip not found")
	}

	return nil
}

// Search performs a generic search with filters and pagination
func (r *tripRepository) Search(ctx context.Context, filters map[string]interface{}, page, limit int, sortBy string, sortOrder string) ([]*domain.SearchTrip, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Calculate skip for pagination
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}

	// Check if this is a geospatial query by looking for $near operator
	// MongoDB's CountDocuments doesn't support $near, so we need different approach
	var total int64
	var needsPostCount bool

	// Check origin coordinates for $near
	if originCoords, ok := filter["origin.coordinates"].(bson.M); ok {
		if _, hasNear := originCoords["$near"]; hasNear {
			needsPostCount = true
		}
	}

	// Check destination coordinates for $near
	if !needsPostCount {
		if destCoords, ok := filter["destination.coordinates"].(bson.M); ok {
			if _, hasNear := destCoords["$near"]; hasNear {
				needsPostCount = true
			}
		}
	}

	// For non-geospatial queries, count normally
	if !needsPostCount {
		var err error
		total, err = r.collection.CountDocuments(ctx, filter)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to count trips: %w", err)
		}
	}

	// Find documents with pagination
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Only apply sorting if sortBy is provided
	// For geospatial queries with $near, MongoDB automatically sorts by distance
	if sortBy != "" {
		sortBson := r.buildSortOptions(sortBy, sortOrder)
		findOptions.SetSort(sortBson)
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search trips: %w", err)
	}
	defer cursor.Close(ctx)

	var trips []*domain.SearchTrip
	if err = cursor.All(ctx, &trips); err != nil {
		return nil, 0, fmt.Errorf("failed to decode trips: %w", err)
	}

	// For geospatial queries, count the results after fetching
	if needsPostCount {
		total = int64(len(trips))
	}

	// Return empty slice instead of nil
	if trips == nil {
		trips = []*domain.SearchTrip{}
	}

	return trips, total, nil
}

// buildSortOptions converts sortBy and sortOrder to MongoDB sort bson.D
// Supports both flexible format (sortBy + sortOrder) and backward compatible shortcuts
func (r *tripRepository) buildSortOptions(sortBy string, sortOrder string) bson.D {
	// Determine sort direction: 1 for ascending, -1 for descending
	direction := 1 // Default to ascending
	if sortOrder == "desc" {
		direction = -1
	}

	// Handle backward compatibility shortcuts (ignore sortOrder for these)
	switch sortBy {
	case "earliest":
		return bson.D{{Key: "departure_datetime", Value: 1}}
	case "cheapest":
		return bson.D{{Key: "price_per_seat", Value: 1}}
	case "best_rated":
		return bson.D{{Key: "driver.rating", Value: -1}}
	}

	// Handle new flexible format (respects sortOrder parameter)
	var field string
	switch sortBy {
	case "price":
		field = "price_per_seat"
	case "departure_time":
		field = "departure_datetime"
	case "rating":
		field = "driver.rating"
	case "popularity":
		field = "popularity_score"
	default:
		// Default to departure_datetime if sortBy is invalid or empty
		field = "departure_datetime"
	}

	return bson.D{{Key: field, Value: direction}}
}

// SearchByLocation performs geospatial search using MongoDB's 2dsphere index
// Uses $near operator to find trips within a radius from the given coordinates
// Results are automatically sorted by distance from the point
// Reference: https://www.mongodb.com/docs/manual/reference/operator/query/near/
func (r *tripRepository) SearchByLocation(ctx context.Context, lat, lng float64, radiusKm int, additionalFilters map[string]interface{}) ([]*domain.SearchTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build geospatial filter with $near operator
	// IMPORTANT: MongoDB GeoJSON uses [longitude, latitude] order (lng first!)
	// The $near operator returns documents from nearest to farthest
	filter := bson.M{
		"origin.coordinates": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat}, // [longitude, latitude]
				},
				"$maxDistance": radiusKm * 1000, // Convert km to meters
			},
		},
		"status":          "published",
		"available_seats": bson.M{"$gte": 1},
	}

	// Add additional filters (e.g., departure_datetime, price_per_seat, preferences)
	for key, value := range additionalFilters {
		filter[key] = value
	}

	// Note: $near already sorts by distance, no additional sort needed
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search trips by location: %w", err)
	}
	defer cursor.Close(ctx)

	var trips []*domain.SearchTrip
	if err = cursor.All(ctx, &trips); err != nil {
		return nil, fmt.Errorf("failed to decode trips: %w", err)
	}

	// Return empty slice instead of nil
	if trips == nil {
		trips = []*domain.SearchTrip{}
	}

	return trips, nil
}

// SearchByRoute searches trips by origin and destination cities
func (r *tripRepository) SearchByRoute(ctx context.Context, originCity, destinationCity string, filters map[string]interface{}) ([]*domain.SearchTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build filter for city-based search
	filter := bson.M{
		"origin.city":      originCity,
		"destination.city": destinationCity,
		"status":           "published",
		"available_seats":  bson.M{"$gte": 1},
	}

	// Add additional filters
	for key, value := range filters {
		filter[key] = value
	}

	// Sort by departure time
	findOptions := options.Find().
		SetSort(bson.D{{Key: "departure_datetime", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to search trips by route: %w", err)
	}
	defer cursor.Close(ctx)

	var trips []*domain.SearchTrip
	if err = cursor.All(ctx, &trips); err != nil {
		return nil, fmt.Errorf("failed to decode trips: %w", err)
	}

	// Return empty slice instead of nil
	if trips == nil {
		trips = []*domain.SearchTrip{}
	}

	return trips, nil
}

// FindByTripID retrieves a trip by its trip_id field (not MongoDB _id)
func (r *tripRepository) FindByTripID(ctx context.Context, tripID string) (*domain.SearchTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var trip domain.SearchTrip
	err := r.collection.FindOne(ctx, bson.M{"trip_id": tripID}).Decode(&trip)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrSearchTripNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find trip by trip_id: %w", err)
	}

	return &trip, nil
}

// UpdateStatusByTripID updates only the status of a trip using trip_id field
func (r *tripRepository) UpdateStatusByTripID(ctx context.Context, tripID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"trip_id": tripID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip status: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrSearchTripNotFound
	}

	return nil
}

// UpdateAvailabilityByTripID updates availability, reserved seats, and status using trip_id field
func (r *tripRepository) UpdateAvailabilityByTripID(ctx context.Context, tripID string, availableSeats int, reservedSeats int, status string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"trip_id": tripID}
	update := bson.M{
		"$set": bson.M{
			"available_seats": availableSeats,
			"reserved_seats":  reservedSeats,
			"status":          status,
			"updated_at":      time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update trip availability: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrSearchTripNotFound
	}

	return nil
}

// DeleteByTripID deletes a trip from the search index by trip_id
func (r *tripRepository) DeleteByTripID(ctx context.Context, tripID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"trip_id": tripID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrSearchTripNotFound
	}

	return nil
}
