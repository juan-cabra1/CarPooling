package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"bookings-api/internal/dao"
)

// HandleTripCancelled processes trip.cancelled events
// Cancels all confirmed bookings for the cancelled trip
func (c *TripsConsumer) HandleTripCancelled(body []byte) error {
	var event TripCancelledEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().
			Err(err).
			Str("raw_body", string(body)).
			Msg("Failed to unmarshal trip.cancelled event")
		// Return nil to ACK - malformed JSON can't be reprocessed
		return nil
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("event_type", event.EventType).
		Str("trip_id", event.TripID).
		Str("correlation_id", event.CorrelationID).
		Str("cancellation_reason", event.CancellationReason).
		Int64("cancelled_by", event.CancelledBy).
		Msg("Processing trip.cancelled event")

	// Check idempotency - skip if already processed
	shouldProcess, err := c.idempotencyService.CheckAndMarkEvent(event.EventID, event.EventType)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Msg("Idempotency check failed")
		return fmt.Errorf("idempotency check failed: %w", err)
	}

	if !shouldProcess {
		log.Info().
			Str("event_id", event.EventID).
			Str("trip_id", event.TripID).
			Msg("Event already processed, skipping")
		return nil
	}

	// Find all confirmed bookings for this trip
	bookings, err := c.bookingRepo.FindByTripID(event.TripID)
	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", event.TripID).
			Msg("Failed to find bookings for trip")
		return fmt.Errorf("failed to find bookings: %w", err)
	}

	// Filter only confirmed bookings (others might already be cancelled/failed)
	confirmedBookings := make([]dao.Booking, 0)
	for _, booking := range bookings {
		if booking.Status == dao.BookingStatusConfirmed {
			confirmedBookings = append(confirmedBookings, booking)
		}
	}

	if len(confirmedBookings) == 0 {
		log.Info().
			Str("trip_id", event.TripID).
			Msg("No confirmed bookings found for cancelled trip")
		return nil
	}

	// Update all confirmed bookings to cancelled
	cancelledCount := 0
	failedCount := 0

	cancellationReason := fmt.Sprintf("Trip cancelled by driver: %s", event.CancellationReason)

	for _, booking := range confirmedBookings {
		err := c.bookingRepo.CancelBooking(booking.BookingUUID, cancellationReason)
		if err != nil {
			log.Error().
				Err(err).
				Str("booking_id", booking.BookingUUID).
				Str("trip_id", event.TripID).
				Msg("Failed to cancel booking")
			failedCount++
			continue
		}

		log.Info().
			Str("booking_id", booking.BookingUUID).
			Str("trip_id", event.TripID).
			Int64("passenger_id", booking.PassengerID).
			Msg("Booking cancelled due to trip cancellation")

		cancelledCount++
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("trip_id", event.TripID).
		Int("total_bookings", len(confirmedBookings)).
		Int("cancelled", cancelledCount).
		Int("failed", failedCount).
		Msg("Completed processing trip.cancelled event")

	// If some bookings failed to cancel, return error to NACK
	if failedCount > 0 {
		return fmt.Errorf("failed to cancel %d out of %d bookings", failedCount, len(confirmedBookings))
	}

	return nil
}

// HandleReservationFailed processes reservation.failed events
// Updates booking status from pending to failed
func (c *TripsConsumer) HandleReservationFailed(body []byte) error {
	var event ReservationFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().
			Err(err).
			Str("raw_body", string(body)).
			Msg("Failed to unmarshal reservation.failed event")
		// Return nil to ACK - malformed JSON can't be reprocessed
		return nil
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("event_type", event.EventType).
		Str("reservation_id", event.ReservationID).
		Str("trip_id", event.TripID).
		Str("correlation_id", event.CorrelationID).
		Str("reason", event.Reason).
		Msg("Processing reservation.failed event")

	// Check idempotency - skip if already processed
	shouldProcess, err := c.idempotencyService.CheckAndMarkEvent(event.EventID, event.EventType)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Msg("Idempotency check failed")
		return fmt.Errorf("idempotency check failed: %w", err)
	}

	if !shouldProcess {
		log.Info().
			Str("event_id", event.EventID).
			Str("reservation_id", event.ReservationID).
			Msg("Event already processed, skipping")
		return nil
	}

	// Find booking by reservation_id (booking_uuid)
	booking, err := c.bookingRepo.FindByID(event.ReservationID)
	if err != nil {
		// Idempotent behavior: if booking doesn't exist, ACK silently
		// This handles race conditions where event arrives before booking creation
		log.Warn().
			Err(err).
			Str("reservation_id", event.ReservationID).
			Str("trip_id", event.TripID).
			Msg("Booking not found for failed reservation, acknowledging")
		return nil
	}

	// Update booking status to failed
	err = c.bookingRepo.UpdateStatus(booking.BookingUUID, dao.BookingStatusFailed)
	if err != nil {
		log.Error().
			Err(err).
			Str("booking_id", booking.BookingUUID).
			Str("trip_id", event.TripID).
			Msg("Failed to update booking status to failed")
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("booking_id", booking.BookingUUID).
		Str("trip_id", event.TripID).
		Int64("passenger_id", booking.PassengerID).
		Str("reason", event.Reason).
		Msg("Booking marked as failed due to reservation failure")

	return nil
}

// HandleReservationConfirmed processes reservation.confirmed events
// Updates booking status from pending to confirmed and sets total price
func (c *TripsConsumer) HandleReservationConfirmed(body []byte) error {
	var event ReservationConfirmedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().
			Err(err).
			Str("raw_body", string(body)).
			Msg("Failed to unmarshal reservation.confirmed event")
		// Return nil to ACK - malformed JSON can't be reprocessed
		return nil
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("event_type", event.EventType).
		Str("reservation_id", event.ReservationID).
		Str("trip_id", event.TripID).
		Str("correlation_id", event.CorrelationID).
		Int("seats_reserved", event.SeatsReserved).
		Float64("total_price", event.TotalPrice).
		Msg("Processing reservation.confirmed event")

	// Check idempotency - skip if already processed
	shouldProcess, err := c.idempotencyService.CheckAndMarkEvent(event.EventID, event.EventType)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Msg("Idempotency check failed")
		return fmt.Errorf("idempotency check failed: %w", err)
	}

	if !shouldProcess {
		log.Info().
			Str("event_id", event.EventID).
			Str("reservation_id", event.ReservationID).
			Msg("Event already processed, skipping")
		return nil
	}

	// Find booking by reservation_id (booking_uuid)
	booking, err := c.bookingRepo.FindByID(event.ReservationID)
	if err != nil {
		// Idempotent behavior: if booking doesn't exist, ACK silently
		// This handles race conditions where event arrives before booking creation
		log.Warn().
			Err(err).
			Str("reservation_id", event.ReservationID).
			Str("trip_id", event.TripID).
			Msg("Booking not found for confirmed reservation, acknowledging")
		return nil
	}

	// Update booking status to confirmed, set total price, and store driver_id
	booking.Status = dao.BookingStatusConfirmed
	booking.TotalPrice = event.TotalPrice
	booking.DriverID = event.DriverID // Store driver for local authorization checks

	err = c.bookingRepo.Update(booking)
	if err != nil {
		log.Error().
			Err(err).
			Str("booking_id", booking.BookingUUID).
			Str("trip_id", event.TripID).
			Float64("total_price", event.TotalPrice).
			Msg("Failed to update booking status to confirmed")
		return fmt.Errorf("failed to update booking: %w", err)
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("booking_id", booking.BookingUUID).
		Str("trip_id", event.TripID).
		Int64("passenger_id", booking.PassengerID).
		Int("seats_reserved", event.SeatsReserved).
		Float64("total_price", event.TotalPrice).
		Msg("✅ Booking confirmed successfully with price")

	return nil
}
