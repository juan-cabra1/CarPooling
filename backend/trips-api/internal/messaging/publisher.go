package messaging

import (
	"context"
	"encoding/json"
	"time"
	"trips-api/internal/domain"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	// Exchange configuration
	exchangeName = "trips.events"
	exchangeType = "topic"

	// Routing keys
	routingKeyTripCreated   = "trip.created"
	routingKeyTripUpdated   = "trip.updated"
	routingKeyTripCancelled = "trip.cancelled"

	// Source service identifier
	sourceService = "trips-api"
)

// Publisher define la interfaz para publicar eventos de viajes a RabbitMQ
type Publisher interface {
	PublishTripCreated(ctx context.Context, trip *domain.Trip)
	PublishTripUpdated(ctx context.Context, trip *domain.Trip)
	PublishTripCancelled(ctx context.Context, trip *domain.Trip, cancelledBy int64, reason string)
	Close() error
}

type publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher crea una nueva instancia del publisher de RabbitMQ
// Establece conexión y declara el exchange necesario
func NewPublisher(rabbitURL string) (Publisher, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declarar exchange (idempotente - no falla si ya existe)
	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type: topic
		true,         // durable: sobrevive a reinicio del broker
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Info().
		Str("exchange", exchangeName).
		Str("type", exchangeType).
		Msg("RabbitMQ exchange declared successfully")

	return &publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

// PublishTripCreated publica un evento trip.created
func (p *publisher) PublishTripCreated(ctx context.Context, trip *domain.Trip) {
	event := TripEvent{
		EventID:        uuid.New().String(),
		EventType:      routingKeyTripCreated,
		TripID:         trip.ID.Hex(),
		DriverID:       trip.DriverID,
		Status:         trip.Status,
		AvailableSeats: trip.AvailableSeats,
		ReservedSeats:  trip.ReservedSeats,
		Timestamp:      time.Now(),
		SourceService:  sourceService,
		CorrelationID:  getCorrelationID(ctx),
	}

	p.publish(ctx, routingKeyTripCreated, event)
}

// PublishTripUpdated publica un evento trip.updated
func (p *publisher) PublishTripUpdated(ctx context.Context, trip *domain.Trip) {
	event := TripEvent{
		EventID:        uuid.New().String(),
		EventType:      routingKeyTripUpdated,
		TripID:         trip.ID.Hex(),
		DriverID:       trip.DriverID,
		Status:         trip.Status,
		AvailableSeats: trip.AvailableSeats,
		ReservedSeats:  trip.ReservedSeats,
		Timestamp:      time.Now(),
		SourceService:  sourceService,
		CorrelationID:  getCorrelationID(ctx),
	}

	p.publish(ctx, routingKeyTripUpdated, event)
}

// PublishTripCancelled publica un evento trip.cancelled con información adicional
func (p *publisher) PublishTripCancelled(ctx context.Context, trip *domain.Trip, cancelledBy int64, reason string) {
	event := TripCancelledEvent{
		TripEvent: TripEvent{
			EventID:        uuid.New().String(),
			EventType:      routingKeyTripCancelled,
			TripID:         trip.ID.Hex(),
			DriverID:       trip.DriverID,
			Status:         trip.Status,
			AvailableSeats: trip.AvailableSeats,
			ReservedSeats:  trip.ReservedSeats,
			Timestamp:      time.Now(),
			SourceService:  sourceService,
			CorrelationID:  getCorrelationID(ctx),
		},
		CancelledBy:        cancelledBy,
		CancellationReason: reason,
	}

	p.publish(ctx, routingKeyTripCancelled, event)
}

// publish es el método interno que serializa y publica eventos a RabbitMQ
// Implementa estrategia fire-and-forget: registra errores pero no los propaga
func (p *publisher) publish(ctx context.Context, routingKey string, event interface{}) {
	// Serializar evento a JSON
	body, err := json.Marshal(event)
	if err != nil {
		log.Error().
			Err(err).
			Str("routing_key", routingKey).
			Msg("Failed to marshal event to JSON")
		return
	}

	// Publicar mensaje con confirmación de contexto
	err = p.channel.PublishWithContext(
		ctx,
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 2 = persistent (sobrevive a reinicio)
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		// Fire-and-forget: registrar error pero no fallar la operación
		log.Error().
			Err(err).
			Str("routing_key", routingKey).
			Str("exchange", exchangeName).
			RawJSON("event", body).
			Msg("Failed to publish event to RabbitMQ")
		return
	}

	// Log exitoso
	log.Info().
		Str("routing_key", routingKey).
		Str("exchange", exchangeName).
		RawJSON("event", body).
		Msg("Event published successfully to RabbitMQ")
}

// Close cierra el canal y la conexión de RabbitMQ
func (p *publisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing RabbitMQ channel")
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing RabbitMQ connection")
			return err
		}
	}
	log.Info().Msg("RabbitMQ connection closed successfully")
	return nil
}

// getCorrelationID extrae o genera un correlation ID para tracing
func getCorrelationID(ctx context.Context) string {
	// Intentar extraer correlation ID del contexto
	if corrID, ok := ctx.Value("correlation_id").(string); ok && corrID != "" {
		return corrID
	}
	// Si no existe, generar uno nuevo
	return uuid.New().String()
}
