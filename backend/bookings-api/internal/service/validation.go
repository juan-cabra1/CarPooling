package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/carpooling-ucc/bookings-api/internal/repository"
)

type ValidationResult struct {
	Success bool
	Error   error
	Data    interface{}
}

type ValidationContext struct {
	TripID        string
	PassengerID   int64
	SeatsReserved int
}

type ValidationService interface {
	RunConcurrentValidations(ctx ValidationContext) (*TripDTO, error)
}

type validationServiceImpl struct {
	tripsClient TripsClient
	usersClient UsersClient
	bookingRepo repository.BookingRepository
}

func NewValidationService(tripsClient TripsClient, usersClient UsersClient, bookingRepo repository.BookingRepository) ValidationService {
	return &validationServiceImpl{
		tripsClient: tripsClient,
		usersClient: usersClient,
		bookingRepo: bookingRepo,
	}
}

func (s *validationServiceImpl) RunConcurrentValidations(ctx ValidationContext) (*TripDTO, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, 4) // Buffer for 4 validations
	tripChan := make(chan *TripDTO, 1)

	// Goroutine 1: Validate trip availability
	wg.Add(1)
	go func() {
		defer wg.Done()
		trip, err := s.validateTripAvailability(ctx.TripID, ctx.SeatsReserved)
		if err != nil {
			errChan <- fmt.Errorf("trip validation failed: %w", err)
			return
		}
		tripChan <- trip
	}()

	// Goroutine 2: Validate user exists
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.usersClient.ValidateUser(ctx.PassengerID)
		if err != nil {
			errChan <- fmt.Errorf("user validation failed: %w", err)
		}
	}()

	// Goroutine 3: Check for duplicate booking
	wg.Add(1)
	go func() {
		defer wg.Done()
		isDuplicate, err := s.bookingRepo.CheckDuplicateBooking(ctx.TripID, ctx.PassengerID)
		if err != nil {
			errChan <- fmt.Errorf("duplicate check failed: %w", err)
			return
		}
		if isDuplicate {
			errChan <- fmt.Errorf("passenger already has a booking for this trip")
		}
	}()

	// Wait for all goroutines with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines completed
		close(errChan)
		close(tripChan)
	case <-time.After(5 * time.Second):
		close(errChan)
		close(tripChan)
		return nil, fmt.Errorf("validation timeout")
	}

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Get trip data
	trip := <-tripChan
	if trip == nil {
		return nil, fmt.Errorf("failed to retrieve trip data")
	}

	// Goroutine 4: Validate passenger is not the driver (needs trip data)
	if trip.DriverID == ctx.PassengerID {
		return nil, fmt.Errorf("passenger cannot be the driver")
	}

	return trip, nil
}

func (s *validationServiceImpl) validateTripAvailability(tripID string, seatsReserved int) (*TripDTO, error) {
	trip, err := s.tripsClient.GetTrip(tripID)
	if err != nil {
		return nil, err
	}

	// Validate trip status
	if trip.Status != "published" {
		return nil, fmt.Errorf("trip is not available for booking (status: %s)", trip.Status)
	}

	// Validate available seats
	if trip.AvailableSeats < seatsReserved {
		return nil, fmt.Errorf("not enough available seats (available: %d, requested: %d)", 
			trip.AvailableSeats, seatsReserved)
	}

	// Validate departure date is in the future
	if trip.DepartureDate.Before(time.Now()) {
		return nil, fmt.Errorf("trip has already departed")
	}

	return trip, nil
}
