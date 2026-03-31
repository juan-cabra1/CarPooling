package repository

import (
	"payments-api/internal/dao"

	"gorm.io/gorm"
)

// EventRepository handles processed event database operations
type EventRepository interface {
	Create(event *dao.ProcessedEvent) error
	Exists(eventID string) (bool, error)
	FindByEventID(eventID string) (*dao.ProcessedEvent, error)

	// With transaction support
	CreateWithTx(tx *gorm.DB, event *dao.ProcessedEvent) error
}

type eventRepository struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

// Create creates a new processed event record
func (r *eventRepository) Create(event *dao.ProcessedEvent) error {
	return r.db.Create(event).Error
}

// Exists checks if an event has already been processed
func (r *eventRepository) Exists(eventID string) (bool, error) {
	var count int64
	err := r.db.Model(&dao.ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

// FindByEventID finds a processed event by its event ID
func (r *eventRepository) FindByEventID(eventID string) (*dao.ProcessedEvent, error) {
	var event dao.ProcessedEvent
	err := r.db.Where("event_id = ?", eventID).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// CreateWithTx creates a processed event within a transaction
func (r *eventRepository) CreateWithTx(tx *gorm.DB, event *dao.ProcessedEvent) error {
	return tx.Create(event).Error
}
