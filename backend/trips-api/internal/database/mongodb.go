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

// ConnectMongoDB establece la conexión con MongoDB y retorna la base de datos
func ConnectMongoDB(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Crear cliente de MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verificar la conexión
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("✅ MongoDB connected successfully to database: %s", dbName)

	return client.Database(dbName), nil
}

// CreateIndexes crea todos los índices necesarios para las colecciones
func CreateIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// ==================== TRIPS COLLECTION INDEXES ====================
	tripsCollection := db.Collection("trips")

	tripIndexes := []mongo.IndexModel{
		// Índice para búsquedas por conductor
		{
			Keys: bson.D{{Key: "driver_id", Value: 1}},
		},
		// Índice para filtrar por estado
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		// Índice para búsquedas por fecha de salida
		{
			Keys: bson.D{{Key: "departure_datetime", Value: 1}},
		},
		// Índice compuesto para búsquedas por ciudad de origen y destino
		{
			Keys: bson.D{
				{Key: "origin.city", Value: 1},
				{Key: "destination.city", Value: 1},
			},
		},
	}

	_, err := tripsCollection.Indexes().CreateMany(ctx, tripIndexes)
	if err != nil {
		return fmt.Errorf("failed to create trips indexes: %w", err)
	}

	log.Println("✅ Trips collection indexes created")

	// ==================== PROCESSED_EVENTS COLLECTION INDEXES ====================
	// CRITICAL: Esta colección es esencial para garantizar idempotencia
	eventsCollection := db.Collection("processed_events")

	eventIndexes := []mongo.IndexModel{
		// ÍNDICE ÚNICO EN event_id - CRÍTICO PARA IDEMPOTENCIA
		// Este índice previene el procesamiento duplicado de eventos
		{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		// Índice para filtrar por tipo de evento
		{
			Keys: bson.D{{Key: "event_type", Value: 1}},
		},
		// Índice para ordenar por fecha de procesamiento
		{
			Keys: bson.D{{Key: "processed_at", Value: 1}},
		},
	}

	_, err = eventsCollection.Indexes().CreateMany(ctx, eventIndexes)
	if err != nil {
		return fmt.Errorf("failed to create processed_events indexes: %w", err)
	}

	log.Println("Processed_events collection indexes created (UNIQUE constraint on event_id)")

	return nil
}
