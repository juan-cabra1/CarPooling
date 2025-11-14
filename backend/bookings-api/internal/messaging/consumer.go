package messaging

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"

	"bookings-api/internal/repository"
	"bookings-api/internal/service"
)

const (
	exchangeName        = "trips.events"
	exchangeType        = "topic"
	queueName           = "bookings.trip-events"
	prefetchCount       = 10
	routingKeyCancelled = "trip.cancelled"
	routingKeyFailed    = "reservation.failed"
	routingKeyConfirmed = "reservation.confirmed"
)

// TripsConsumer handles RabbitMQ messages from trips-api
type TripsConsumer struct {
	conn               *amqp.Connection
	channel            *amqp.Channel
	bookingRepo        repository.BookingRepository
	idempotencyService service.IdempotencyService
}

// NewTripsConsumer creates a new RabbitMQ consumer for trips events
func NewTripsConsumer(
	rabbitMQURL string,
	bookingRepo repository.BookingRepository,
	idempotencyService service.IdempotencyService,
) (*TripsConsumer, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange (idempotent)
	err = channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue (idempotent)
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange for trip.cancelled events
	err = channel.QueueBind(
		queue.Name,          // queue name
		routingKeyCancelled, // routing key
		exchangeName,        // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue for trip.cancelled: %w", err)
	}

	// Bind queue to exchange for reservation.failed events
	err = channel.QueueBind(
		queue.Name,       // queue name
		routingKeyFailed, // routing key
		exchangeName,     // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue for reservation.failed: %w", err)
	}

	// Bind queue to exchange for reservation.confirmed events
	err = channel.QueueBind(
		queue.Name,          // queue name
		routingKeyConfirmed, // routing key
		exchangeName,        // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue for reservation.confirmed: %w", err)
	}

	// Set QoS prefetch count
	err = channel.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Info().
		Str("exchange", exchangeName).
		Str("queue", queueName).
		Int("prefetch", prefetchCount).
		Msg("RabbitMQ consumer initialized successfully")

	return &TripsConsumer{
		conn:               conn,
		channel:            channel,
		bookingRepo:        bookingRepo,
		idempotencyService: idempotencyService,
	}, nil
}

// Start begins consuming messages from RabbitMQ
// Blocks until context is cancelled or an error occurs
func (c *TripsConsumer) Start(ctx context.Context) error {
	// Start consuming messages
	messages, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer tag (auto-generated)
		false,     // auto-ack (manual ACK for reliability)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Info().
		Str("queue", queueName).
		Msg("Started consuming messages from RabbitMQ")

	// Process messages until context is cancelled
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Consumer context cancelled, shutting down")
			return nil

		case msg, ok := <-messages:
			if !ok {
				log.Warn().Msg("Message channel closed, stopping consumer")
				return fmt.Errorf("message channel closed")
			}

			// Process message
			c.handleMessage(msg)
		}
	}
}

// handleMessage processes a single RabbitMQ message
func (c *TripsConsumer) handleMessage(msg amqp.Delivery) {
	log.Debug().
		Str("routing_key", msg.RoutingKey).
		Str("correlation_id", msg.CorrelationId).
		Msg("Received message")

	// Route to appropriate handler based on routing key
	var err error
	switch msg.RoutingKey {
	case routingKeyCancelled:
		err = c.HandleTripCancelled(msg.Body)
	case routingKeyFailed:
		err = c.HandleReservationFailed(msg.Body)
	case routingKeyConfirmed:
		err = c.HandleReservationConfirmed(msg.Body)
	default:
		log.Warn().
			Str("routing_key", msg.RoutingKey).
			Msg("Unknown routing key, acknowledging message")
		msg.Ack(false) // ACK unknown messages to avoid blocking queue
		return
	}

	// Handle processing result
	if err != nil {
		log.Error().
			Err(err).
			Str("routing_key", msg.RoutingKey).
			Str("correlation_id", msg.CorrelationId).
			Msg("Failed to process message, negative acknowledging")

		// NACK with requeue for system errors
		// This allows retry in case of temporary failures (DB connection, etc.)
		msg.Nack(false, true)
		return
	}

	// ACK successful processing
	if err := msg.Ack(false); err != nil {
		log.Error().
			Err(err).
			Str("routing_key", msg.RoutingKey).
			Msg("Failed to acknowledge message")
	}
}

// Close gracefully shuts down the consumer
func (c *TripsConsumer) Close() error {
	log.Info().Msg("Closing RabbitMQ consumer")

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing channel")
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing connection")
			return err
		}
	}

	return nil
}
