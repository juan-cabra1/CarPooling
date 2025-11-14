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

// ============================================================================
// INCOMING EVENTS (Consumed from bookings-api)
// ============================================================================

// ReservationCreatedEvent representa un evento de reserva creada (incoming from bookings-api)
type ReservationCreatedEvent struct {
	EventID       string    `json:"event_id"`        // UUID v4 - CRÍTICO para idempotencia
	EventType     string    `json:"event_type"`      // "reservation.created"
	TripID        string    `json:"trip_id"`         // MongoDB ObjectID como string
	PassengerID   int64     `json:"passenger_id"`    // ID del pasajero
	SeatsReserved int       `json:"seats_reserved"`  // Número de asientos a reservar
	ReservationID string    `json:"reservation_id"`  // UUID de bookings-api
	Timestamp     time.Time `json:"timestamp"`       // Timestamp del evento
}

// ReservationCancelledEvent representa un evento de reserva cancelada (incoming from bookings-api)
type ReservationCancelledEvent struct {
	EventID       string    `json:"event_id"`        // UUID v4 - CRÍTICO para idempotencia
	EventType     string    `json:"event_type"`      // "reservation.cancelled"
	TripID        string    `json:"trip_id"`         // MongoDB ObjectID como string
	SeatsReleased int       `json:"seats_released"`  // Número de asientos a liberar
	ReservationID string    `json:"reservation_id"`  // UUID de bookings-api
	Timestamp     time.Time `json:"timestamp"`       // Timestamp del evento
}

// ============================================================================
// OUTGOING COMPENSATING EVENTS
// ============================================================================

// ReservationFailedEvent representa un evento de compensación cuando falla una reserva
type ReservationFailedEvent struct {
	EventID        string    `json:"event_id"`        // Nuevo UUID v4
	EventType      string    `json:"event_type"`      // "reservation.failed"
	ReservationID  string    `json:"reservation_id"`  // UUID de la reserva que falló
	TripID         string    `json:"trip_id"`         // MongoDB ObjectID como string
	Reason         string    `json:"reason"`          // "No seats available" | "Version conflict"
	AvailableSeats int       `json:"available_seats"` // Cantidad actual de asientos disponibles
	SourceService  string    `json:"source_service"`  // "trips-api"
	CorrelationID  string    `json:"correlation_id"`  // Para tracing de requests
	Timestamp      time.Time `json:"timestamp"`       // Timestamp del evento
}

// ReservationConfirmedEvent representa un evento de confirmación de reserva exitosa
type ReservationConfirmedEvent struct {
	EventID        string    `json:"event_id"`        // Nuevo UUID v4
	EventType      string    `json:"event_type"`      // "reservation.confirmed"
	ReservationID  string    `json:"reservation_id"`  // UUID de la reserva confirmada
	TripID         string    `json:"trip_id"`         // MongoDB ObjectID como string
	PassengerID    int64     `json:"passenger_id"`    // ID del pasajero
	DriverID       int64     `json:"driver_id"`       // ID del conductor del viaje
	SeatsReserved  int       `json:"seats_reserved"`  // Número de asientos reservados
	TotalPrice     float64   `json:"total_price"`     // Precio total de la reserva
	AvailableSeats int       `json:"available_seats"` // Asientos disponibles después de reserva
	SourceService  string    `json:"source_service"`  // "trips-api"
	CorrelationID  string    `json:"correlation_id"`  // Para tracing de requests
	Timestamp      time.Time `json:"timestamp"`       // Timestamp del evento
}
