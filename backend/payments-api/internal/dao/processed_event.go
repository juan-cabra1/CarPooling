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
// If it exists -> skip processing (idempotent behavior)
// If it doesn't exist -> process and insert record
//
// The UNIQUE constraint on event_id ensures that even with concurrent
// consumers trying to process the same event, only ONE will succeed in
// inserting the record (the others will get a duplicate key error).
type ProcessedEvent struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// EventID is the unique identifier of the event (UUID from RabbitMQ message)
	// This field has a UNIQUE constraint - this is the CORE of idempotency
	EventID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"event_id"`

	// EventType categorizes the event for logging and debugging
	// Examples: "booking.cancelled", "trip.completed", "payment.approved"
	EventType string `gorm:"type:varchar(50);index;not null" json:"event_type"`

	// Result stores the outcome of processing this event
	// Possible values: "success", "skipped", "failed"
	Result string `gorm:"type:varchar(20);not null" json:"result"`

	// ErrorMessage stores the error details if Result is "failed"
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`

	// ProcessedAt is the timestamp when event was processed
	ProcessedAt time.Time `gorm:"autoCreateTime;index" json:"processed_at"`

	// CreatedAt is automatically managed by GORM
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the custom table name
func (ProcessedEvent) TableName() string {
	return "processed_events"
}

// Event result constants
const (
	// EventResultSuccess - Event was processed successfully
	EventResultSuccess = "success"

	// EventResultSkipped - Event was skipped (duplicate)
	EventResultSkipped = "skipped"

	// EventResultFailed - Event processing failed
	EventResultFailed = "failed"
)

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
