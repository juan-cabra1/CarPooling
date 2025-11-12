package dao

import (
	"time"
)

// ProcessedEvent represents an event that has been processed by the system
//
// This table is CRITICAL for implementing idempotency in event-driven systems.
//
// Problem:
// RabbitMQ may deliver the same message multiple times due to:
//   - Network failures during ACK
//   - Consumer crashes after processing but before ACK
//   - Message broker restarts
//
// Solution:
// Before processing any event, we check if it's already in this table.
// If it exists → skip processing (idempotent behavior)
// If it doesn't exist → process and insert record
//
// The UNIQUE constraint on event_id ensures that even with concurrent
// consumers trying to process the same event, only ONE will succeed in
// inserting the record (the others will get a duplicate key error).
//
// Example flow:
//  1. Consumer receives event with event_id="abc-123"
//  2. BEGIN TRANSACTION
//  3. Try to INSERT INTO processed_events (event_id="abc-123")
//  4. If INSERT succeeds → process event, COMMIT
//  5. If INSERT fails (duplicate key) → skip processing, ROLLBACK
//  6. ACK message to RabbitMQ
//
// This ensures exactly-once processing semantics even with at-least-once delivery.
type ProcessedEvent struct {
	// ID is the internal database primary key (auto-increment)
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// EventID is the unique identifier of the event (UUID from RabbitMQ message)
	// This field has a UNIQUE constraint - this is the CORE of idempotency
	//
	// MySQL will reject duplicate event_id inserts with error:
	// "Error 1062: Duplicate entry 'abc-123' for key 'event_id'"
	//
	// We use this error to detect duplicate events and skip reprocessing
	EventID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"event_id"`

	// EventType categorizes the event for logging and debugging
	// Examples: "trip.updated", "trip.cancelled", "reservation.failed"
	//
	// This helps with:
	//   - Debugging: "show me all trip.cancelled events we've processed"
	//   - Monitoring: "count events by type in last hour"
	//   - Auditing: "when did we last process a trip.updated event?"
	EventType string `gorm:"type:varchar(50);index;not null" json:"event_type"`

	// Result stores the outcome of processing this event
	// Possible values: "success", "skipped", "failed"
	//
	// - "success": Event was processed and booking updated
	// - "skipped": Duplicate event, skipped to maintain idempotency
	// - "failed": Error during processing (may be retried)
	//
	// This field is useful for troubleshooting and metrics
	Result string `gorm:"type:varchar(20);not null" json:"result"`

	// ErrorMessage stores the error details if Result is "failed"
	// Nullable field - only populated when processing fails
	//
	// Examples:
	//   - "Booking not found: uuid-123"
	//   - "Database connection timeout"
	//   - "Invalid event payload: missing trip_id"
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`

	// ProcessedAt is the timestamp when event was processed
	// Automatically set by GORM on insert
	//
	// Indexed for queries like:
	//   - "Show events processed in last hour"
	//   - "Clean up old processed_events records older than 30 days"
	ProcessedAt time.Time `gorm:"autoCreateTime;index" json:"processed_at"`

	// CreatedAt is automatically managed by GORM (same as ProcessedAt in this case)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the custom table name for the ProcessedEvent model
//
// Returns:
//   - "processed_events": The MySQL table name for event idempotency tracking
func (ProcessedEvent) TableName() string {
	return "processed_events"
}

// Event result constants for the Result field
const (
	// EventResultSuccess - Event was processed successfully and changes were applied
	EventResultSuccess = "success"

	// EventResultSkipped - Event was skipped because it was already processed (idempotency)
	EventResultSkipped = "skipped"

	// EventResultFailed - Event processing failed (error in ErrorMessage field)
	EventResultFailed = "failed"
)

// Helper methods for result checking

// IsSuccess checks if event was processed successfully
func (pe *ProcessedEvent) IsSuccess() bool {
	return pe.Result == EventResultSuccess
}

// IsSkipped checks if event was skipped (duplicate)
func (pe *ProcessedEvent) IsSkipped() bool {
	return pe.Result == EventResultSkipped
}

// IsFailed checks if event processing failed
func (pe *ProcessedEvent) IsFailed() bool {
	return pe.Result == EventResultFailed
}
