package service

import (
	"context"
	"fmt"
	"time"
	"trips-api/internal/clients"
	"trips-api/internal/domain"
	"trips-api/internal/messaging"
	"trips-api/internal/repository"

	"github.com/rs/zerolog/log"
)

// TripService define las operaciones de lógica de negocio para viajes
type TripService interface {
	// CreateTrip crea un nuevo viaje con validaciones de negocio
	CreateTrip(ctx context.Context, driverID int64, request domain.CreateTripRequest) (*domain.Trip, error)

	// GetTrip obtiene un viaje por su ID
	GetTrip(ctx context.Context, tripID string) (*domain.Trip, error)

	// ListTrips lista viajes con filtros y paginación
	ListTrips(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error)

	// UpdateTrip actualiza un viaje existente (solo el dueño)
	UpdateTrip(ctx context.Context, tripID string, userID int64, request domain.UpdateTripRequest) (*domain.Trip, error)

	// DeleteTrip elimina un viaje (solo el dueño)
	DeleteTrip(ctx context.Context, tripID string, userID int64) error

	// CancelTrip cancela un viaje (solo el dueño)
	CancelTrip(ctx context.Context, tripID string, userID int64, request domain.CancelTripRequest) error
}

type tripService struct {
	tripRepo           repository.TripRepository
	idempotencyService IdempotencyService
	usersClient        clients.UsersClient
	publisher          messaging.Publisher
}

// NewTripService crea una nueva instancia del servicio de viajes
func NewTripService(
	tripRepo repository.TripRepository,
	idempotencyService IdempotencyService,
	usersClient clients.UsersClient,
	publisher messaging.Publisher,
) TripService {
	return &tripService{
		tripRepo:           tripRepo,
		idempotencyService: idempotencyService,
		usersClient:        usersClient,
		publisher:          publisher,
	}
}

// CreateTrip implementa la creación de viajes con todas las validaciones de negocio
//
// Validaciones:
// - departure_datetime debe ser en el futuro
// - total_seats debe estar entre 1-8
// - driver_id debe existir (llamada a users-api)
//
// Valores iniciales:
// - available_seats = total_seats
// - reserved_seats = 0
// - status = "published"
// - availability_version = 1
//
// Ejemplo de uso:
//
//	trip, err := service.CreateTrip(ctx, userID, createRequest)
//	if err != nil {
//	    if errors.Is(err, domain.ErrPastDeparture) {
//	        return c.JSON(400, gin.H{"error": "departure must be in future"})
//	    }
//	    return c.JSON(500, gin.H{"error": err.Error()})
//	}
func (s *tripService) CreateTrip(ctx context.Context, driverID int64, request domain.CreateTripRequest) (*domain.Trip, error) {
	// Validación 1: Parsear fecha de salida
	departureTime, err := time.Parse(time.RFC3339, request.DepartureDatetime)
	if err != nil {
		return nil, fmt.Errorf("invalid departure_datetime format: %w", err)
	}

	// Validación 2: Parsear fecha de llegada estimada
	arrivalTime, err := time.Parse(time.RFC3339, request.EstimatedArrivalDatetime)
	if err != nil {
		return nil, fmt.Errorf("invalid estimated_arrival_datetime format: %w", err)
	}

	// Validación 3: La salida debe ser en el futuro
	if departureTime.Before(time.Now()) {
		return nil, domain.ErrPastDeparture
	}

	// Validación 4: La llegada debe ser después de la salida
	if arrivalTime.Before(departureTime) {
		return nil, fmt.Errorf("arrival time must be after departure time")
	}

	// Validación 5: total_seats debe estar entre 1-8 (validado por binding, pero verificamos)
	if request.TotalSeats < 1 || request.TotalSeats > 8 {
		return nil, fmt.Errorf("total_seats must be between 1 and 8")
	}

	// Validación 6: Verificar que el driver existe en users-api
	_, err = s.usersClient.GetUser(ctx, driverID)
	if err != nil {
		// Si es ErrDriverNotFound, mantener ese error específico
		return nil, fmt.Errorf("failed to validate driver: %w", err)
	}

	// Construir el trip con valores iniciales
	trip := &domain.Trip{
		DriverID:                 driverID,
		Origin:                   request.Origin,
		Destination:              request.Destination,
		DepartureDatetime:        departureTime,
		EstimatedArrivalDatetime: arrivalTime,
		PricePerSeat:             request.PricePerSeat,
		TotalSeats:               request.TotalSeats,
		Car:                      request.Car,
		Preferences:              request.Preferences,
		Description:              request.Description,

		// Valores iniciales CRÍTICOS
		AvailableSeats:      request.TotalSeats, // Todos los asientos disponibles inicialmente
		ReservedSeats:       0,                  // Sin reservas inicialmente
		Status:              "published",        // Estado inicial
		AvailabilityVersion: 1,                  // Versión inicial para optimistic locking
	}

	// Crear el trip en la base de datos
	if err := s.tripRepo.Create(ctx, trip); err != nil {
		log.Error().Err(err).Int64("driver_id", driverID).Msg("Failed to create trip")
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	log.Info().Str("trip_id", trip.ID.Hex()).Int64("driver_id", driverID).Msg("Trip created")

	// Publicar evento trip.created (fire-and-forget)
	s.publisher.PublishTripCreated(ctx, trip)

	return trip, nil
}

// GetTrip obtiene un viaje por su ID
func (s *tripService) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	return trip, nil
}

// ListTrips lista viajes con filtros y paginación
func (s *tripService) ListTrips(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	trips, total, err := s.tripRepo.FindAll(ctx, filters, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list trips: %w", err)
	}

	return trips, total, nil
}

// UpdateTrip actualiza un viaje existente con validaciones de negocio
//
// Validaciones:
// - Solo el dueño puede actualizar (userID == driver_id)
// - No se puede actualizar si reserved_seats > 0
// - No se puede cambiar total_seats a menos que reserved_seats
// - Las fechas deben ser válidas si se proporcionan
func (s *tripService) UpdateTrip(ctx context.Context, tripID string, userID int64, request domain.UpdateTripRequest) (*domain.Trip, error) {
	// Obtener el trip actual
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return nil, err
	}

	// Validación 1: Solo el dueño puede actualizar
	if trip.DriverID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Validación 2: No se puede actualizar si hay reservas
	if trip.ReservedSeats > 0 {
		return nil, domain.ErrHasReservations
	}

	// Aplicar actualizaciones opcionales
	if request.Origin != nil {
		trip.Origin = *request.Origin
	}

	if request.Destination != nil {
		trip.Destination = *request.Destination
	}

	if request.DepartureDatetime != nil {
		departureTime, err := time.Parse(time.RFC3339, *request.DepartureDatetime)
		if err != nil {
			return nil, fmt.Errorf("invalid departure_datetime format: %w", err)
		}
		if departureTime.Before(time.Now()) {
			return nil, domain.ErrPastDeparture
		}
		trip.DepartureDatetime = departureTime
	}

	if request.EstimatedArrivalDatetime != nil {
		arrivalTime, err := time.Parse(time.RFC3339, *request.EstimatedArrivalDatetime)
		if err != nil {
			return nil, fmt.Errorf("invalid estimated_arrival_datetime format: %w", err)
		}
		if arrivalTime.Before(trip.DepartureDatetime) {
			return nil, fmt.Errorf("arrival time must be after departure time")
		}
		trip.EstimatedArrivalDatetime = arrivalTime
	}

	if request.PricePerSeat != nil {
		if *request.PricePerSeat < 0 {
			return nil, fmt.Errorf("price_per_seat must be non-negative")
		}
		trip.PricePerSeat = *request.PricePerSeat
	}

	if request.TotalSeats != nil {
		// Validación 3: No se puede reducir total_seats por debajo de reserved_seats
		if *request.TotalSeats < trip.ReservedSeats {
			return nil, fmt.Errorf("cannot set total_seats below reserved_seats (%d)", trip.ReservedSeats)
		}
		if *request.TotalSeats < 1 || *request.TotalSeats > 8 {
			return nil, fmt.Errorf("total_seats must be between 1 and 8")
		}

		// Recalcular available_seats
		oldTotalSeats := trip.TotalSeats
		trip.TotalSeats = *request.TotalSeats
		trip.AvailableSeats = trip.AvailableSeats + (trip.TotalSeats - oldTotalSeats)
	}

	if request.Car != nil {
		trip.Car = *request.Car
	}

	if request.Preferences != nil {
		trip.Preferences = *request.Preferences
	}

	if request.Description != nil {
		trip.Description = *request.Description
	}

	// Actualizar en la base de datos
	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to update trip")
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	log.Info().Str("trip_id", tripID).Int64("user_id", userID).Msg("Trip updated")

	// Publicar evento trip.updated (fire-and-forget)
	s.publisher.PublishTripUpdated(ctx, trip)

	return trip, nil
}

// DeleteTrip elimina un viaje (solo el dueño)
//
// Validaciones:
// - Solo el dueño puede eliminar (userID == driver_id)
func (s *tripService) DeleteTrip(ctx context.Context, tripID string, userID int64) error {
	// Obtener el trip para verificar ownership
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return err
	}

	// Validación: Solo el dueño puede eliminar
	if trip.DriverID != userID {
		return domain.ErrUnauthorized
	}

	// Eliminar del repositorio
	if err := s.tripRepo.Delete(ctx, tripID); err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}

// CancelTrip cancela un viaje (solo el dueño)
//
// Validaciones:
// - Solo el dueño puede cancelar (userID == driver_id)
//
// Acciones:
// - Establece status = 'cancelled'
// - Registra cancelled_at, cancelled_by, cancellation_reason
// - TODO: Publicar evento trip.cancelled a RabbitMQ (Fase 5)
func (s *tripService) CancelTrip(ctx context.Context, tripID string, userID int64, request domain.CancelTripRequest) error {
	// Obtener el trip para verificar ownership
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return err
	}

	// Validación: Solo el dueño puede cancelar
	if trip.DriverID != userID {
		return domain.ErrUnauthorized
	}

	// Cancelar usando el método del repositorio
	if err := s.tripRepo.Cancel(ctx, tripID, userID, request.Reason); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to cancel trip")
		return fmt.Errorf("failed to cancel trip: %w", err)
	}

	log.Info().Str("trip_id", tripID).Int64("user_id", userID).Str("reason", request.Reason).Msg("Trip cancelled")

	// Obtener el trip actualizado para publicar el evento con el estado correcto
	cancelledTrip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to fetch cancelled trip for event publishing")
		// No retornamos el error porque el trip ya fue cancelado exitosamente
	} else {
		// Publicar evento trip.cancelled (fire-and-forget)
		s.publisher.PublishTripCancelled(ctx, cancelledTrip, userID, request.Reason)
	}

	return nil
}
