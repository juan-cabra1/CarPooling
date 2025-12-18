package service

import (
	"context"
	"fmt"
	"math"
	"time"
	"trips-api/internal/domain"
	"trips-api/internal/repository"

	"github.com/rs/zerolog/log"
)

// TrackingService define las operaciones de tracking de viajes
type TrackingService interface {
	// StartTrip inicia un viaje y marca el timestamp de inicio
	StartTrip(ctx context.Context, tripID string, driverID int64) error

	// UpdateLocation actualiza la ubicación actual del conductor
	UpdateLocation(ctx context.Context, tripID string, driverID int64, location domain.LocationPoint) error

	// CompleteTrip marca un viaje como completado
	CompleteTrip(ctx context.Context, tripID string, driverID int64) error

	// GetTripTracking obtiene el estado actual de tracking de un viaje
	GetTripTracking(ctx context.Context, tripID string) (*domain.Trip, error)
}

type trackingService struct {
	tripRepo repository.TripRepository
}

// NewTrackingService crea una nueva instancia del servicio de tracking
func NewTrackingService(tripRepo repository.TripRepository) TrackingService {
	return &trackingService{
		tripRepo: tripRepo,
	}
}

// StartTrip inicia un viaje y marca el timestamp de inicio
func (s *trackingService) StartTrip(ctx context.Context, tripID string, driverID int64) error {
	// Obtener el viaje
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}

	// Verificar que el usuario sea el conductor
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}

	// Verificar que el viaje esté en estado published
	if trip.Status != "published" {
		return fmt.Errorf("trip is not in published state (current: %s)", trip.Status)
	}

	// Actualizar el viaje
	now := time.Now()
	trip.Status = "in_progress"
	trip.StartedAt = &now

	// Guardar cambios
	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	log.Info().
		Str("trip_id", tripID).
		Int64("driver_id", driverID).
		Msg("Trip started")

	return nil
}

// UpdateLocation actualiza la ubicación actual del conductor
func (s *trackingService) UpdateLocation(ctx context.Context, tripID string, driverID int64, location domain.LocationPoint) error {
	// Obtener el viaje
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}

	// Verificar que el usuario sea el conductor
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}

	// Verificar que el viaje esté en progreso
	if trip.Status != "in_progress" {
		return fmt.Errorf("trip is not in progress (current: %s)", trip.Status)
	}

	// Actualizar ubicación actual
	coords := &domain.Coordinates{
		Lat: location.Lat,
		Lng: location.Lng,
	}
	trip.CurrentLocation = coords

	// Agregar a historial de ubicaciones
	if trip.LocationHistory == nil {
		trip.LocationHistory = []domain.LocationPoint{}
	}
	trip.LocationHistory = append(trip.LocationHistory, location)

	// Calcular progreso del viaje
	if trip.TripProgress == nil {
		trip.TripProgress = &domain.TripProgress{}
	}

	// Calcular distancia recorrida desde el origen
	distanceTraveled := calculateDistance(
		trip.Origin.Coordinates.Lat,
		trip.Origin.Coordinates.Lng,
		location.Lat,
		location.Lng,
	)
	trip.TripProgress.DistanceTraveled = distanceTraveled

	// Calcular distancia restante al destino
	distanceRemaining := calculateDistance(
		location.Lat,
		location.Lng,
		trip.Destination.Coordinates.Lat,
		trip.Destination.Coordinates.Lng,
	)

	// Estimar tiempo restante (basado en velocidad o velocidad promedio)
	speed := location.Speed
	if speed == 0 {
		speed = 50.0 // Velocidad promedio de 50 km/h si no hay datos
	}
	etaMinutes := int((distanceRemaining / speed) * 60)
	trip.TripProgress.EstimatedTimeRemaining = etaMinutes
	trip.TripProgress.LastUpdated = time.Now()

	// Guardar cambios
	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		return fmt.Errorf("failed to update trip location: %w", err)
	}

	log.Debug().
		Str("trip_id", tripID).
		Float64("lat", location.Lat).
		Float64("lng", location.Lng).
		Float64("distance_traveled", distanceTraveled).
		Float64("distance_remaining", distanceRemaining).
		Int("eta_minutes", etaMinutes).
		Msg("Trip location updated")

	return nil
}

// CompleteTrip marca un viaje como completado
func (s *trackingService) CompleteTrip(ctx context.Context, tripID string, driverID int64) error {
	// Obtener el viaje
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}

	// Verificar que el usuario sea el conductor
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}

	// Verificar que el viaje esté en progreso
	if trip.Status != "in_progress" {
		return fmt.Errorf("trip is not in progress (current: %s)", trip.Status)
	}

	// Actualizar el viaje
	now := time.Now()
	trip.Status = "completed"
	trip.CompletedAt = &now

	// Guardar cambios
	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	log.Info().
		Str("trip_id", tripID).
		Int64("driver_id", driverID).
		Msg("Trip completed")

	return nil
}

// GetTripTracking obtiene el estado actual de tracking de un viaje
func (s *trackingService) GetTripTracking(ctx context.Context, tripID string) (*domain.Trip, error) {
	// Obtener el viaje
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to find trip: %w", err)
	}

	return trip, nil
}

// calculateDistance calcula la distancia entre dos puntos usando la fórmula de Haversine
// Retorna la distancia en kilómetros
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0 // Radio de la Tierra en km

	// Convertir a radianes
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// Diferencias
	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	// Fórmula de Haversine
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c

	// Redondear a 2 decimales
	return math.Round(distance*100) / 100
}
