package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProcessedEvent tracks processed events for idempotency
// Ensures that events from RabbitMQ are only processed once
type ProcessedEvent struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EventID      string             `json:"event_id" bson:"event_id"` // UNIQUE index required
	EventType    string             `json:"event_type" bson:"event_type"`
	ProcessedAt  time.Time          `json:"processed_at" bson:"processed_at"`
	Result       string             `json:"result" bson:"result"` // success, skipped, failed
	ErrorMessage string             `json:"error_message,omitempty" bson:"error_message,omitempty"`
}
