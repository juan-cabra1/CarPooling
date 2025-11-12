package messaging

import "time"

// TripCancelledEvent represents a trip cancellation event from trips-api
// Published when a driver cancels a trip
type TripCancelledEvent struct {
	EventID            string    `json:"event_id"`             // UUID for idempotency
	EventType          string    `json:"event_type"`           // "trip.cancelled"
	TripID             string    `json:"trip_id"`              // MongoDB ObjectID
	DriverID           int64     `json:"driver_id"`            // Driver who owns the trip
	Status             string    `json:"status"`               // Should be "cancelled"
	AvailableSeats     int       `json:"available_seats"`      // Remaining seats (likely 0)
	ReservedSeats      int       `json:"reserved_seats"`       // Reserved seats count
	Timestamp          time.Time `json:"timestamp"`            // Event creation time
	SourceService      string    `json:"source_service"`       // "trips-api"
	CorrelationID      string    `json:"correlation_id"`       // For request tracing
	CancelledBy        int64     `json:"cancelled_by"`         // User ID who cancelled
	CancellationReason string    `json:"cancellation_reason"`  // Human-readable reason
}

// ReservationFailedEvent represents a failed reservation attempt
// Published when trips-api cannot reserve seats for a booking
type ReservationFailedEvent struct {
	EventID        string    `json:"event_id"`         // UUID for idempotency
	EventType      string    `json:"event_type"`       // "reservation.failed"
	ReservationID  string    `json:"reservation_id"`   // Booking UUID from bookings-api
	TripID         string    `json:"trip_id"`          // MongoDB ObjectID
	Reason         string    `json:"reason"`           // Failure reason (e.g., "No seats available")
	AvailableSeats int       `json:"available_seats"`  // Current available seats
	SourceService  string    `json:"source_service"`   // "trips-api"
	CorrelationID  string    `json:"correlation_id"`   // For request tracing
	Timestamp      time.Time `json:"timestamp"`        // Event creation time
}

// ReservationConfirmedEvent represents a successful reservation confirmation
// Published when trips-api successfully reserves seats for a booking
type ReservationConfirmedEvent struct {
	EventID        string    `json:"event_id"`         // UUID for idempotency
	EventType      string    `json:"event_type"`       // "reservation.confirmed"
	ReservationID  string    `json:"reservation_id"`   // Booking UUID from bookings-api
	TripID         string    `json:"trip_id"`          // MongoDB ObjectID
	PassengerID    int64     `json:"passenger_id"`     // User ID of passenger
	SeatsReserved  int       `json:"seats_reserved"`   // Number of seats reserved
	TotalPrice     float64   `json:"total_price"`      // Total price for the reservation
	AvailableSeats int       `json:"available_seats"`  // Remaining seats after reservation
	SourceService  string    `json:"source_service"`   // "trips-api"
	CorrelationID  string    `json:"correlation_id"`   // For request tracing
	Timestamp      time.Time `json:"timestamp"`        // Event creation time
}
