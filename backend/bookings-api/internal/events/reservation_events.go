package events

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// OUTBOUND EVENT TYPE CONSTANTS (Published by bookings-api)
// ============================================================================
// These constants define the types of events that bookings-api publishes
// to RabbitMQ for consumption by trips-api.
//
// Event Flow:
//   1. User creates booking → reservation.created event published
//   2. User cancels booking → reservation.cancelled event published
//   3. trips-api consumes events and updates available seats
//
// Event naming convention: <resource>.<action>
const (
	// EventTypeReservationCreated - Published when a new booking is created
	// trips-api will decrement available seats upon receiving this event
	EventTypeReservationCreated = "reservation.created"

	// EventTypeReservationCancelled - Published when a booking is cancelled
	// trips-api will increment available seats upon receiving this event
	EventTypeReservationCancelled = "reservation.cancelled"
)

// ============================================================================
// BASE EVENT STRUCTURE
// ============================================================================

// BaseEvent contains fields common to all events published by bookings-api
//
// This base structure ensures:
//   - Event idempotency via EventID (UUID v4)
//   - Event type identification
//   - Event timestamp for ordering and debugging
//
// All specific event types embed this base structure to inherit these fields.
//
// Idempotency:
// The EventID is a unique UUID v4 generated for each event. trips-api uses
// this ID to track processed events and prevent duplicate processing.
type BaseEvent struct {
	// EventID is a unique identifier (UUID v4) for this event
	// Used by consumers for idempotency (prevent duplicate processing)
	// Generated automatically by NewBaseEvent()
	EventID string `json:"event_id"`

	// EventType categorizes the event (e.g., "reservation.created")
	// Consumers use this field to route events to appropriate handlers
	EventType string `json:"event_type"`

	// Timestamp indicates when the event was generated
	// Used for event ordering, debugging, and auditing
	// Set automatically by NewBaseEvent() to current UTC time
	Timestamp time.Time `json:"timestamp"`
}

// ============================================================================
// RESERVATION CREATED EVENT (Outbound from bookings-api)
// ============================================================================

// ReservationCreatedEvent is published when a new booking is created
//
// This event informs trips-api that seats have been reserved on a trip.
// trips-api will:
//   1. Validate the event (check for duplicates using EventID)
//   2. Decrement available_seats by SeatsReserved
//   3. Persist the booking reference
//   4. ACK the message if successful
//
// Idempotency:
// If trips-api receives the same EventID twice (e.g., due to retry),
// it will skip processing to prevent double-decrementing seats.
type ReservationCreatedEvent struct {
	// Embed BaseEvent to inherit EventID, EventType, Timestamp
	BaseEvent

	// TripID identifies which trip the reservation is for
	// This is a MongoDB ObjectID (string) from trips-api
	TripID string `json:"trip_id"`

	// SeatsReserved indicates how many seats were reserved
	// trips-api will decrement available_seats by this amount
	// Must be > 0 and <= trip.available_seats
	SeatsReserved int `json:"seats_reserved"`

	// ReservationID is the booking UUID from bookings-api
	// Used for tracking and debugging (links event to booking record)
	ReservationID string `json:"reservation_id"`
}

// ============================================================================
// RESERVATION CANCELLED EVENT (Outbound from bookings-api)
// ============================================================================

// ReservationCancelledEvent is published when a booking is cancelled
//
// This event informs trips-api that previously reserved seats are now available.
// trips-api will:
//   1. Validate the event (check for duplicates using EventID)
//   2. Increment available_seats by SeatsReleased
//   3. Remove the booking reference
//   4. ACK the message if successful
//
// Idempotency:
// If trips-api receives the same EventID twice (e.g., due to retry),
// it will skip processing to prevent double-incrementing seats.
type ReservationCancelledEvent struct {
	// Embed BaseEvent to inherit EventID, EventType, Timestamp
	BaseEvent

	// TripID identifies which trip the cancellation is for
	// This is a MongoDB ObjectID (string) from trips-api
	TripID string `json:"trip_id"`

	// SeatsReleased indicates how many seats are being freed
	// trips-api will increment available_seats by this amount
	// Must match the original SeatsReserved from ReservationCreatedEvent
	SeatsReleased int `json:"seats_released"`

	// ReservationID is the booking UUID from bookings-api
	// Used for tracking and debugging (links event to booking record)
	ReservationID string `json:"reservation_id"`
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// NewBaseEvent creates a new BaseEvent with auto-generated ID and timestamp
//
// This function is a convenience constructor that:
//   1. Generates a unique UUID v4 for EventID (used for idempotency)
//   2. Sets Timestamp to current UTC time
//   3. Sets EventType to the provided value
//
// Usage:
//
//	baseEvent := NewBaseEvent(EventTypeReservationCreated)
//	event := ReservationCreatedEvent{
//	    BaseEvent:     baseEvent,
//	    TripID:        "trip-123",
//	    SeatsReserved: 2,
//	    ReservationID: booking.BookingUUID,
//	}
//
// Parameters:
//   - eventType: The type of event (use constants like EventTypeReservationCreated)
//
// Returns:
//   - BaseEvent: Populated with EventID, EventType, and Timestamp
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		EventID:   uuid.New().String(), // Generate UUID v4 for idempotency
		EventType: eventType,
		Timestamp: time.Now().UTC(), // Use UTC for consistency across timezones
	}
}
