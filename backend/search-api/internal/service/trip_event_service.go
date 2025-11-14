package service

import (
	"context"
	"fmt"
	"time"

	"search-api/internal/cache"
	"search-api/internal/clients"
	"search-api/internal/domain"
	"search-api/internal/repository"

	"github.com/rs/zerolog/log"
)

// TripEventService handles trip events from RabbitMQ with idempotency
type TripEventService struct {
	tripRepo    repository.TripRepository
	eventRepo   repository.EventRepository
	tripsClient clients.TripsClient
	usersClient clients.UsersClient
	solrClient  *clients.SolrClient
	cache       cache.Cache
}

// NewTripEventService creates a new TripEventService
func NewTripEventService(
	tripRepo repository.TripRepository,
	eventRepo repository.EventRepository,
	tripsClient clients.TripsClient,
	usersClient clients.UsersClient,
	solrClient *clients.SolrClient,
	cache cache.Cache,
) *TripEventService {
	return &TripEventService{
		tripRepo:    tripRepo,
		eventRepo:   eventRepo,
		tripsClient: tripsClient,
		usersClient: usersClient,
		solrClient:  solrClient,
		cache:       cache,
	}
}

// HandleTripCreated processes trip.created events
func (s *TripEventService) HandleTripCreated(ctx context.Context, eventID, tripID string, driverID int64) error {
	log.Info().
		Str("event_id", eventID).
		Str("event_type", "trip.created").
		Str("trip_id", tripID).
		Int64("driver_id", driverID).
		Msg("Processing trip.created event")

	// Check idempotency - if already processed, skip
	processed, err := s.eventRepo.IsEventProcessed(ctx, eventID)
	if err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to check event idempotency")
		return fmt.Errorf("idempotency check failed: %w", err)
	}
	if processed {
		log.Info().Str("event_id", eventID).Msg("Event already processed, skipping")
		return nil
	}

	// Fetch trip data from trips-api
	trip, err := s.tripsClient.GetTrip(ctx, tripID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			// Permanent error - trip doesn't exist
			log.Warn().Str("trip_id", tripID).Msg("Trip not found in trips-api, marking event as processed")
			processedEvent := &domain.ProcessedEvent{
				EventID:     eventID,
				EventType:   "trip.created",
				ProcessedAt: time.Now(),
				Result:      "skipped",
			}
			if markErr := s.eventRepo.MarkEventProcessed(ctx, processedEvent); markErr != nil {
				log.Error().Err(markErr).Str("event_id", eventID).Msg("Failed to mark event as processed")
			}
			return domain.ErrTripNotFound
		}
		// Transient error - retry
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to fetch trip from trips-api")
		return fmt.Errorf("fetch trip failed: %w", err)
	}

	// Fetch driver data from users-api
	driver, err := s.usersClient.GetUser(ctx, driverID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			// Permanent error - driver doesn't exist
			log.Warn().Int64("driver_id", driverID).Msg("Driver not found in users-api, marking event as processed")
			processedEvent := &domain.ProcessedEvent{
				EventID:     eventID,
				EventType:   "trip.created",
				ProcessedAt: time.Now(),
				Result:      "skipped",
			}
			if markErr := s.eventRepo.MarkEventProcessed(ctx, processedEvent); markErr != nil {
				log.Error().Err(markErr).Str("event_id", eventID).Msg("Failed to mark event as processed")
			}
			return domain.ErrUserNotFound
		}
		// Transient error - retry
		log.Error().Err(err).Int64("driver_id", driverID).Msg("Failed to fetch driver from users-api")
		return fmt.Errorf("fetch driver failed: %w", err)
	}

	// Build denormalized SearchTrip using existing ToSearchTrip method
	searchTrip := trip.ToSearchTrip(driver.ToDriver())
	searchTrip.PopularityScore = 0.0 // Initial popularity score
	searchTrip.CreatedAt = time.Now()
	searchTrip.UpdatedAt = time.Now()

	// Store in MongoDB
	if err := s.tripRepo.Create(ctx, searchTrip); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to create trip in MongoDB")
		return fmt.Errorf("mongodb create failed: %w", err)
	}

	log.Info().Str("trip_id", tripID).Msg("Trip created in MongoDB successfully")

	// Index in Solr (optional - log error but continue)
	if s.solrClient != nil {
		if err := s.solrClient.Index(ctx, searchTrip); err != nil {
			log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to index trip in Solr (continuing)")
			// Don't return error - MongoDB is source of truth
		} else {
			log.Info().Str("trip_id", tripID).Msg("Trip indexed in Solr successfully")
		}
	}

	// Mark event as processed
	processedEvent := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   "trip.created",
		ProcessedAt: time.Now(),
		Result:      "success",
	}
	if err := s.eventRepo.MarkEventProcessed(ctx, processedEvent); err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to mark event as processed")
		return fmt.Errorf("mark event processed failed: %w", err)
	}

	log.Info().
		Str("event_id", eventID).
		Str("trip_id", tripID).
		Msg("trip.created event processed successfully")

	return nil
}

// HandleTripUpdated processes trip.updated events
func (s *TripEventService) HandleTripUpdated(ctx context.Context, eventID, tripID string, availableSeats, reservedSeats int, status string) error {
	log.Info().
		Str("event_id", eventID).
		Str("event_type", "trip.updated").
		Str("trip_id", tripID).
		Int("available_seats", availableSeats).
		Int("reserved_seats", reservedSeats).
		Str("status", status).
		Msg("Processing trip.updated event")

	// Check idempotency
	processed, err := s.eventRepo.IsEventProcessed(ctx, eventID)
	if err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to check event idempotency")
		return fmt.Errorf("idempotency check failed: %w", err)
	}
	if processed {
		log.Info().Str("event_id", eventID).Msg("Event already processed, skipping")
		return nil
	}

	// Update availability and status in MongoDB
	if err := s.tripRepo.UpdateAvailabilityByTripID(ctx, tripID, availableSeats, reservedSeats, status); err != nil {
		if domain.IsNotFoundError(err) {
			// Permanent error - trip doesn't exist
			log.Warn().Str("trip_id", tripID).Msg("Trip not found in MongoDB, marking event as processed")
			processedEvent := &domain.ProcessedEvent{
				EventID:     eventID,
				EventType:   "trip.updated",
				ProcessedAt: time.Now(),
				Result:      "skipped",
			}
			if markErr := s.eventRepo.MarkEventProcessed(ctx, processedEvent); markErr != nil {
				log.Error().Err(markErr).Str("event_id", eventID).Msg("Failed to mark event as processed")
			}
			return domain.ErrSearchTripNotFound
		}
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to update trip in MongoDB")
		return fmt.Errorf("mongodb update failed: %w", err)
	}

	log.Info().Str("trip_id", tripID).Msg("Trip updated in MongoDB successfully")

	// Update in Solr (optional - log error but continue)
	if s.solrClient != nil {
		searchTrip, err := s.tripRepo.FindByTripID(ctx, tripID)
		if err == nil {
			if err := s.solrClient.Index(ctx, searchTrip); err != nil {
				log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to update trip in Solr (continuing)")
			} else {
				log.Info().Str("trip_id", tripID).Msg("Trip updated in Solr successfully")
			}
		} else {
			log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to fetch trip for Solr update")
		}
	}

	// Invalidate cache for this trip
	if s.cache != nil {
		cacheKey := fmt.Sprintf("trip:%s", tripID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Error().Err(err).Str("cache_key", cacheKey).Msg("Failed to invalidate cache (continuing)")
		} else {
			log.Info().Str("cache_key", cacheKey).Msg("Cache invalidated successfully")
		}
	}

	// Mark event as processed
	processedEvent := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   "trip.updated",
		ProcessedAt: time.Now(),
		Result:      "success",
	}
	if err := s.eventRepo.MarkEventProcessed(ctx, processedEvent); err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to mark event as processed")
		return fmt.Errorf("mark event processed failed: %w", err)
	}

	log.Info().
		Str("event_id", eventID).
		Str("trip_id", tripID).
		Msg("trip.updated event processed successfully")

	return nil
}

// HandleTripCancelled processes trip.cancelled events
func (s *TripEventService) HandleTripCancelled(ctx context.Context, eventID, tripID, cancellationReason string) error {
	log.Info().
		Str("event_id", eventID).
		Str("event_type", "trip.cancelled").
		Str("trip_id", tripID).
		Str("cancellation_reason", cancellationReason).
		Msg("Processing trip.cancelled event")

	// Check idempotency
	processed, err := s.eventRepo.IsEventProcessed(ctx, eventID)
	if err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to check event idempotency")
		return fmt.Errorf("idempotency check failed: %w", err)
	}
	if processed {
		log.Info().Str("event_id", eventID).Msg("Event already processed, skipping")
		return nil
	}

	// Update status to cancelled in MongoDB
	if err := s.tripRepo.UpdateStatusByTripID(ctx, tripID, "cancelled"); err != nil {
		if domain.IsNotFoundError(err) {
			// Permanent error - trip doesn't exist
			log.Warn().Str("trip_id", tripID).Msg("Trip not found in MongoDB, marking event as processed")
			processedEvent := &domain.ProcessedEvent{
				EventID:     eventID,
				EventType:   "trip.cancelled",
				ProcessedAt: time.Now(),
				Result:      "skipped",
			}
			if markErr := s.eventRepo.MarkEventProcessed(ctx, processedEvent); markErr != nil {
				log.Error().Err(markErr).Str("event_id", eventID).Msg("Failed to mark event as processed")
			}
			return domain.ErrSearchTripNotFound
		}
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to update trip status in MongoDB")
		return fmt.Errorf("mongodb update failed: %w", err)
	}

	log.Info().Str("trip_id", tripID).Msg("Trip status updated to cancelled in MongoDB")

	// Update in Solr (optional - log error but continue)
	if s.solrClient != nil {
		searchTrip, err := s.tripRepo.FindByTripID(ctx, tripID)
		if err == nil {
			if err := s.solrClient.Index(ctx, searchTrip); err != nil {
				log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to update trip in Solr (continuing)")
			} else {
				log.Info().Str("trip_id", tripID).Msg("Trip cancelled in Solr successfully")
			}
		} else {
			log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to fetch trip for Solr update")
		}
	}

	// Invalidate cache for this trip
	if s.cache != nil {
		cacheKey := fmt.Sprintf("trip:%s", tripID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Error().Err(err).Str("cache_key", cacheKey).Msg("Failed to invalidate cache (continuing)")
		} else {
			log.Info().Str("cache_key", cacheKey).Msg("Cache invalidated successfully")
		}
	}

	// Mark event as processed
	processedEvent := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   "trip.cancelled",
		ProcessedAt: time.Now(),
		Result:      "success",
	}
	if err := s.eventRepo.MarkEventProcessed(ctx, processedEvent); err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to mark event as processed")
		return fmt.Errorf("mark event processed failed: %w", err)
	}

	log.Info().
		Str("event_id", eventID).
		Str("trip_id", tripID).
		Msg("trip.cancelled event processed successfully")

	return nil
}

// HandleTripDeleted processes trip.deleted events
func (s *TripEventService) HandleTripDeleted(ctx context.Context, eventID, tripID, reason string) error {
	log.Info().
		Str("event_id", eventID).
		Str("event_type", "trip.deleted").
		Str("trip_id", tripID).
		Str("reason", reason).
		Msg("Processing trip.deleted event")

	// Check idempotency
	processed, err := s.eventRepo.IsEventProcessed(ctx, eventID)
	if err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to check event idempotency")
		return fmt.Errorf("idempotency check failed: %w", err)
	}
	if processed {
		log.Info().Str("event_id", eventID).Msg("Event already processed, skipping")
		return nil
	}

	// Delete from MongoDB
	if err := s.tripRepo.DeleteByTripID(ctx, tripID); err != nil {
		if domain.IsNotFoundError(err) {
			// Permanent error - trip doesn't exist (already deleted or never existed)
			log.Warn().Str("trip_id", tripID).Msg("Trip not found in MongoDB, marking event as processed")
			processedEvent := &domain.ProcessedEvent{
				EventID:     eventID,
				EventType:   "trip.deleted",
				ProcessedAt: time.Now(),
				Result:      "skipped",
			}
			if markErr := s.eventRepo.MarkEventProcessed(ctx, processedEvent); markErr != nil {
				log.Error().Err(markErr).Str("event_id", eventID).Msg("Failed to mark event as processed")
			}
			return domain.ErrSearchTripNotFound
		}
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to delete trip from MongoDB")
		return fmt.Errorf("mongodb delete failed: %w", err)
	}

	log.Info().Str("trip_id", tripID).Msg("Trip deleted from MongoDB successfully")

	// Delete from Solr (optional - log error but continue)
	if s.solrClient != nil {
		if err := s.solrClient.Delete(ctx, tripID); err != nil {
			log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to delete trip from Solr (continuing)")
			// Don't return error - MongoDB is source of truth
		} else {
			log.Info().Str("trip_id", tripID).Msg("Trip deleted from Solr successfully")
		}
	}

	// Invalidate cache for this trip
	if s.cache != nil {
		cacheKey := fmt.Sprintf("trip:%s", tripID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Error().Err(err).Str("cache_key", cacheKey).Msg("Failed to invalidate cache (continuing)")
		} else {
			log.Info().Str("cache_key", cacheKey).Msg("Cache invalidated successfully")
		}
	}

	// Mark event as processed
	processedEvent := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   "trip.deleted",
		ProcessedAt: time.Now(),
		Result:      "success",
	}
	if err := s.eventRepo.MarkEventProcessed(ctx, processedEvent); err != nil {
		log.Error().Err(err).Str("event_id", eventID).Msg("Failed to mark event as processed")
		return fmt.Errorf("mark event processed failed: %w", err)
	}

	log.Info().
		Str("event_id", eventID).
		Str("trip_id", tripID).
		Msg("trip.deleted event processed successfully")

	return nil
}
