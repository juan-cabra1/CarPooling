package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupPopularRouteTestDB creates a test MongoDB connection for popular route repository tests
func setupPopularRouteTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err, "Failed to connect to MongoDB")

	db := client.Database("search_api_test_routes")

	// Create unique compound index on (origin_city, destination_city)
	_, err = db.Collection("popular_routes").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "origin_city", Value: 1},
			{Key: "destination_city", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	require.NoError(t, err, "Failed to create unique index on popular routes")

	cleanup := func() {
		ctx := context.Background()
		_ = db.Collection("popular_routes").Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestPopularRouteRepository_IncrementSearchCount_NewRoute(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Increment search count for new route
	err := repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	require.NoError(t, err, "Should create and increment new route")

	// Verify route was created with count = 1
	routes, err := repo.GetPopularRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get popular routes")
	require.Len(t, routes, 1, "Should have 1 route")
	assert.Equal(t, "Buenos Aires", routes[0].OriginCity, "Origin should match")
	assert.Equal(t, "La Plata", routes[0].DestinationCity, "Destination should match")
	assert.Equal(t, 1, routes[0].SearchCount, "Search count should be 1")
}

func TestPopularRouteRepository_IncrementSearchCount_ExistingRoute(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	origin := "Buenos Aires"
	destination := "Córdoba"

	// Increment 3 times
	err := repo.IncrementSearchCount(context.Background(), origin, destination)
	require.NoError(t, err, "First increment should succeed")

	err = repo.IncrementSearchCount(context.Background(), origin, destination)
	require.NoError(t, err, "Second increment should succeed")

	err = repo.IncrementSearchCount(context.Background(), origin, destination)
	require.NoError(t, err, "Third increment should succeed")

	// Verify count is 3
	routes, err := repo.GetPopularRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get popular routes")
	require.Len(t, routes, 1, "Should have 1 route")
	assert.Equal(t, 3, routes[0].SearchCount, "Search count should be 3")
}

func TestPopularRouteRepository_GetPopularRoutes_OrderBySearchCount(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create multiple routes with different search counts
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")

	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Córdoba")
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Córdoba")

	_ = repo.IncrementSearchCount(context.Background(), "Rosario", "Mendoza")

	// Get all routes
	routes, err := repo.GetPopularRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get popular routes")
	assert.Len(t, routes, 3, "Should have 3 different routes")

	// Verify ordering by search count (descending)
	assert.Equal(t, 3, routes[0].SearchCount, "Most popular should have 3 searches")
	assert.Equal(t, 2, routes[1].SearchCount, "Second should have 2 searches")
	assert.Equal(t, 1, routes[2].SearchCount, "Least popular should have 1 search")
}

func TestPopularRouteRepository_GetPopularRoutes_Limit(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create 5 routes with different search counts
	cities := [][]string{
		{"Buenos Aires", "La Plata"},
		{"Buenos Aires", "Córdoba"},
		{"Rosario", "Mendoza"},
		{"Mar del Plata", "Buenos Aires"},
		{"La Plata", "Rosario"},
	}

	for i, route := range cities {
		// Different search counts: 5, 4, 3, 2, 1
		for j := 0; j <= (4 - i); j++ {
			_ = repo.IncrementSearchCount(context.Background(), route[0], route[1])
		}
	}

	// Get top 3 routes
	routes, err := repo.GetPopularRoutes(context.Background(), 3)
	require.NoError(t, err, "Should get popular routes")
	assert.Len(t, routes, 3, "Should return only 3 routes (limit)")

	// Verify ordering
	assert.GreaterOrEqual(t, routes[0].SearchCount, routes[1].SearchCount, "Should be in descending order")
	assert.GreaterOrEqual(t, routes[1].SearchCount, routes[2].SearchCount, "Should be in descending order")
}

func TestPopularRouteRepository_GetPopularRoutes_Empty(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Get routes when collection is empty
	routes, err := repo.GetPopularRoutes(context.Background(), 10)
	require.NoError(t, err, "Should not return error on empty collection")
	assert.NotNil(t, routes, "Should return non-nil slice")
	assert.Len(t, routes, 0, "Should return empty slice")
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_OriginMatch(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create routes with different cities
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	_ = repo.IncrementSearchCount(context.Background(), "Barcelona", "Madrid")
	_ = repo.IncrementSearchCount(context.Background(), "Córdoba", "Rosario")

	// Search for cities starting with "B" (case-insensitive)
	cities, err := repo.GetAutocompleteSuggestions(context.Background(), "B", 10)
	require.NoError(t, err, "Should get autocomplete suggestions")
	assert.GreaterOrEqual(t, len(cities), 2, "Should find at least 2 cities starting with B")
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_CaseInsensitive(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Mendoza")

	// Search with lowercase
	cities, err := repo.GetAutocompleteSuggestions(context.Background(), "buenos", 10)
	require.NoError(t, err, "Should get suggestions")
	assert.GreaterOrEqual(t, len(cities), 1, "Should find matches case-insensitively")

	// Search with uppercase
	cities, err = repo.GetAutocompleteSuggestions(context.Background(), "BUENOS", 10)
	require.NoError(t, err, "Should get suggestions")
	assert.GreaterOrEqual(t, len(cities), 1, "Should find matches case-insensitively")
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_Unique(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create routes where "Buenos Aires" appears multiple times
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Córdoba")
	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Rosario")
	_ = repo.IncrementSearchCount(context.Background(), "La Plata", "Buenos Aires")

	// Search for "Buenos"
	cities, err := repo.GetAutocompleteSuggestions(context.Background(), "Buenos", 10)
	require.NoError(t, err, "Should get suggestions")

	// Verify uniqueness
	uniqueCities := make(map[string]bool)
	for _, city := range cities {
		assert.False(t, uniqueCities[city], "City should appear only once: "+city)
		uniqueCities[city] = true
	}
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_Limit(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create many routes starting with "C"
	cities := []string{"Córdoba", "Corrientes", "Catamarca", "Comodoro Rivadavia", "Concordia"}
	for _, city := range cities {
		_ = repo.IncrementSearchCount(context.Background(), city, "Buenos Aires")
	}

	// Get only 3 suggestions
	suggestions, err := repo.GetAutocompleteSuggestions(context.Background(), "C", 3)
	require.NoError(t, err, "Should get suggestions")
	assert.LessOrEqual(t, len(suggestions), 3, "Should respect limit of 3")
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_NoMatch(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")

	// Search for non-existent prefix
	cities, err := repo.GetAutocompleteSuggestions(context.Background(), "XYZ", 10)
	require.NoError(t, err, "Should not return error")
	assert.NotNil(t, cities, "Should return non-nil slice")
	assert.Len(t, cities, 0, "Should return empty slice for no matches")
}

func TestPopularRouteRepository_GetAutocompleteSuggestions_AlphabeticalOrder(t *testing.T) {
	db, cleanup := setupPopularRouteTestDB(t)
	defer cleanup()

	repo := NewPopularRouteRepository(db)

	// Create routes in non-alphabetical order
	_ = repo.IncrementSearchCount(context.Background(), "Mendoza", "Rosario")
	_ = repo.IncrementSearchCount(context.Background(), "Mar del Plata", "Buenos Aires")

	// Get suggestions for "M"
	cities, err := repo.GetAutocompleteSuggestions(context.Background(), "M", 10)
	require.NoError(t, err, "Should get suggestions")

	// Verify alphabetical ordering
	for i := 0; i < len(cities)-1; i++ {
		assert.LessOrEqual(t, cities[i], cities[i+1], "Cities should be in alphabetical order")
	}
}
