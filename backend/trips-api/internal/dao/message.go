package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a chat message in a trip
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TripID    string             `bson:"trip_id" json:"trip_id"`
	UserID    int64              `bson:"user_id" json:"user_id"`
	UserName  string             `bson:"user_name" json:"user_name"`
	Message   string             `bson:"message" json:"message"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// CollectionName returns the MongoDB collection name for messages
func (Message) CollectionName() string {
	return "messages"
}
