package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Use test database
	db := client.Database("search_api_test")

	t.Cleanup(func() {
		CleanupTestDB(t, db)
		if err := client.Disconnect(context.Background()); err != nil {
			t.Logf("Failed to disconnect from MongoDB: %v", err)
		}
	})

	return db
}

// CleanupTestDB drops all collections in the test database
func CleanupTestDB(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		t.Logf("Failed to list collections: %v", err)
		return
	}

	for _, collection := range collections {
		if err := db.Collection(collection).Drop(ctx); err != nil {
			t.Logf("Failed to drop collection %s: %v", collection, err)
		}
	}
}

// CreateTestIndexes creates necessary indexes for testing
func CreateTestIndexes(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	// Create 2dsphere index for trips collection
	tripsCollection := db.Collection("trips")
	_, err := tripsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"origin.coordinates": "2dsphere",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create origin geospatial index: %v", err)
	}

	_, err = tripsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"destination.coordinates": "2dsphere",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create destination geospatial index: %v", err)
	}

	// Create unique index for event_id in processed_events collection
	eventsCollection := db.Collection("processed_events")
	_, err = eventsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"event_id": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		t.Fatalf("Failed to create event_id unique index: %v", err)
	}

	// Create compound unique index for popular_routes collection
	routesCollection := db.Collection("popular_routes")
	_, err = routesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"origin_city":      1,
			"destination_city": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		t.Fatalf("Failed to create popular_routes compound index: %v", err)
	}
}

// WaitForMongoDB waits for MongoDB to be ready (useful for CI/CD)
func WaitForMongoDB(uri string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Disconnect(context.Background())

	// Retry ping with backoff
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for MongoDB")
		case <-ticker.C:
			if err := client.Ping(ctx, nil); err == nil {
				return nil
			}
		}
	}
}
