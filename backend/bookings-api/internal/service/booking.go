package service

import (
	"fmt"
	"time"

	"github.com/carpooling-ucc/bookings-api/internal/dao"
	"github.com/carpooling-ucc/bookings-api/internal/domain"
	"github.com/carpooling-ucc/bookings-api/internal/repository"
	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(req domain.CreateBookingRequest, passengerID int64) (*domain.BookingDTO, error)
	GetBookingByID(id string, userID int64) (*domain.BookingDTO, error)
	CancelBooking(id string, userID int64) (*domain.BookingDTO, error)
	ConfirmArrival(id string, userID int64) (*domain.BookingDTO, error)
	GetUserBookings(userID int64, page, limit int, status string) (*domain.BookingListResponse, error)
	GetTripBookings(tripID string, userID int64) ([]*domain.BookingDTO, error)
	// Internal methods (no auth required)
	GetBookingByIDInternal(id string) (*domain.BookingDTO, error)
	CompleteBookingInternal(id string) (*domain.BookingDTO, error)
}

type bookingServiceImpl struct {
	bookingRepo       repository.BookingRepository
	validationService ValidationService
	tripsClient       TripsClient
	rabbitClient      RabbitMQClient
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	validationService ValidationService,
	tripsClient TripsClient,
	rabbitClient RabbitMQClient,
) BookingService {
	return &bookingServiceImpl{
		bookingRepo:       bookingRepo,
		validationService: validationService,
		tripsClient:       tripsClient,
		rabbitClient:      rabbitClient,
	}
}

// CreateBooking creates a new booking with concurrent validations
func (s *bookingServiceImpl) CreateBooking(req domain.CreateBookingRequest, passengerID int64) (*domain.BookingDTO, error) {
	// Run concurrent validations
	valCtx := ValidationContext{
		TripID:        req.TripID,
		PassengerID:   passengerID,
		SeatsReserved: req.SeatsReserved,
	}

	trip, err := s.validationService.RunConcurrentValidations(valCtx)
	if err != nil {
		return nil, err
	}

	// Create booking entity
	booking := &dao.BookingDAO{
		ID:            uuid.New().String(),
		TripID:        req.TripID,
		PassengerID:   passengerID,
		DriverID:      trip.DriverID,
		SeatsReserved: req.SeatsReserved,
		PricePerSeat:  trip.PricePerSeat,
		TotalAmount:   float64(req.SeatsReserved) * trip.PricePerSeat,
		Status:        "pending",
		PaymentStatus: "pending",
		ArrivedSafely: false,
	}

	// Save to database
	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Update trip seats
	if err := s.tripsClient.UpdateTripSeats(req.TripID, -req.SeatsReserved); err != nil {
		// Rollback - delete the booking
		s.bookingRepo.Delete(booking.ID)
		return nil, fmt.Errorf("failed to update trip seats: %w", err)
	}

	// Publish event to RabbitMQ
	eventPayload := map[string]interface{}{
		"reservation_id": booking.ID,
		"trip_id":        booking.TripID,
		"passenger_id":   booking.PassengerID,
		"seats_reserved": booking.SeatsReserved,
	}
	s.rabbitClient.PublishEvent("reservation.created", eventPayload)

	return toBookingDTO(booking), nil
}

func (s *bookingServiceImpl) GetBookingByID(id string, userID int64) (*domain.BookingDTO, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Check if user has permission to view this booking
	if booking.PassengerID != userID && booking.DriverID != userID {
		return nil, fmt.Errorf("forbidden: you don't have permission to view this booking")
	}

	return toBookingDTO(booking), nil
}

func (s *bookingServiceImpl) CancelBooking(id string, userID int64) (*domain.BookingDTO, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Only passenger can cancel
	if booking.PassengerID != userID {
		return nil, fmt.Errorf("forbidden: only the passenger can cancel the booking")
	}

	// Check if already cancelled or completed
	if booking.Status == "cancelled" {
		return nil, fmt.Errorf("booking is already cancelled")
	}
	if booking.Status == "completed" {
		return nil, fmt.Errorf("cannot cancel a completed booking")
	}

	// Get trip to check departure date
	trip, err := s.tripsClient.GetTrip(booking.TripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip details: %w", err)
	}

	// Check if cancellation is within 24 hours
	timeUntilDeparture := time.Until(trip.DepartureDate)
	if timeUntilDeparture < 24*time.Hour {
		return nil, fmt.Errorf("cannot cancel booking less than 24 hours before departure")
	}

	// Update booking status
	booking.Status = "cancelled"
	booking.PaymentStatus = "refunded"

	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	// Restore seats in trip
	if err := s.tripsClient.UpdateTripSeats(booking.TripID, booking.SeatsReserved); err != nil {
		// Log error but don't fail the cancellation
		fmt.Printf("Warning: failed to restore trip seats: %v\n", err)
	}

	// Publish event
	eventPayload := map[string]interface{}{
		"reservation_id": booking.ID,
		"trip_id":        booking.TripID,
		"seats_restored": booking.SeatsReserved,
	}
	s.rabbitClient.PublishEvent("reservation.cancelled", eventPayload)

	return toBookingDTO(booking), nil
}

func (s *bookingServiceImpl) ConfirmArrival(id string, userID int64) (*domain.BookingDTO, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Only passenger can confirm arrival
	if booking.PassengerID != userID {
		return nil, fmt.Errorf("forbidden: only the passenger can confirm arrival")
	}

	// Check if booking is confirmed
	if booking.Status != "confirmed" {
		return nil, fmt.Errorf("booking must be in confirmed status")
	}

	// Update arrival confirmation
	now := time.Now()
	booking.ArrivedSafely = true
	booking.ArrivalConfirmedAt = &now

	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	// Check if all passengers confirmed - if so, mark as completed
	// For now, we'll mark this specific booking as completed
	booking.Status = "completed"
	s.bookingRepo.Update(booking)

	// Publish event
	eventPayload := map[string]interface{}{
		"reservation_id": booking.ID,
		"trip_id":        booking.TripID,
	}
	s.rabbitClient.PublishEvent("reservation.completed", eventPayload)

	return toBookingDTO(booking), nil
}

func (s *bookingServiceImpl) GetUserBookings(userID int64, page, limit int, status string) (*domain.BookingListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	bookings, total, err := s.bookingRepo.FindByPassengerID(userID, page, limit, status)
	if err != nil {
		return nil, err
	}

	bookingDTOs := make([]*domain.BookingDTO, len(bookings))
	for i, booking := range bookings {
		bookingDTOs[i] = toBookingDTO(booking)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &domain.BookingListResponse{
		Bookings: bookingDTOs,
		Pagination: domain.PaginationInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *bookingServiceImpl) GetTripBookings(tripID string, userID int64) ([]*domain.BookingDTO, error) {
	// Validate that user is the driver
	driverID, err := s.tripsClient.GetTripDriver(tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip driver: %w", err)
	}

	if driverID != userID {
		return nil, fmt.Errorf("forbidden: only the trip driver can view all bookings")
	}

	bookings, err := s.bookingRepo.FindByTripID(tripID)
	if err != nil {
		return nil, err
	}

	bookingDTOs := make([]*domain.BookingDTO, len(bookings))
	for i, booking := range bookings {
		bookingDTOs[i] = toBookingDTO(booking)
	}

	return bookingDTOs, nil
}

// Internal methods

func (s *bookingServiceImpl) GetBookingByIDInternal(id string) (*domain.BookingDTO, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}
	return toBookingDTO(booking), nil
}

func (s *bookingServiceImpl) CompleteBookingInternal(id string) (*domain.BookingDTO, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	booking.Status = "completed"
	if err := s.bookingRepo.Update(booking); err != nil {
		return nil, err
	}

	return toBookingDTO(booking), nil
}

// Helper function to convert DAO to DTO
func toBookingDTO(dao *dao.BookingDAO) *domain.BookingDTO {
	return &domain.BookingDTO{
		ID:                 dao.ID,
		TripID:             dao.TripID,
		PassengerID:        dao.PassengerID,
		DriverID:           dao.DriverID,
		SeatsReserved:      dao.SeatsReserved,
		PricePerSeat:       dao.PricePerSeat,
		TotalAmount:        dao.TotalAmount,
		Status:             dao.Status,
		PaymentStatus:      dao.PaymentStatus,
		ArrivedSafely:      dao.ArrivedSafely,
		ArrivalConfirmedAt: dao.ArrivalConfirmedAt,
		CreatedAt:          dao.CreatedAt,
		UpdatedAt:          dao.UpdatedAt,
	}
}
