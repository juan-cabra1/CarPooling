package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Booking status constants
// These define the possible states of a booking throughout its lifecycle
const (
	// BookingStatusPending - Initial state when booking is created, waiting for trip confirmation
	BookingStatusPending = "pending"

	// BookingStatusConfirmed - Trip has confirmed the reservation, seats are reserved
	BookingStatusConfirmed = "confirmed"

	// BookingStatusCancelled - Booking was cancelled by passenger or due to trip cancellation
	BookingStatusCancelled = "cancelled"

	// BookingStatusCompleted - Trip has been completed successfully
	BookingStatusCompleted = "completed"

	// BookingStatusFailed - Reservation failed (e.g., no seats available)
	BookingStatusFailed = "failed"
)

// Booking represents a passenger's reservation for a trip in the database
//
// This is the GORM model (DAO - Data Access Object) that maps directly to the MySQL table.
// It uses GORM struct tags to define:
//   - Database column types and constraints
//   - Indexes for query performance
//   - JSON serialization format for API responses
//
// Key Design Decisions:
//   - ID: Internal auto-increment primary key (not exposed to external APIs)
//   - BookingUUID: External unique identifier (UUID v4) used in API responses
//   - TripID: Reference to trip (stored as string, MongoDB ObjectID from trips-api)
//   - PassengerID: Reference to user (stored as int64 from users-api)
//   - Price: Total booking price (seats_requested * price_per_seat from trip)
//   - Status: Current state with index for efficient filtering
//   - CancelledAt: Nullable timestamp, only set when status becomes 'cancelled'
//   - CancellationReason: Optional text explaining why booking was cancelled
//
// Indexes:
//   - booking_uuid (unique): Fast lookup by external ID
//   - trip_id: Find all bookings for a trip
//   - passenger_id: Find all bookings by a user
//   - status: Filter by booking state (confirmed, cancelled, etc.)
type Booking struct {
	// ID is the internal database primary key (auto-increment)
	// Not exposed in API responses (use BookingUUID instead)
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// BookingUUID is the external unique identifier (UUID v4)
	// This is used in all API responses and external references
	// Generated automatically in BeforeCreate hook if not set
	BookingUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"id"`

	// TripID references a trip from trips-api (MongoDB ObjectID as string)
	// Indexed for efficient queries like "find all bookings for this trip"
	TripID string `gorm:"type:varchar(36);index;not null" json:"trip_id"`

	// PassengerID references a user from users-api (user ID as int64)
	// Indexed for efficient queries like "find all my bookings"
	PassengerID int64 `gorm:"index;not null" json:"passenger_id"`

	// DriverID references the trip's driver from users-api (user ID as int64)
	// Populated when reservation is confirmed via reservation.confirmed event
	// Used for authorization checks (driver can also cancel bookings on their trips)
	// Default 0 for pending bookings (populated on confirmation)
	DriverID int64 `gorm:"index;default:0" json:"driver_id"`

	// SeatsRequested is the number of seats reserved for this booking
	// Must be between 1 and available seats on the trip
	SeatsRequested int `gorm:"not null" json:"seats_requested"`

	// TotalPrice is the total cost (seats_requested * price_per_seat)
	// Stored as DECIMAL(10,2) for precise currency calculations
	// Example: 2 seats * $5000.00 = $10000.00
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`

	// Status is the current state of the booking
	// Indexed for efficient filtering (e.g., "show only confirmed bookings")
	// Possible values: pending, confirmed, cancelled, completed, failed
	// See constants: BookingStatusPending, BookingStatusConfirmed, etc.
	Status string `gorm:"type:varchar(20);index;not null;default:pending" json:"status"`

	// CancelledAt is the timestamp when booking was cancelled (nullable)
	// Only populated when Status becomes 'cancelled'
	CancelledAt *time.Time `gorm:"index" json:"cancelled_at,omitempty"`

	// CancellationReason explains why booking was cancelled (nullable)
	// Examples: "Change of plans", "Trip cancelled by driver", etc.
	CancellationReason string `gorm:"type:text" json:"cancellation_reason,omitempty"`

	// CreatedAt is automatically managed by GORM (timestamp when row inserted)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// UpdatedAt is automatically managed by GORM (timestamp when row updated)
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the custom table name for the Booking model
//
// GORM convention: By default, GORM would create table "bookings" (pluralized struct name)
// We explicitly define it here for clarity and to avoid naming ambiguity
//
// Returns:
//   - "bookings": The MySQL table name where booking records are stored
func (Booking) TableName() string {
	return "bookings"
}

// BeforeCreate is a GORM hook that runs before inserting a new booking record
//
// This hook automatically generates a UUID v4 for BookingUUID if it's not already set.
// UUIDs are used as external identifiers because:
//   - They're globally unique (safe for distributed systems)
//   - They don't expose internal database auto-increment IDs
//   - They're compatible with external APIs and microservices
//
// Hook execution order:
//  1. Application calls db.Create(&booking)
//  2. GORM calls BeforeCreate (this function)
//  3. UUID generated if BookingUUID is empty
//  4. GORM inserts row into database with generated UUID
//
// Parameters:
//   - tx: GORM database transaction (not used here, but required by interface)
//
// Returns:
//   - error: Always nil in current implementation (UUID generation cannot fail)
//
// Example:
//
//	booking := &Booking{TripID: "trip123", PassengerID: 456}
//	db.Create(booking)
//	// booking.BookingUUID is now set to a UUID like "550e8400-e29b-41d4-a716-446655440000"
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID only if not already set (allows manual UUID assignment for testing)
	if b.BookingUUID == "" {
		b.BookingUUID = uuid.New().String()
	}
	return nil
}

// Helper methods for status validation (optional but recommended)

// IsPending checks if booking is in pending state (waiting for confirmation)
func (b *Booking) IsPending() bool {
	return b.Status == BookingStatusPending
}

// IsConfirmed checks if booking is confirmed (seats reserved)
func (b *Booking) IsConfirmed() bool {
	return b.Status == BookingStatusConfirmed
}

// IsCancelled checks if booking was cancelled
func (b *Booking) IsCancelled() bool {
	return b.Status == BookingStatusCancelled
}

// IsCompleted checks if trip was completed successfully
func (b *Booking) IsCompleted() bool {
	return b.Status == BookingStatusCompleted
}

// IsFailed checks if reservation failed (e.g., no seats available)
func (b *Booking) IsFailed() bool {
	return b.Status == BookingStatusFailed
}

// CanBeCancelled checks if booking can be cancelled by user
// Rules: Can only cancel if status is 'pending' or 'confirmed'
func (b *Booking) CanBeCancelled() bool {
	return b.IsPending() || b.IsConfirmed()
}
