package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

const (
	// Inbound exchange/queue configuration (from bookings-api)
	consumerExchange = "reservations.events" // Topic exchange
	consumerQueue    = "trips.reservations"  // Durable queue
	bindingKey       = "reservation.*"       // Matches reservation.created, reservation.cancelled

	// Consumer settings
	prefetchCount = 10                   // Process 10 messages concurrently
	consumerTag   = "trips-api-consumer" // Consumer identifier
)

// TripServiceInterface define los métodos necesarios del trip service para el consumer
// Esto evita import cycles entre messaging y service
type TripServiceInterface interface {
	ProcessReservationCreated(ctx context.Context, event ReservationCreatedEvent) error
	ProcessReservationCancelled(ctx context.Context, event ReservationCancelledEvent) error
}

// IdempotencyServiceInterface define los métodos necesarios del idempotency service
type IdempotencyServiceInterface interface {
	CheckAndMarkEvent(ctx context.Context, eventID, eventType string) (shouldProcess bool, err error)
}

// ReservationConsumer define la interfaz para consumir eventos de reservas
type ReservationConsumer interface {
	Start(ctx context.Context) error
	Close() error
}

type reservationConsumer struct {
	conn               *amqp.Connection
	channel            *amqp.Channel
	tripService        TripServiceInterface
	idempotencyService IdempotencyServiceInterface
	publisher          Publisher
}

// NewReservationConsumer crea una nueva instancia del consumer de reservas
func NewReservationConsumer(
	rabbitURL string,
	tripService TripServiceInterface,
	idempotencyService IdempotencyServiceInterface,
	publisher Publisher,
) (ReservationConsumer, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declarar exchange (idempotente - no falla si ya existe)
	err = ch.ExchangeDeclare(
		consumerExchange, // name
		exchangeType,     // type: topic
		true,             // durable: sobrevive a reinicio del broker
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declarar queue durable
	queue, err := ch.QueueDeclare(
		consumerQueue, // name
		true,          // durable: sobrevive a reinicio
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange con routing key pattern
	err = ch.QueueBind(
		queue.Name,       // queue name
		bindingKey,       // routing key pattern
		consumerExchange, // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// Configurar QoS (prefetch count)
	err = ch.Qos(
		prefetchCount, // prefetch count: procesar N mensajes concurrentemente
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Info().
		Str("exchange", consumerExchange).
		Str("queue", consumerQueue).
		Str("binding_key", bindingKey).
		Int("prefetch", prefetchCount).
		Msg("RabbitMQ consumer configured successfully")

	return &reservationConsumer{
		conn:               conn,
		channel:            ch,
		tripService:        tripService,
		idempotencyService: idempotencyService,
		publisher:          publisher,
	}, nil
}

// Start comienza a consumir mensajes (bloquea hasta que el contexto sea cancelado)
func (c *reservationConsumer) Start(ctx context.Context) error {
	// Iniciar consumo de mensajes
	deliveries, err := c.channel.Consume(
		consumerQueue, // queue
		consumerTag,   // consumer tag
		false,         // auto-ack: FALSE - manual ACK
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Info().
		Str("queue", consumerQueue).
		Str("consumer_tag", consumerTag).
		Msg("Started consuming reservation events")

	// Loop de procesamiento de mensajes
	for {
		select {
		case <-ctx.Done():
			// Contexto cancelado - graceful shutdown
			log.Info().Msg("Consumer context cancelled, stopping...")
			return nil

		case delivery, ok := <-deliveries:
			if !ok {
				// Canal cerrado - conexión perdida
				log.Warn().Msg("Deliveries channel closed, consumer stopping")
				return fmt.Errorf("deliveries channel closed")
			}

			// Procesar mensaje
			err := c.handleDelivery(ctx, delivery)
			if err != nil {
				// Error de sistema - NACK con requeue
				log.Error().
					Err(err).
					Str("message_id", delivery.MessageId).
					Msg("Failed to process message, will retry")
				delivery.Nack(false, true) // requeue=true
			} else {
				// Éxito o fallo manejado - ACK
				delivery.Ack(false)
			}
		}
	}
}

// handleDelivery procesa un solo mensaje
func (c *reservationConsumer) handleDelivery(ctx context.Context, delivery amqp.Delivery) error {
	// Log del mensaje recibido
	log.Debug().
		Str("routing_key", delivery.RoutingKey).
		Str("content_type", delivery.ContentType).
		Int("body_size", len(delivery.Body)).
		Msg("Received message")

	// Parsear evento base para obtener event_id y event_type
	var baseEvent struct {
		EventID   string `json:"event_id"`
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(delivery.Body, &baseEvent); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal base event")
		return nil // ACK - mensaje inválido no se puede procesar
	}

	// CRÍTICO: Verificar idempotencia PRIMERO
	shouldProcess, err := c.idempotencyService.CheckAndMarkEvent(ctx, baseEvent.EventID, baseEvent.EventType)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_id", baseEvent.EventID).
			Str("event_type", baseEvent.EventType).
			Msg("Idempotency check failed")
		return err // NACK - retry
	}

	if !shouldProcess {
		log.Info().
			Str("event_id", baseEvent.EventID).
			Str("event_type", baseEvent.EventType).
			Msg("Event already processed, skipping")
		return nil // ACK - skip processing
	}

	// Rutear al handler apropiado según el tipo de evento
	switch baseEvent.EventType {
	case "reservation.created":
		return c.handleReservationCreated(ctx, delivery.Body)

	case "reservation.cancelled":
		return c.handleReservationCancelled(ctx, delivery.Body)

	default:
		log.Warn().
			Str("event_type", baseEvent.EventType).
			Msg("Unknown event type, ignoring")
		return nil // ACK - tipo desconocido
	}
}

// handleReservationCreated procesa eventos de reserva creada
func (c *reservationConsumer) handleReservationCreated(ctx context.Context, body []byte) error {
	var event ReservationCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal reservation.created event")
		return nil // ACK - JSON inválido
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("trip_id", event.TripID).
		Str("reservation_id", event.ReservationID).
		Int("seats_reserved", event.SeatsReserved).
		Msg("Processing reservation.created event")

	// Delegar al servicio de negocio
	err := c.tripService.ProcessReservationCreated(ctx, event)
	if err != nil {
		// Error de sistema - NACK
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Msg("Failed to process reservation.created")
		return err
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("trip_id", event.TripID).
		Msg("Successfully processed reservation.created event")
	return nil
}

// handleReservationCancelled procesa eventos de reserva cancelada
func (c *reservationConsumer) handleReservationCancelled(ctx context.Context, body []byte) error {
	var event ReservationCancelledEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal reservation.cancelled event")
		return nil // ACK - JSON inválido
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("trip_id", event.TripID).
		Str("reservation_id", event.ReservationID).
		Int("seats_released", event.SeatsReleased).
		Msg("Processing reservation.cancelled event")

	// Delegar al servicio de negocio
	err := c.tripService.ProcessReservationCancelled(ctx, event)
	if err != nil {
		// Error de sistema - NACK
		log.Error().
			Err(err).
			Str("event_id", event.EventID).
			Msg("Failed to process reservation.cancelled")
		return err
	}

	log.Info().
		Str("event_id", event.EventID).
		Str("trip_id", event.TripID).
		Msg("Successfully processed reservation.cancelled event")
	return nil
}

// Close cierra el canal y la conexión de RabbitMQ
func (c *reservationConsumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing consumer channel")
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing consumer connection")
			return err
		}
	}
	log.Info().Msg("RabbitMQ consumer closed successfully")
	return nil
}
