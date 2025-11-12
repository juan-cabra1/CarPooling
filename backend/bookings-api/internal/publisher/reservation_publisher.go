package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bookings-api/internal/config"
	"bookings-api/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// ============================================================================
// RABBITMQ PUBLISHER FOR RESERVATION EVENTS
// ============================================================================
// This publisher is responsible for emitting reservation events from
// bookings-api to RabbitMQ for consumption by trips-api.
//
// Architecture:
//   bookings-api → RabbitMQ Exchange → trips-api Consumer
//
// Events Published:
//   - reservation.created  (when booking is created)
//   - reservation.cancelled (when booking is cancelled)
//
// Exchange Configuration:
//   - Name: "bookings.events"
//   - Type: "topic"
//   - Durable: true (survives RabbitMQ restart)
//
// Message Persistence:
//   - Delivery Mode: Persistent (messages survive broker restart)
//
// Routing Keys:
//   - "reservation.created" for new bookings
//   - "reservation.cancelled" for cancellations
//
// Idempotency:
// Each event includes a unique UUID (event_id) to allow consumers
// to detect and skip duplicate messages.
// ============================================================================

const (
	// ExchangeName is the RabbitMQ exchange for bookings events
	ExchangeName = "bookings.events"

	// ExchangeType defines the exchange type (topic allows routing key patterns)
	ExchangeType = "topic"

	// RoutingKeyReservationCreated is used when publishing reservation.created events
	RoutingKeyReservationCreated = "reservation.created"

	// RoutingKeyReservationCancelled is used when publishing reservation.cancelled events
	RoutingKeyReservationCancelled = "reservation.cancelled"
)

// ============================================================================
// RESERVATION PUBLISHER INTERFACE
// ============================================================================

// Publisher defines the interface for publishing reservation events to RabbitMQ
// This interface allows for easy mocking in tests without requiring actual RabbitMQ connection
type Publisher interface {
	// PublishReservationCreated publishes a reservation.created event
	PublishReservationCreated(tripID string, seatsReserved int, reservationID string) error

	// PublishReservationCancelled publishes a reservation.cancelled event
	PublishReservationCancelled(tripID string, seatsReleased int, reservationID string) error

	// Close closes the RabbitMQ connection and channel
	Close() error
}

// ============================================================================
// RESERVATION PUBLISHER STRUCT
// ============================================================================

// ReservationPublisher manages RabbitMQ connection and publishes reservation events
//
// This publisher:
//   - Maintains a persistent connection to RabbitMQ
//   - Declares the bookings.events exchange on initialization
//   - Provides methods to publish reservation events
//   - Handles errors gracefully without panicking
//   - Logs all operations with structured logging (zerolog)
//
// Usage:
//
//	publisher, err := NewReservationPublisher(cfg, logger)
//	if err != nil {
//	    log.Fatal().Err(err).Msg("Failed to create publisher")
//	}
//	defer publisher.Close()
//
//	err = publisher.PublishReservationCreated("trip-123", 2, "booking-456")
//	if err != nil {
//	    log.Error().Err(err).Msg("Failed to publish event")
//	}
//
// Thread Safety:
// This publisher is NOT thread-safe. If using from multiple goroutines,
// consider creating a channel pool or using mutex for synchronization.
type ReservationPublisher struct {
	// conn is the RabbitMQ connection
	conn *amqp.Connection

	// channel is the RabbitMQ channel for publishing messages
	channel *amqp.Channel

	// exchangeName is the name of the exchange to publish to
	exchangeName string

	// logger is the structured logger (zerolog)
	logger zerolog.Logger
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================

// NewReservationPublisher creates a new RabbitMQ publisher for reservation events
//
// This constructor:
//   1. Connects to RabbitMQ using the URL from configuration
//   2. Opens a channel for publishing
//   3. Declares the "bookings.events" topic exchange (durable)
//   4. Returns the publisher instance ready for use
//
// Connection URL Format:
//
//	amqp://username:password@host:port/vhost
//	Example: amqp://guest:guest@localhost:5672/
//
// Exchange Configuration:
//   - Name: "bookings.events"
//   - Type: "topic" (allows routing key patterns)
//   - Durable: true (survives RabbitMQ restart)
//   - Auto-delete: false (persists even if no consumers)
//
// Error Handling:
//   - Connection failures → return error with sanitized URL
//   - Channel creation failures → return error
//   - Exchange declaration failures → return error
//   - All errors are wrapped with context for debugging
//
// Parameters:
//   - cfg: Application configuration containing RabbitMQURL
//   - logger: Structured logger for operation logging
//
// Returns:
//   - *ReservationPublisher: Ready-to-use publisher instance
//   - error: Non-nil if connection or setup fails
//
// Example:
//
//	cfg := &config.Config{RabbitMQURL: "amqp://guest:guest@localhost:5672/"}
//	logger := zerolog.New(os.Stderr)
//	publisher, err := NewReservationPublisher(cfg, logger)
//	if err != nil {
//	    log.Fatal().Err(err).Msg("Failed to initialize publisher")
//	}
//	defer publisher.Close()
func NewReservationPublisher(cfg *config.Config, logger zerolog.Logger) (*ReservationPublisher, error) {
	// Validate configuration
	if cfg.RabbitMQURL == "" {
		return nil, fmt.Errorf("RabbitMQ URL is required but not provided in configuration")
	}

	// Log connection attempt (sanitize credentials for security)
	sanitizedURL := sanitizeRabbitMQURL(cfg.RabbitMQURL)
	logger.Info().
		Str("url", sanitizedURL).
		Msg("🔌 Connecting to RabbitMQ for event publishing...")

	// ========================================================================
	// STEP 1: Connect to RabbitMQ
	// ========================================================================
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ at %s: %w", sanitizedURL, err)
	}

	logger.Info().Msg("✅ RabbitMQ connection established")

	// ========================================================================
	// STEP 2: Open a channel
	// ========================================================================
	// Channels are lightweight and allow concurrent operations
	channel, err := conn.Channel()
	if err != nil {
		conn.Close() // Clean up connection on failure
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	logger.Info().Msg("✅ RabbitMQ channel opened")

	// ========================================================================
	// STEP 3: Declare exchange
	// ========================================================================
	// This ensures the exchange exists before publishing
	// If exchange already exists with same config, this is idempotent
	err = channel.ExchangeDeclare(
		ExchangeName, // name
		ExchangeType, // type (topic)
		true,         // durable (survives RabbitMQ restart)
		false,        // auto-deleted (don't delete when unused)
		false,        // internal (no, can be published to directly)
		false,        // no-wait (wait for server confirmation)
		nil,          // arguments
	)
	if err != nil {
		channel.Close() // Clean up channel on failure
		conn.Close()    // Clean up connection on failure
		return nil, fmt.Errorf("failed to declare exchange '%s': %w", ExchangeName, err)
	}

	logger.Info().
		Str("exchange", ExchangeName).
		Str("type", ExchangeType).
		Msg("✅ Exchange declared successfully")

	// ========================================================================
	// STEP 4: Return publisher instance
	// ========================================================================
	return &ReservationPublisher{
		conn:         conn,
		channel:      channel,
		exchangeName: ExchangeName,
		logger:       logger,
	}, nil
}

// ============================================================================
// PUBLISH METHODS
// ============================================================================

// PublishReservationCreated publishes a reservation.created event to RabbitMQ
//
// This method:
//   1. Creates a ReservationCreatedEvent with unique UUID and timestamp
//   2. Marshals the event to JSON
//   3. Publishes to "bookings.events" exchange with routing key "reservation.created"
//   4. Uses persistent delivery mode (survives broker restart)
//   5. Logs the published event with structured fields
//
// Event Flow:
//   bookings-api (this method) → RabbitMQ Exchange → trips-api Consumer
//
// trips-api will receive this event and:
//   - Decrement available_seats by seatsReserved
//   - Track the reservation
//   - Use event_id for idempotency
//
// Parameters:
//   - tripID: MongoDB ObjectID of the trip (string)
//   - seatsReserved: Number of seats reserved (must be > 0)
//   - reservationID: Booking UUID from bookings table
//
// Returns:
//   - error: Non-nil if marshaling or publishing fails
//
// Error Handling:
//   - JSON marshal failure → return error
//   - RabbitMQ publish failure → return error
//   - Logs all errors with context
//
// Example:
//
//	err := publisher.PublishReservationCreated("trip-123", 2, "booking-456")
//	if err != nil {
//	    log.Error().Err(err).Msg("Failed to publish reservation.created event")
//	}
//
// Idempotency:
// Each event gets a unique event_id (UUID v4). If trips-api receives
// the same event_id twice, it will skip processing.
func (p *ReservationPublisher) PublishReservationCreated(tripID string, passengerID int64, seatsReserved int, reservationID string) error {
	// ========================================================================
	// STEP 1: Create event structure
	// ========================================================================
	event := events.ReservationCreatedEvent{
		BaseEvent:     events.NewBaseEvent(events.EventTypeReservationCreated),
		TripID:        tripID,
		PassengerID:   passengerID,
		SeatsReserved: seatsReserved,
		ReservationID: reservationID,
	}

	// ========================================================================
	// STEP 2: Marshal to JSON
	// ========================================================================
	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("event_type", events.EventTypeReservationCreated).
			Str("trip_id", tripID).
			Msg("❌ Failed to marshal reservation.created event")
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// ========================================================================
	// STEP 3: Publish to RabbitMQ
	// ========================================================================
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		p.exchangeName,                // exchange
		RoutingKeyReservationCreated,  // routing key
		false,                         // mandatory (don't return if no queue is bound)
		false,                         // immediate (don't wait for consumer confirmation)
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 2 = persistent (survives broker restart)
			Timestamp:    event.Timestamp,
			MessageId:    event.EventID, // Use event_id as message_id for tracing
		},
	)

	if err != nil {
		p.logger.Error().
			Err(err).
			Str("event_id", event.EventID).
			Str("event_type", events.EventTypeReservationCreated).
			Str("trip_id", tripID).
			Int("seats_reserved", seatsReserved).
			Str("reservation_id", reservationID).
			Msg("❌ Failed to publish reservation.created event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// ========================================================================
	// STEP 4: Log success
	// ========================================================================
	p.logger.Info().
		Str("event_id", event.EventID).
		Str("event_type", events.EventTypeReservationCreated).
		Str("trip_id", tripID).
		Int("seats_reserved", seatsReserved).
		Str("reservation_id", reservationID).
		Msg("✅ Published reservation.created event")

	return nil
}

// PublishReservationCancelled publishes a reservation.cancelled event to RabbitMQ
//
// This method:
//   1. Creates a ReservationCancelledEvent with unique UUID and timestamp
//   2. Marshals the event to JSON
//   3. Publishes to "bookings.events" exchange with routing key "reservation.cancelled"
//   4. Uses persistent delivery mode (survives broker restart)
//   5. Logs the published event with structured fields
//
// Event Flow:
//   bookings-api (this method) → RabbitMQ Exchange → trips-api Consumer
//
// trips-api will receive this event and:
//   - Increment available_seats by seatsReleased
//   - Remove the reservation tracking
//   - Use event_id for idempotency
//
// Parameters:
//   - tripID: MongoDB ObjectID of the trip (string)
//   - seatsReleased: Number of seats being released (must be > 0)
//   - reservationID: Booking UUID from bookings table
//
// Returns:
//   - error: Non-nil if marshaling or publishing fails
//
// Error Handling:
//   - JSON marshal failure → return error
//   - RabbitMQ publish failure → return error
//   - Logs all errors with context
//
// Example:
//
//	err := publisher.PublishReservationCancelled("trip-123", 2, "booking-456")
//	if err != nil {
//	    log.Error().Err(err).Msg("Failed to publish reservation.cancelled event")
//	}
//
// Idempotency:
// Each event gets a unique event_id (UUID v4). If trips-api receives
// the same event_id twice, it will skip processing.
func (p *ReservationPublisher) PublishReservationCancelled(tripID string, seatsReleased int, reservationID string) error {
	// ========================================================================
	// STEP 1: Create event structure
	// ========================================================================
	event := events.ReservationCancelledEvent{
		BaseEvent:     events.NewBaseEvent(events.EventTypeReservationCancelled),
		TripID:        tripID,
		SeatsReleased: seatsReleased,
		ReservationID: reservationID,
	}

	// ========================================================================
	// STEP 2: Marshal to JSON
	// ========================================================================
	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("event_type", events.EventTypeReservationCancelled).
			Str("trip_id", tripID).
			Msg("❌ Failed to marshal reservation.cancelled event")
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// ========================================================================
	// STEP 3: Publish to RabbitMQ
	// ========================================================================
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		p.exchangeName,                 // exchange
		RoutingKeyReservationCancelled, // routing key
		false,                          // mandatory
		false,                          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 2 = persistent
			Timestamp:    event.Timestamp,
			MessageId:    event.EventID, // Use event_id as message_id for tracing
		},
	)

	if err != nil {
		p.logger.Error().
			Err(err).
			Str("event_id", event.EventID).
			Str("event_type", events.EventTypeReservationCancelled).
			Str("trip_id", tripID).
			Int("seats_released", seatsReleased).
			Str("reservation_id", reservationID).
			Msg("❌ Failed to publish reservation.cancelled event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// ========================================================================
	// STEP 4: Log success
	// ========================================================================
	p.logger.Info().
		Str("event_id", event.EventID).
		Str("event_type", events.EventTypeReservationCancelled).
		Str("trip_id", tripID).
		Int("seats_released", seatsReleased).
		Str("reservation_id", reservationID).
		Msg("✅ Published reservation.cancelled event")

	return nil
}

// ============================================================================
// CONNECTION MANAGEMENT
// ============================================================================

// Close gracefully closes the RabbitMQ channel and connection
//
// This method should be called during application shutdown to:
//   - Close the publishing channel
//   - Close the RabbitMQ connection
//   - Release resources
//
// Usage:
//
//	publisher, _ := NewReservationPublisher(cfg, logger)
//	defer publisher.Close()
//
// Or in graceful shutdown:
//
//	// On shutdown signal
//	if err := publisher.Close(); err != nil {
//	    log.Error().Err(err).Msg("Error closing publisher")
//	}
//
// Error Handling:
//   - Channel close errors are logged but don't stop connection closure
//   - Connection close errors are returned
//
// Returns:
//   - error: Non-nil if connection closure fails
func (p *ReservationPublisher) Close() error {
	p.logger.Info().Msg("🔌 Closing RabbitMQ publisher...")

	// Close channel first
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			p.logger.Error().
				Err(err).
				Msg("⚠️  Error closing RabbitMQ channel")
			// Continue to close connection even if channel close fails
		} else {
			p.logger.Info().Msg("✅ RabbitMQ channel closed")
		}
	}

	// Close connection
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			p.logger.Error().
				Err(err).
				Msg("❌ Error closing RabbitMQ connection")
			return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
		}
		p.logger.Info().Msg("✅ RabbitMQ connection closed")
	}

	p.logger.Info().Msg("✅ Publisher closed successfully")
	return nil
}

// ============================================================================
// HELPER FUNCTIONS (Private)
// ============================================================================

// sanitizeRabbitMQURL removes the password from RabbitMQ URL for safe logging
//
// Example:
//
//	Input:  "amqp://guest:secret@localhost:5672/"
//	Output: "amqp://guest:***@localhost:5672/"
//
// This prevents password leaks in logs while still showing connection details.
func sanitizeRabbitMQURL(url string) string {
	// Find password section (between : and @)
	// Example URL: amqp://username:password@host:port/vhost
	if idx := strings.Index(url, "://"); idx != -1 {
		// Skip protocol
		rest := url[idx+3:]
		if idx2 := strings.Index(rest, ":"); idx2 != -1 {
			// Found username:password section
			if idx3 := strings.Index(rest[idx2:], "@"); idx3 != -1 {
				// Replace password with ***
				return url[:idx+3+idx2+1] + "***" + rest[idx2+idx3:]
			}
		}
	}
	// If format doesn't match, return as-is (no password to hide)
	return url
}
