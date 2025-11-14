package repository

import (
	"context"
	"time"

	"search-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventRepository handles idempotency tracking for processed events
type EventRepository interface {
	IsEventProcessed(ctx context.Context, eventID string) (bool, error)
	MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error
}

type eventRepository struct {
	collection *mongo.Collection
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *mongo.Database) EventRepository {
	return &eventRepository{
		collection: db.Collection("processed_events"),
	}
}

// IsEventProcessed checks if an event has already been processed
func (r *eventRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var event domain.ProcessedEvent
	err := r.collection.FindOne(ctx, bson.M{"event_id": eventID}).Decode(&event)

	if err == mongo.ErrNoDocuments {
		return false, nil // Event not processed yet
	}

	if err != nil {
		return false, err // Database error
	}

	return true, nil // Event already processed
}

// MarkEventProcessed marks an event as processed
// If the event was already processed (duplicate key error), it returns nil (idempotent)
func (r *eventRepository) MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if event.ProcessedAt.IsZero() {
		event.ProcessedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, event)

	// If duplicate key error, event was already processed - return success (idempotency)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}

	return err
}
