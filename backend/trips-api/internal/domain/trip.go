package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Trip representa un viaje en el sistema de carpooling
type Trip struct {
	ID                       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DriverID                 int64              `json:"driver_id" bson:"driver_id"`

	Origin                   Location `json:"origin" bson:"origin"`
	Destination              Location `json:"destination" bson:"destination"`

	DepartureDatetime        time.Time `json:"departure_datetime" bson:"departure_datetime"`
	EstimatedArrivalDatetime time.Time `json:"estimated_arrival_datetime" bson:"estimated_arrival_datetime"`

	PricePerSeat             float64     `json:"price_per_seat" bson:"price_per_seat"`
	TotalSeats               int         `json:"total_seats" bson:"total_seats"`
	ReservedSeats            int         `json:"reserved_seats" bson:"reserved_seats"`
	AvailableSeats           int         `json:"available_seats" bson:"available_seats"`
	AvailabilityVersion      int         `json:"availability_version" bson:"availability_version"` // For optimistic locking

	Car         Car         `json:"car" bson:"car"`
	Preferences Preferences `json:"preferences" bson:"preferences"`

	Status      string `json:"status" bson:"status"` // draft, published, full, in_progress, completed, cancelled
	Description string `json:"description" bson:"description"`

	CancelledAt        *time.Time `json:"cancelled_at,omitempty" bson:"cancelled_at,omitempty"`
	CancelledBy        *int64     `json:"cancelled_by,omitempty" bson:"cancelled_by,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateTripRequest representa la solicitud para crear un nuevo viaje
type CreateTripRequest struct {
	Origin                   Location    `json:"origin" binding:"required"`
	Destination              Location    `json:"destination" binding:"required"`
	DepartureDatetime        string      `json:"departure_datetime" binding:"required"`        // RFC3339 format
	EstimatedArrivalDatetime string      `json:"estimated_arrival_datetime" binding:"required"` // RFC3339 format
	PricePerSeat             float64     `json:"price_per_seat" binding:"required,min=0"`
	TotalSeats               int         `json:"total_seats" binding:"required,min=1,max=8"`
	Car                      Car         `json:"car" binding:"required"`
	Preferences              Preferences `json:"preferences"`
	Description              string      `json:"description"`
}

// UpdateTripRequest representa la solicitud para actualizar un viaje existente
type UpdateTripRequest struct {
	Origin                   *Location    `json:"origin"`
	Destination              *Location    `json:"destination"`
	DepartureDatetime        *string      `json:"departure_datetime"`        // RFC3339 format
	EstimatedArrivalDatetime *string      `json:"estimated_arrival_datetime"` // RFC3339 format
	PricePerSeat             *float64     `json:"price_per_seat"`
	TotalSeats               *int         `json:"total_seats"`
	Car                      *Car         `json:"car"`
	Preferences              *Preferences `json:"preferences"`
	Description              *string      `json:"description"`
}

