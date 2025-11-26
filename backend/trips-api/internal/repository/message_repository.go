package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"trips-api/internal/dao"
)

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	Create(ctx context.Context, message *dao.Message) error
	FindByTripID(ctx context.Context, tripID string, limit int) ([]*dao.Message, error)
}

type mongoMessageRepository struct {
	db *mongo.Database
}

// NewMessageRepository creates a new MongoDB message repository
func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &mongoMessageRepository{db: db}
}

// Create saves a new message to the database
func (r *mongoMessageRepository) Create(ctx context.Context, message *dao.Message) error {
	collection := r.db.Collection(dao.Message{}.CollectionName())
	message.CreatedAt = time.Now()
	_, err := collection.InsertOne(ctx, message)
	return err
}

// FindByTripID retrieves messages for a specific trip
func (r *mongoMessageRepository) FindByTripID(ctx context.Context, tripID string, limit int) ([]*dao.Message, error) {
	collection := r.db.Collection(dao.Message{}.CollectionName())

	filter := bson.M{"trip_id": tripID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}). // Most recent first
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*dao.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
