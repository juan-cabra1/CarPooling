package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PopularRoute tracks popular routes for trending and analytics
// Unique index on (origin_city, destination_city) ensures one record per route
type PopularRoute struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OriginCity      string             `json:"origin_city" bson:"origin_city"`
	DestinationCity string             `json:"destination_city" bson:"destination_city"`
	SearchCount     int                `json:"search_count" bson:"search_count"`
	LastSearched    time.Time          `json:"last_searched" bson:"last_searched"`
}
