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

// setupPopularRouteTest creates a test MongoDB connection and repository
func setupPopularRouteTest(t *testing.T) (PopularRouteRepository, func()) {
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

	repo := NewPopularRouteRepository(db)

	cleanup := func() {
		ctx := context.Background()
		_ = db.Collection("popular_routes").Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return repo, cleanup
}

func TestPopularRouteRepository_IncrementSearchCount_NewRoute(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	// Increment search count for new route
	err := repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	require.NoError(t, err, "Should create and increment new route")

	// Verify route was created with count = 1
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get top routes")
	require.Len(t, routes, 1, "Should have 1 route")
	assert.Equal(t, "Buenos Aires", routes[0].OriginCity, "Origin should match")
	assert.Equal(t, "La Plata", routes[0].DestinationCity, "Destination should match")
	assert.Equal(t, 1, routes[0].SearchCount, "Search count should be 1")
}

func TestPopularRouteRepository_IncrementSearchCount_ExistingRoute(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	origin := "Buenos Aires"
	destination := "Córdoba"

	// Increment 3 times
	for i := 0; i < 3; i++ {
		err := repo.IncrementSearchCount(context.Background(), origin, destination)
		require.NoError(t, err, "Increment should succeed")
	}

	// Verify count is 3
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get top routes")
	require.Len(t, routes, 1, "Should have 1 route")
	assert.Equal(t, 3, routes[0].SearchCount, "Search count should be 3")
}

func TestPopularRouteRepository_GetTopRoutes_OrderBySearchCount(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	// Create multiple routes with different search counts
	// Buenos Aires -> La Plata: 3 searches
	for i := 0; i < 3; i++ {
		_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "La Plata")
	}

	// Buenos Aires -> Córdoba: 2 searches
	for i := 0; i < 2; i++ {
		_ = repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Córdoba")
	}

	// Rosario -> Mendoza: 1 search
	_ = repo.IncrementSearchCount(context.Background(), "Rosario", "Mendoza")

	// Get all routes
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err, "Should get top routes")
	assert.Len(t, routes, 3, "Should have 3 different routes")

	// Verify ordering by search count (descending)
	assert.Equal(t, 3, routes[0].SearchCount, "Most popular should have 3 searches")
	assert.Equal(t, 2, routes[1].SearchCount, "Second should have 2 searches")
	assert.Equal(t, 1, routes[2].SearchCount, "Least popular should have 1 search")
}

func TestPopularRouteRepository_GetTopRoutes_Limit(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

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
	routes, err := repo.GetTopRoutes(context.Background(), 3)
	require.NoError(t, err, "Should get top routes")
	assert.Len(t, routes, 3, "Should return only 3 routes (limit)")

	// Verify ordering
	assert.GreaterOrEqual(t, routes[0].SearchCount, routes[1].SearchCount, "Should be in descending order")
	assert.GreaterOrEqual(t, routes[1].SearchCount, routes[2].SearchCount, "Should be in descending order")
}

func TestPopularRouteRepository_GetTopRoutes_Empty(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	// Get routes when collection is empty
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err, "Should not return error on empty collection")
	assert.NotNil(t, routes, "Should return non-nil slice")
	assert.Len(t, routes, 0, "Should return empty slice")
}

func TestPopularRouteRepository_IncrementSearchCount_UpdatesLastSearched(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	origin := "Buenos Aires"
	destination := "Córdoba"

	// First increment
	err := repo.IncrementSearchCount(context.Background(), origin, destination)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Second increment
	err = repo.IncrementSearchCount(context.Background(), origin, destination)
	require.NoError(t, err)

	// Verify route exists with count = 2
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, 2, routes[0].SearchCount)
	assert.False(t, routes[0].LastSearched.IsZero(), "LastSearched should be set")
}

func TestPopularRouteRepository_IncrementSearchCount_Concurrent(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	origin := "Buenos Aires"
	destination := "Rosario"
	concurrentCalls := 10

	// Simulate concurrent increments
	errors := make(chan error, concurrentCalls)

	for i := 0; i < concurrentCalls; i++ {
		go func() {
			err := repo.IncrementSearchCount(context.Background(), origin, destination)
			errors <- err
		}()
	}

	// Collect results
	for i := 0; i < concurrentCalls; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	// Verify final count is correct
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, concurrentCalls, routes[0].SearchCount, "All concurrent increments should be counted")
}

func TestPopularRouteRepository_GetTopRoutes_MultipleRoutes(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	// Create multiple routes with specific counts
	testData := []struct {
		origin      string
		destination string
		count       int
	}{
		{"Buenos Aires", "Córdoba", 100},
		{"Buenos Aires", "Rosario", 50},
		{"Mendoza", "San Juan", 75},
		{"La Plata", "Mar del Plata", 25},
	}

	for _, td := range testData {
		for i := 0; i < td.count; i++ {
			err := repo.IncrementSearchCount(context.Background(), td.origin, td.destination)
			require.NoError(t, err)
		}
	}

	// Get all routes
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, routes, 4)

	// Verify they're sorted by count descending
	assert.Equal(t, 100, routes[0].SearchCount, "First should have 100")
	assert.Equal(t, 75, routes[1].SearchCount, "Second should have 75")
	assert.Equal(t, 50, routes[2].SearchCount, "Third should have 50")
	assert.Equal(t, 25, routes[3].SearchCount, "Fourth should have 25")
}

func TestPopularRouteRepository_DifferentRoutes(t *testing.T) {
	repo, cleanup := setupPopularRouteTest(t)
	defer cleanup()

	// Create route A -> B
	err := repo.IncrementSearchCount(context.Background(), "Buenos Aires", "Córdoba")
	require.NoError(t, err)

	// Create route B -> A (should be different)
	err = repo.IncrementSearchCount(context.Background(), "Córdoba", "Buenos Aires")
	require.NoError(t, err)

	// Should have 2 different routes
	routes, err := repo.GetTopRoutes(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, routes, 2, "Reverse routes should be treated as different")
}
