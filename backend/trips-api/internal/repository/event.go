package repository

import (
	"context"
	"fmt"
	"time"
	"trips-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventRepository define las operaciones para manejar eventos procesados (idempotencia)
type EventRepository interface {
	IsEventProcessed(ctx context.Context, eventID string) (bool, error)
	MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error
}

type eventRepository struct {
	collection *mongo.Collection
}

// NewEventRepository crea una nueva instancia del repositorio de eventos
func NewEventRepository(db *mongo.Database) EventRepository {
	return &eventRepository{
		collection: db.Collection("processed_events"),
	}
}

// IsEventProcessed verifica si un evento ya fue procesado
func (r *eventRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	filter := bson.M{"event_id": eventID}

	var event domain.ProcessedEvent
	err := r.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// El evento NO ha sido procesado
			return false, nil
		}
		return false, fmt.Errorf("failed to check event: %w", err)
	}

	// El evento ya fue procesado
	return true, nil
}

// MarkEventProcessed marca un evento como procesado
// CRÍTICO: Este método debe manejar duplicate key errors (E11000) de MongoDB
// Si el evento ya existe (índice UNIQUE en event_id), retorna nil sin error
func (r *eventRepository) MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Establecer timestamp si no existe
	if event.ProcessedAt.IsZero() {
		event.ProcessedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		// Verificar si es un error de duplicate key (E11000)
		if mongo.IsDuplicateKeyError(err) {
			// El evento ya fue procesado por otro proceso/goroutine
			// Esto es esperado en escenarios de RabbitMQ con reintentos
			// Retornar nil para indicar éxito (idempotencia)
			return nil
		}
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}

	return nil
}
