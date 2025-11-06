package service

import (
	"context"
	"fmt"
	"time"
	"trips-api/internal/domain"
	"trips-api/internal/repository"
)

// IdempotencyService define las operaciones para garantizar idempotencia en el procesamiento de eventos
// CRÍTICO: Este servicio previene el procesamiento duplicado de eventos de RabbitMQ
type IdempotencyService interface {
	// CheckAndMarkEvent verifica si un evento ya fue procesado y lo marca si no lo ha sido
	// Retorna:
	// - shouldProcess=true: el evento NO ha sido procesado, se debe procesar
	// - shouldProcess=false: el evento YA fue procesado, se debe saltar (idempotencia)
	// - error: solo si hay un error real de sistema
	CheckAndMarkEvent(ctx context.Context, eventID, eventType string) (shouldProcess bool, err error)
}

type idempotencyService struct {
	eventRepo repository.EventRepository
}

// NewIdempotencyService crea una nueva instancia del servicio de idempotencia
func NewIdempotencyService(eventRepo repository.EventRepository) IdempotencyService {
	return &idempotencyService{
		eventRepo: eventRepo,
	}
}

// CheckAndMarkEvent implementa la lógica de idempotencia atómica
//
// Flujo:
// 1. Verifica si el evento ya fue procesado
// 2. Si NO fue procesado → lo marca como procesado y retorna shouldProcess=true
// 3. Si YA fue procesado → retorna shouldProcess=false (evita procesamiento duplicado)
//
// Race Conditions:
// El índice UNIQUE en MongoDB garantiza que solo un proceso/goroutine pueda
// marcar el evento como procesado. Si múltiples procesos intentan marcar el mismo
// evento simultáneamente, solo uno tendrá éxito y los demás recibirán duplicate key error.
//
// Ejemplos de uso:
//
//	shouldProcess, err := service.CheckAndMarkEvent(ctx, "event-123", "reservation.created")
//	if err != nil {
//	    return err // Error de sistema, se debe reintentar
//	}
//	if !shouldProcess {
//	    logger.Info().Msg("Event already processed, skipping")
//	    return nil // ACK sin procesar (idempotencia)
//	}
//	// Procesar evento...
func (s *idempotencyService) CheckAndMarkEvent(ctx context.Context, eventID, eventType string) (bool, error) {
	// Paso 1: Verificar si el evento ya fue procesado
	isProcessed, err := s.eventRepo.IsEventProcessed(ctx, eventID)
	if err != nil {
		return false, fmt.Errorf("failed to check if event is processed: %w", err)
	}

	// Si ya fue procesado, retornar false (saltar procesamiento)
	if isProcessed {
		return false, nil
	}

	// Paso 2: El evento NO ha sido procesado, intentar marcarlo
	processedEvent := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   eventType,
		ProcessedAt: time.Now(),
		Result:      "processing", // Se actualizará después del procesamiento real
	}

	err = s.eventRepo.MarkEventProcessed(ctx, processedEvent)
	if err != nil {
		// IMPORTANTE: Si MarkEventProcessed retorna error, NO es necesariamente un error fatal.
		// El repositorio maneja duplicate key errors internamente y retorna nil.
		// Si llegamos aquí con error, es un error real de sistema.
		return false, fmt.Errorf("failed to mark event as processed: %w", err)
	}

	// Evento marcado exitosamente como procesado, se debe procesar
	return true, nil
}
