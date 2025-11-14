package messaging

import "time"

// TripCreatedEvent represents a trip creation event from trips-api
type TripCreatedEvent struct {
	EventID           string    `json:"event_id"`
	EventType         string    `json:"event_type"`
	TripID            string    `json:"trip_id"`
	DriverID          int64     `json:"driver_id"`
	OriginCity        string    `json:"origin_city"`
	DestinationCity   string    `json:"destination_city"`
	DepartureDatetime time.Time `json:"departure_datetime"`
	AvailableSeats    int       `json:"available_seats"`
	Status            string    `json:"status"`
	Timestamp         time.Time `json:"timestamp"`
}

// TripUpdatedEvent represents a trip update event from trips-api
type TripUpdatedEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	TripID         string    `json:"trip_id"`
	AvailableSeats int       `json:"available_seats"`
	ReservedSeats  int       `json:"reserved_seats"`
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
}

// TripCancelledEvent represents a trip cancellation event from trips-api
type TripCancelledEvent struct {
	EventID            string    `json:"event_id"`
	EventType          string    `json:"event_type"`
	TripID             string    `json:"trip_id"`
	CancellationReason string    `json:"cancellation_reason"`
	Timestamp          time.Time `json:"timestamp"`
}
