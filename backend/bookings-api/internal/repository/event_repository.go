package repository

import (
	"bookings-api/internal/dao"
	"strings"

	"gorm.io/gorm"
)

// EventRepository defines the interface for event idempotency operations
type EventRepository interface {
	// IsEventProcessed checks if an event has already been processed
	// Returns true if the event exists, false otherwise
	IsEventProcessed(eventID string) (bool, error)

	// MarkEventAsSuccess marks an event as successfully processed
	// If the event already exists (duplicate), returns nil (idempotency)
	MarkEventAsSuccess(eventID, eventType string) error

	// MarkEventAsFailed marks an event as failed with an error message
	// If the event already exists (duplicate), returns nil (idempotency)
	MarkEventAsFailed(eventID, eventType, errorMsg string) error
}

// eventRepository implements EventRepository using GORM
type eventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new instance of EventRepository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

// IsEventProcessed checks if an event has already been processed
func (r *eventRepository) IsEventProcessed(eventID string) (bool, error) {
	var count int64
	err := r.db.Model(&dao.ProcessedEvent{}).
		Where("event_id = ?", eventID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// MarkEventAsSuccess marks an event as successfully processed
func (r *eventRepository) MarkEventAsSuccess(eventID, eventType string) error {
	event := &dao.ProcessedEvent{
		EventID:   eventID,
		EventType: eventType,
		Result:    dao.EventResultSuccess,
	}

	err := r.db.Create(event).Error
	if err != nil {
		// CRITICAL: Handle duplicate key errors for idempotency
		// If the event_id already exists (UNIQUE constraint violation),
		// treat it as success (event already processed)
		if isDuplicateKeyError(err) {
			return nil // Idempotency: event already processed, return success
		}
		return err
	}

	return nil
}

// MarkEventAsFailed marks an event as failed with an error message
func (r *eventRepository) MarkEventAsFailed(eventID, eventType, errorMsg string) error {
	event := &dao.ProcessedEvent{
		EventID:      eventID,
		EventType:    eventType,
		Result:       dao.EventResultFailed,
		ErrorMessage: errorMsg,
	}

	err := r.db.Create(event).Error
	if err != nil {
		// CRITICAL: Handle duplicate key errors for idempotency
		// If the event_id already exists (UNIQUE constraint violation),
		// treat it as success (event already processed)
		if isDuplicateKeyError(err) {
			return nil // Idempotency: event already processed, return success
		}
		return err
	}

	return nil
}

// isDuplicateKeyError checks if the error is a MySQL duplicate key error
// MySQL returns error 1062 for duplicate entry violations
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Check for MySQL duplicate entry error message
	// Error 1062: Duplicate entry 'value' for key 'index_name'
	return strings.Contains(errMsg, "Duplicate entry") ||
		strings.Contains(errMsg, "1062") ||
		strings.Contains(errMsg, "duplicate key")
}
