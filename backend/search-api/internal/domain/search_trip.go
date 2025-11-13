package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchTrip represents a denormalized trip document stored in MongoDB for search purposes
// This includes driver information and search-specific fields
type SearchTrip struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Original trip ID from trips-api
	TripID   string `json:"trip_id" bson:"trip_id"`
	DriverID int64  `json:"driver_id" bson:"driver_id"`

	// Denormalized driver information (fetched from users-api)
	Driver Driver `json:"driver" bson:"driver"`

	// Location information
	Origin      Location `json:"origin" bson:"origin"`
	Destination Location `json:"destination" bson:"destination"`

	// Trip timing
	DepartureDatetime        time.Time `json:"departure_datetime" bson:"departure_datetime"`
	EstimatedArrivalDatetime time.Time `json:"estimated_arrival_datetime" bson:"estimated_arrival_datetime"`

	// Pricing and availability
	PricePerSeat   float64 `json:"price_per_seat" bson:"price_per_seat"`
	TotalSeats     int     `json:"total_seats" bson:"total_seats"`
	AvailableSeats int     `json:"available_seats" bson:"available_seats"`

	// Vehicle and preferences
	Car         Car         `json:"car" bson:"car"`
	Preferences Preferences `json:"preferences" bson:"preferences"`

	// Trip details
	Status      string `json:"status" bson:"status"` // published, full, in_progress, completed, cancelled
	Description string `json:"description" bson:"description"`

	// Search-specific fields
	SearchText      string  `json:"search_text,omitempty" bson:"search_text,omitempty"`           // Concatenated text for backup text search
	PopularityScore float64 `json:"popularity_score,omitempty" bson:"popularity_score,omitempty"` // For ranking popular trips

	// Timestamps
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
