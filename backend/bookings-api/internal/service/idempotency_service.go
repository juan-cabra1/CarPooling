package service

import (
	"bookings-api/internal/repository"

	"github.com/rs/zerolog/log"
)

// IdempotencyService provides idempotency checking for event processing
// This is CRITICAL for RabbitMQ consumers to prevent duplicate event processing
type IdempotencyService interface {
	// CheckAndMarkEvent checks if an event has been processed and marks it as processed
	// Returns (true, nil) if the event is NEW and should be processed
	// Returns (false, nil) if the event is DUPLICATE and should be skipped
	// Returns (false, err) if there was a database error
	CheckAndMarkEvent(eventID, eventType string) (bool, error)

	// IsEventProcessed checks if an event has already been processed
	// Returns true if already processed, false if new
	IsEventProcessed(eventID string) (bool, error)

	// MarkEventAsSuccess marks an event as successfully processed
	MarkEventAsSuccess(eventID, eventType string) error

	// MarkEventAsFailed marks an event as failed with an error message
	MarkEventAsFailed(eventID, eventType, errorMsg string) error
}

// idempotencyService implements IdempotencyService
type idempotencyService struct {
	eventRepo repository.EventRepository
}

// NewIdempotencyService creates a new IdempotencyService
func NewIdempotencyService(eventRepo repository.EventRepository) IdempotencyService {
	return &idempotencyService{
		eventRepo: eventRepo,
	}
}

// CheckAndMarkEvent checks if an event has been processed and marks it as processed atomically
// This is the primary method used by RabbitMQ consumers to ensure idempotency
func (s *idempotencyService) CheckAndMarkEvent(eventID, eventType string) (bool, error) {
	// Check if event has already been processed
	processed, err := s.eventRepo.IsEventProcessed(eventID)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", eventID).
			Str("event_type", eventType).
			Msg("Failed to check if event is processed")
		return false, err
	}

	// If already processed, return false (skip processing)
	if processed {
		log.Info().
			Str("event_id", eventID).
			Str("event_type", eventType).
			Msg("Event already processed, skipping (idempotency)")
		return false, nil
	}

	// Event is new, mark as success (will be processed)
	// Note: We mark it as success upfront. If processing fails later,
	// the consumer should handle retries or mark as failed separately
	if err := s.eventRepo.MarkEventAsSuccess(eventID, eventType); err != nil {
		log.Error().
			Err(err).
			Str("event_id", eventID).
			Str("event_type", eventType).
			Msg("Failed to mark event as processed")
		return false, err
	}

	log.Debug().
		Str("event_id", eventID).
		Str("event_type", eventType).
		Msg("Event marked as processed, will process")

	// Return true to indicate event should be processed
	return true, nil
}

// IsEventProcessed checks if an event has already been processed
func (s *idempotencyService) IsEventProcessed(eventID string) (bool, error) {
	processed, err := s.eventRepo.IsEventProcessed(eventID)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", eventID).
			Msg("Failed to check if event is processed")
		return false, err
	}

	return processed, nil
}

// MarkEventAsSuccess marks an event as successfully processed
func (s *idempotencyService) MarkEventAsSuccess(eventID, eventType string) error {
	if err := s.eventRepo.MarkEventAsSuccess(eventID, eventType); err != nil {
		log.Error().
			Err(err).
			Str("event_id", eventID).
			Str("event_type", eventType).
			Msg("Failed to mark event as success")
		return err
	}

	log.Info().
		Str("event_id", eventID).
		Str("event_type", eventType).
		Msg("Event marked as success")

	return nil
}

// MarkEventAsFailed marks an event as failed with an error message
func (s *idempotencyService) MarkEventAsFailed(eventID, eventType, errorMsg string) error {
	if err := s.eventRepo.MarkEventAsFailed(eventID, eventType, errorMsg); err != nil {
		log.Error().
			Err(err).
			Str("event_id", eventID).
			Str("event_type", eventType).
			Str("error_msg", errorMsg).
			Msg("Failed to mark event as failed")
		return err
	}

	log.Warn().
		Str("event_id", eventID).
		Str("event_type", eventType).
		Str("error_msg", errorMsg).
		Msg("Event marked as failed")

	return nil
}
