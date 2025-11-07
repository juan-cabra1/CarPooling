package messaging

import "time"

// TripEvent representa el evento base para eventos de viajes
// Usado para trip.created y trip.updated
type TripEvent struct {
	EventID        string    `json:"event_id"`         // UUID v4 - CRÍTICO para idempotencia
	EventType      string    `json:"event_type"`       // trip.created, trip.updated, trip.cancelled
	TripID         string    `json:"trip_id"`          // MongoDB ObjectID como string
	DriverID       int64     `json:"driver_id"`        // ID del conductor
	Status         string    `json:"status"`           // Estado actual del viaje
	AvailableSeats int       `json:"available_seats"`  // Asientos disponibles
	ReservedSeats  int       `json:"reserved_seats"`   // Asientos reservados
	Timestamp      time.Time `json:"timestamp"`        // Timestamp del evento
	SourceService  string    `json:"source_service"`   // Siempre "trips-api"
	CorrelationID  string    `json:"correlation_id"`   // ID para tracing de requests
}

// TripCancelledEvent representa el evento de cancelación de viaje
// Extiende TripEvent con información adicional de cancelación
type TripCancelledEvent struct {
	TripEvent
	CancelledBy        int64  `json:"cancelled_by"`         // ID del usuario que canceló
	CancellationReason string `json:"cancellation_reason"`  // Razón de la cancelación
}
