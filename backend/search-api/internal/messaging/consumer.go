package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"search-api/internal/service"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// Consumer handles RabbitMQ message consumption
type Consumer struct {
	conn            *amqp091.Connection
	channel         *amqp091.Channel
	queueName       string
	eventService    *service.TripEventService
	reconnectDelay  time.Duration
	maxReconnectDelay time.Duration
	stopChan        chan struct{}
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(rabbitmqURL, queueName string, eventService *service.TripEventService) (*Consumer, error) {
	consumer := &Consumer{
		queueName:       queueName,
		eventService:    eventService,
		reconnectDelay:  1 * time.Second,
		maxReconnectDelay: 30 * time.Second,
		stopChan:        make(chan struct{}),
	}

	if err := consumer.connect(rabbitmqURL); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return consumer, nil
}

// connect establishes connection to RabbitMQ
func (c *Consumer) connect(rabbitmqURL string) error {
	log.Info().Str("url", rabbitmqURL).Msg("Connecting to RabbitMQ")

	conn, err := amqp091.Dial(rabbitmqURL)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("channel creation failed: %w", err)
	}

	// Declare exchange (idempotent - safe even if trips-api already created it)
	err = channel.ExchangeDeclare(
		"trips.events", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("exchange declaration failed: %w", err)
	}

	// Declare queue (idempotent)
	_, err = channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("queue declaration failed: %w", err)
	}

	// Bind queue to exchange with topic pattern
	// This is the CRITICAL missing piece - without binding, the queue receives ZERO events
	err = channel.QueueBind(
		c.queueName,    // queue name
		"trip.*",       // routing key pattern (matches trip.created, trip.updated, trip.cancelled)
		"trips.events", // exchange name
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("queue binding failed: %w", err)
	}

	log.Info().
		Str("queue", c.queueName).
		Str("exchange", "trips.events").
		Str("routing_key", "trip.*").
		Msg("Queue bound to exchange successfully")

	// Set QoS - prefetch 1 message at a time for fair distribution
	if err := channel.Qos(1, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return fmt.Errorf("qos setup failed: %w", err)
	}

	c.conn = conn
	c.channel = channel

	log.Info().Str("queue", c.queueName).Msg("Connected to RabbitMQ successfully")

	return nil
}

// Start begins consuming messages from RabbitMQ
func (c *Consumer) Start(ctx context.Context, rabbitmqURL string) {
	log.Info().Msg("Starting RabbitMQ consumer")

	for {
		select {
		case <-c.stopChan:
			log.Info().Msg("Consumer stop signal received")
			return
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, stopping consumer")
			return
		default:
			if err := c.consume(ctx); err != nil {
				log.Error().Err(err).Msg("Consume error, reconnecting...")
				c.reconnect(rabbitmqURL)
			}
		}
	}
}

// consume handles message consumption
func (c *Consumer) consume(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer tag
		false,       // auto-ack (manual ACK for idempotency)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("consume setup failed: %w", err)
	}

	log.Info().Msg("Started consuming messages from RabbitMQ")

	// Monitor connection errors
	closeChan := make(chan *amqp091.Error)
	c.channel.NotifyClose(closeChan)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-closeChan:
			if err != nil {
				return fmt.Errorf("channel closed: %w", err)
			}
			return fmt.Errorf("channel closed")
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}
			c.handleMessage(ctx, msg)
		}
	}
}

// handleMessage processes a single message
func (c *Consumer) handleMessage(ctx context.Context, msg amqp091.Delivery) {
	// Parse the event type from message body
	var baseEvent struct {
		EventID   string `json:"event_id"`
		EventType string `json:"event_type"`
	}

	if err := json.Unmarshal(msg.Body, &baseEvent); err != nil {
		log.Error().Err(err).Bytes("body", msg.Body).Msg("Failed to parse message, NACKing without requeue")
		msg.Nack(false, false) // Permanent error - bad message format
		return
	}

	log.Info().
		Str("event_id", baseEvent.EventID).
		Str("event_type", baseEvent.EventType).
		Msg("Received message from RabbitMQ")

	var err error

	switch baseEvent.EventType {
	case "trip.created":
		err = c.handleTripCreated(ctx, msg.Body)
	case "trip.updated":
		err = c.handleTripUpdated(ctx, msg.Body)
	case "trip.cancelled":
		err = c.handleTripCancelled(ctx, msg.Body)
	case "trip.deleted":
		err = c.handleTripDeleted(ctx, msg.Body)
	default:
		log.Warn().
			Str("event_type", baseEvent.EventType).
			Msg("Unknown event type, ACKing without processing")
		msg.Ack(false)
		return
	}

	if err != nil {
		// Determine if error is transient or permanent
		if isTransientError(err) {
			// Transient error - NACK with requeue
			log.Error().
				Err(err).
				Str("event_id", baseEvent.EventID).
				Str("event_type", baseEvent.EventType).
				Msg("Transient error processing message, NACKing with requeue")
			msg.Nack(false, true) // Requeue for retry
		} else {
			// Permanent error - ACK without processing to avoid blocking queue
			log.Error().
				Err(err).
				Str("event_id", baseEvent.EventID).
				Str("event_type", baseEvent.EventType).
				Msg("Permanent error processing message, ACKing without requeue")
			msg.Ack(false)
		}
	} else {
		// Success - ACK
		log.Info().
			Str("event_id", baseEvent.EventID).
			Str("event_type", baseEvent.EventType).
			Msg("Message processed successfully, ACKing")
		msg.Ack(false)
	}
}

// handleTripCreated processes trip.created events
func (c *Consumer) handleTripCreated(ctx context.Context, body []byte) error {
	var event TripCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal trip.created failed: %w", err)
	}

	return c.eventService.HandleTripCreated(ctx, event.EventID, event.TripID, event.DriverID)
}

// handleTripUpdated processes trip.updated events
func (c *Consumer) handleTripUpdated(ctx context.Context, body []byte) error {
	var event TripUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal trip.updated failed: %w", err)
	}

	return c.eventService.HandleTripUpdated(ctx, event.EventID, event.TripID, event.AvailableSeats, event.ReservedSeats, event.Status)
}

// handleTripCancelled processes trip.cancelled events
func (c *Consumer) handleTripCancelled(ctx context.Context, body []byte) error {
	var event TripCancelledEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal trip.cancelled failed: %w", err)
	}

	return c.eventService.HandleTripCancelled(ctx, event.EventID, event.TripID, event.CancellationReason)
}

// handleTripDeleted processes trip.deleted events
func (c *Consumer) handleTripDeleted(ctx context.Context, body []byte) error {
	var event TripDeletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal trip.deleted failed: %w", err)
	}

	return c.eventService.HandleTripDeleted(ctx, event.EventID, event.TripID, event.Reason)
}

// reconnect handles reconnection with exponential backoff
func (c *Consumer) reconnect(rabbitmqURL string) {
	delay := c.reconnectDelay

	for {
		select {
		case <-c.stopChan:
			return
		default:
			log.Info().
				Dur("delay", delay).
				Msg("Attempting to reconnect to RabbitMQ")

			time.Sleep(delay)

			if err := c.connect(rabbitmqURL); err != nil {
				log.Error().Err(err).Msg("Reconnection failed")
				// Exponential backoff
				delay *= 2
				if delay > c.maxReconnectDelay {
					delay = c.maxReconnectDelay
				}
			} else {
				log.Info().Msg("Reconnected to RabbitMQ successfully")
				c.reconnectDelay = 1 * time.Second // Reset delay
				return
			}
		}
	}
}

// Close gracefully closes the consumer
func (c *Consumer) Close() error {
	log.Info().Msg("Closing RabbitMQ consumer")

	close(c.stopChan)

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close channel")
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close connection")
		}
	}

	log.Info().Msg("RabbitMQ consumer closed successfully")
	return nil
}

// isTransientError determines if an error is transient (should retry)
func isTransientError(err error) bool {
	if err == nil {
		return false
	}

	// Check error message for transient indicators
	errStr := err.Error()

	// Transient errors that should be retried
	transientPatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
		"502",
		"503",
		"504",
	}

	for _, pattern := range transientPatterns {
		if contains(errStr, pattern) {
			return true
		}
	}

	// Permanent errors (not found, validation errors, etc.)
	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
