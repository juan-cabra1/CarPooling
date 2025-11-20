package repository

import (
	"context"
	"testing"
	"time"

	"search-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB creates a test MongoDB connection and returns a clean database
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB (adjust URI as needed for your test environment)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err, "Failed to connect to MongoDB")

	// Use a test database
	db := client.Database("search_api_test")

	// Create 2dsphere index on origin.coordinates for geospatial queries
	_, err = db.Collection("trips").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "origin.coordinates", Value: "2dsphere"}},
	})
	require.NoError(t, err, "Failed to create geospatial index")

	// Cleanup function
	cleanup := func() {
		ctx := context.Background()
		_ = db.Collection("trips").Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return db, cleanup
}

// createTestTrip creates a sample trip for testing
func createTestTrip() *domain.SearchTrip {
	return &domain.SearchTrip{
		TripID:   "trip123",
		DriverID: 1001,
		Driver: domain.Driver{
			ID:         1001,
			Name:       "Juan Pérez",
			Email:      "juan@example.com",
			Rating:     4.8,
			TotalTrips: 25,
		},
		Origin: domain.Location{
			City:     "Buenos Aires",
			Province: "Buenos Aires",
			Address:  "Av. Corrientes 1000",
			Coordinates: domain.NewGeoJSONPoint(-34.6037, -58.3816), // Buenos Aires coordinates
		},
		Destination: domain.Location{
			City:     "La Plata",
			Province: "Buenos Aires",
			Address:  "Calle 7 y 50",
			Coordinates: domain.NewGeoJSONPoint(-34.9214, -57.9544), // La Plata coordinates
		},
		DepartureDatetime:        time.Now().Add(24 * time.Hour),
		EstimatedArrivalDatetime: time.Now().Add(26 * time.Hour),
		PricePerSeat:             500.0,
		TotalSeats:               4,
		AvailableSeats:           3,
		Car: domain.Car{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  2020,
			Color: "Blanco",
			Plate: "ABC123",
		},
		Preferences: domain.Preferences{
			PetsAllowed:    false,
			SmokingAllowed: false,
			MusicAllowed:   true,
		},
		Status:      "published",
		Description: "Viaje cómodo a La Plata",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestTripRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)
	trip := createTestTrip()

	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")
	assert.False(t, trip.ID.IsZero(), "Trip ID should be set after creation")
}

func TestTripRepository_FindByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)
	trip := createTestTrip()

	// Create trip first
	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Find by ID
	found, err := repo.FindByID(context.Background(), trip.ID.Hex())
	require.NoError(t, err, "Failed to find trip")
	require.NotNil(t, found, "Trip should be found")
	assert.Equal(t, trip.TripID, found.TripID, "Trip IDs should match")
	assert.Equal(t, trip.Origin.City, found.Origin.City, "Origin cities should match")
}

func TestTripRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)
	trip := createTestTrip()

	// Create trip
	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Update trip
	trip.Description = "Updated description"
	trip.AvailableSeats = 2
	err = repo.Update(context.Background(), trip)
	require.NoError(t, err, "Failed to update trip")

	// Verify update
	found, err := repo.FindByID(context.Background(), trip.ID.Hex())
	require.NoError(t, err, "Failed to find updated trip")
	assert.Equal(t, "Updated description", found.Description, "Description should be updated")
	assert.Equal(t, 2, found.AvailableSeats, "Available seats should be updated")
}

func TestTripRepository_UpdateStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)
	trip := createTestTrip()

	// Create trip
	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Update status
	err = repo.UpdateStatus(context.Background(), trip.ID.Hex(), "completed")
	require.NoError(t, err, "Failed to update status")

	// Verify status update
	found, err := repo.FindByID(context.Background(), trip.ID.Hex())
	require.NoError(t, err, "Failed to find trip")
	assert.Equal(t, "completed", found.Status, "Status should be updated")
}

func TestTripRepository_UpdateAvailability(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)
	trip := createTestTrip()

	// Create trip
	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Update availability
	err = repo.UpdateAvailability(context.Background(), trip.ID.Hex(), 1)
	require.NoError(t, err, "Failed to update availability")

	// Verify availability update
	found, err := repo.FindByID(context.Background(), trip.ID.Hex())
	require.NoError(t, err, "Failed to find trip")
	assert.Equal(t, 1, found.AvailableSeats, "Available seats should be updated to 1")
}

func TestTripRepository_Search(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)

	// Create multiple trips
	trip1 := createTestTrip()
	trip1.TripID = "trip1"
	trip2 := createTestTrip()
	trip2.TripID = "trip2"
	trip2.Status = "completed"

	err := repo.Create(context.Background(), trip1)
	require.NoError(t, err, "Failed to create trip1")
	err = repo.Create(context.Background(), trip2)
	require.NoError(t, err, "Failed to create trip2")

	// Search for published trips
	filters := map[string]interface{}{
		"status": "published",
	}
	trips, total, err := repo.Search(context.Background(), filters, 1, 10, "popularity", "desc")
	require.NoError(t, err, "Failed to search trips")
	assert.Equal(t, int64(1), total, "Should find 1 published trip")
	assert.Len(t, trips, 1, "Should return 1 trip")
	assert.Equal(t, "trip1", trips[0].TripID, "Should find trip1")
}

func TestTripRepository_SearchByRoute(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)

	// Create trips with different routes
	trip1 := createTestTrip()
	trip1.Origin.City = "Buenos Aires"
	trip1.Destination.City = "La Plata"

	trip2 := createTestTrip()
	trip2.Origin.City = "Buenos Aires"
	trip2.Destination.City = "Córdoba"

	err := repo.Create(context.Background(), trip1)
	require.NoError(t, err, "Failed to create trip1")
	err = repo.Create(context.Background(), trip2)
	require.NoError(t, err, "Failed to create trip2")

	// Search by route
	trips, err := repo.SearchByRoute(context.Background(), "Buenos Aires", "La Plata", nil)
	require.NoError(t, err, "Failed to search by route")
	assert.Len(t, trips, 1, "Should find 1 trip for Buenos Aires -> La Plata")
	assert.Equal(t, "La Plata", trips[0].Destination.City, "Destination should be La Plata")
}

func TestTripRepository_SearchByLocation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)

	// Create trip near Buenos Aires
	trip := createTestTrip()
	trip.Origin.Coordinates = domain.NewGeoJSONPoint(-34.6037, -58.3816) // Buenos Aires

	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Search within 50km radius of Buenos Aires
	lat := -34.6037
	lng := -58.3816
	radiusKm := 50

	trips, err := repo.SearchByLocation(context.Background(), lat, lng, radiusKm, nil)
	require.NoError(t, err, "Failed to search by location")
	assert.GreaterOrEqual(t, len(trips), 1, "Should find at least 1 trip near Buenos Aires")

	// Verify coordinates are in GeoJSON format
	if len(trips) > 0 {
		assert.Equal(t, "Point", trips[0].Origin.Coordinates.Type, "Should be GeoJSON Point type")
		assert.Len(t, trips[0].Origin.Coordinates.Coordinates, 2, "Should have 2 coordinates")
		assert.Equal(t, lng, trips[0].Origin.Coordinates.Lng(), "Longitude should match")
		assert.Equal(t, lat, trips[0].Origin.Coordinates.Lat(), "Latitude should match")
	}
}

func TestTripRepository_SearchByLocation_OutOfRange(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)

	// Create trip in Buenos Aires
	trip := createTestTrip()
	trip.Origin.Coordinates = domain.NewGeoJSONPoint(-34.6037, -58.3816) // Buenos Aires

	err := repo.Create(context.Background(), trip)
	require.NoError(t, err, "Failed to create trip")

	// Search in Córdoba (far from Buenos Aires, > 700km away)
	lat := -31.4201  // Córdoba latitude
	lng := -64.1888  // Córdoba longitude
	radiusKm := 50   // 50km radius

	trips, err := repo.SearchByLocation(context.Background(), lat, lng, radiusKm, nil)
	require.NoError(t, err, "Failed to search by location")
	assert.Len(t, trips, 0, "Should find 0 trips near Córdoba (too far from Buenos Aires)")
}

func TestTripRepository_Pagination(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewTripRepository(db)

	// Create 5 trips
	for i := 1; i <= 5; i++ {
		trip := createTestTrip()
		trip.TripID = primitive.NewObjectID().Hex()
		err := repo.Create(context.Background(), trip)
		require.NoError(t, err, "Failed to create trip")
	}

	// Test pagination: page 1, limit 2
	trips, total, err := repo.Search(context.Background(), map[string]interface{}{}, 1, 2, "", "")
	require.NoError(t, err, "Failed to search with pagination")
	assert.Equal(t, int64(5), total, "Total should be 5")
	assert.Len(t, trips, 2, "Should return 2 trips on page 1")

	// Test pagination: page 2, limit 2
	trips, total, err = repo.Search(context.Background(), map[string]interface{}{}, 2, 2, "", "")
	require.NoError(t, err, "Failed to search with pagination")
	assert.Equal(t, int64(5), total, "Total should still be 5")
	assert.Len(t, trips, 2, "Should return 2 trips on page 2")

	// Test pagination: page 3, limit 2
	trips, total, err = repo.Search(context.Background(), map[string]interface{}{}, 3, 2, "", "")
	require.NoError(t, err, "Failed to search with pagination")
	assert.Equal(t, int64(5), total, "Total should still be 5")
	assert.Len(t, trips, 1, "Should return 1 trip on page 3")
}
