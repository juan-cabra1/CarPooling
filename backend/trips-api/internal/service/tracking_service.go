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
	StartTrip(ctx context.Context, tripID string, driverID int64) error
	UpdateLocation(ctx context.Context, tripID string, driverID int64, location domain.LocationPoint) error
	CompleteTrip(ctx context.Context, tripID string, driverID int64) error
	GetTripTracking(ctx context.Context, tripID string) (*domain.Trip, error)
}

type trackingService struct {
	tripRepo repository.TripRepository
}

func NewTrackingService(tripRepo repository.TripRepository) TrackingService {
	return &trackingService{tripRepo: tripRepo}
}

func (s *trackingService) StartTrip(ctx context.Context, tripID string, driverID int64) error {
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}
	if trip.Status != "published" {
		return fmt.Errorf("trip is not in published state (current: %s)", trip.Status)
	}

	now := time.Now()
	trip.Status = "in_progress"
	trip.StartedAt = &now

	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	log.Info().Str("trip_id", tripID).Int64("driver_id", driverID).Msg("Trip started")
	return nil
}

func (s *trackingService) UpdateLocation(ctx context.Context, tripID string, driverID int64, location domain.LocationPoint) error {
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}
	if trip.Status != "in_progress" {
		return fmt.Errorf("trip is not in progress (current: %s)", trip.Status)
	}

	trip.CurrentLocation = &domain.Coordinates{
		Lat: location.Lat,
		Lng: location.Lng,
	}

	if trip.LocationHistory == nil {
		trip.LocationHistory = []domain.LocationPoint{}
	}
	trip.LocationHistory = append(trip.LocationHistory, location)

	if trip.TripProgress == nil {
		trip.TripProgress = &domain.TripProgress{}
	}

	distanceTraveled := calculateDistance(
		trip.Origin.Coordinates.Lat,
		trip.Origin.Coordinates.Lng,
		location.Lat,
		location.Lng,
	)
	trip.TripProgress.DistanceTraveled = distanceTraveled

	distanceRemaining := calculateDistance(
		location.Lat,
		location.Lng,
		trip.Destination.Coordinates.Lat,
		trip.Destination.Coordinates.Lng,
	)

	speed := location.Speed
	if speed == 0 {
		speed = 50.0
	}
	etaMinutes := int((distanceRemaining / speed) * 60)
	trip.TripProgress.EstimatedTimeRemaining = etaMinutes
	trip.TripProgress.LastUpdated = time.Now()

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

func (s *trackingService) CompleteTrip(ctx context.Context, tripID string, driverID int64) error {
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to find trip: %w", err)
	}
	if trip.DriverID != driverID {
		return fmt.Errorf("user %d is not the driver of trip %s", driverID, tripID)
	}
	if trip.Status != "in_progress" {
		return fmt.Errorf("trip is not in progress (current: %s)", trip.Status)
	}

	now := time.Now()
	trip.Status = "completed"
	trip.CompletedAt = &now

	if err := s.tripRepo.Update(ctx, tripID, trip); err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	log.Info().Str("trip_id", tripID).Int64("driver_id", driverID).Msg("Trip completed")
	return nil
}

func (s *trackingService) GetTripTracking(ctx context.Context, tripID string) (*domain.Trip, error) {
	trip, err := s.tripRepo.FindByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to find trip: %w", err)
	}
	return trip, nil
}

// calculateDistance uses the Haversine formula to compute distance in km
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return math.Round(earthRadius*c*100) / 100
}
