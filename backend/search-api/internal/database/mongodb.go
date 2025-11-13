package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongoDB establishes connection to MongoDB and returns the database instance
func ConnectMongoDB(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("✅ MongoDB connected successfully to database: %s", dbName)
	return client.Database(dbName), nil
}

// CreateIndexes creates all required indexes for the search-api collections
func CreateIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ==================== TRIPS COLLECTION INDEXES ====================
	tripsCollection := db.Collection("trips")

	tripIndexes := []mongo.IndexModel{
		// Compound index for status and departure time filtering
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "departure_datetime", Value: 1},
			},
		},
		// Compound index for city-to-city route searches
		{
			Keys: bson.D{
				{Key: "origin.city", Value: 1},
				{Key: "destination.city", Value: 1},
			},
		},
		// 2dsphere index for geospatial queries on origin coordinates
		{
			Keys: bson.D{
				{Key: "origin.coordinates", Value: "2dsphere"},
			},
		},
		// 2dsphere index for geospatial queries on destination coordinates
		{
			Keys: bson.D{
				{Key: "destination.coordinates", Value: "2dsphere"},
			},
		},
	}

	_, err := tripsCollection.Indexes().CreateMany(ctx, tripIndexes)
	if err != nil {
		return fmt.Errorf("failed to create trips indexes: %w", err)
	}
	log.Println("✅ Trips collection indexes created successfully")

	// ==================== PROCESSED_EVENTS COLLECTION INDEXES ====================
	eventsCollection := db.Collection("processed_events")

	eventIndexes := []mongo.IndexModel{
		// UNIQUE index on event_id for idempotency - CRITICAL
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// Index for filtering by event type
		{
			Keys: bson.D{{Key: "event_type", Value: 1}},
		},
		// Index for timestamp-based queries
		{
			Keys: bson.D{{Key: "processed_at", Value: 1}},
		},
	}

	_, err = eventsCollection.Indexes().CreateMany(ctx, eventIndexes)
	if err != nil {
		return fmt.Errorf("failed to create processed_events indexes: %w", err)
	}
	log.Println("✅ Processed events collection indexes created successfully")

	// ==================== POPULAR_ROUTES COLLECTION INDEXES ====================
	popularRoutesCollection := db.Collection("popular_routes")

	popularRouteIndexes := []mongo.IndexModel{
		// UNIQUE compound index on origin_city and destination_city
		{
			Keys: bson.D{
				{Key: "origin_city", Value: 1},
				{Key: "destination_city", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = popularRoutesCollection.Indexes().CreateMany(ctx, popularRouteIndexes)
	if err != nil {
		return fmt.Errorf("failed to create popular_routes indexes: %w", err)
	}
	log.Println("✅ Popular routes collection indexes created successfully")

	log.Println("✅ All MongoDB indexes created successfully")
	return nil
}
