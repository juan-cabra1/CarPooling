package service

import (
	"context"
	"testing"
	"time"
	"trips-api/internal/clients"
	"trips-api/internal/domain"
	"trips-api/internal/messaging"
	"trips-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTripRepository is a mock implementation of TripRepository
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepository) FindByID(ctx context.Context, id string) (*domain.Trip, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trip), args.Error(1)
}

func (m *MockTripRepository) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error) {
	args := m.Called(ctx, filters, page, limit)
	return args.Get(0).([]domain.Trip), args.Get(1).(int64), args.Error(2)
}

func (m *MockTripRepository) Update(ctx context.Context, id string, trip *domain.Trip) error {
	args := m.Called(ctx, id, trip)
	return args.Error(0)
}

func (m *MockTripRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTripRepository) Cancel(ctx context.Context, id string, userID int64, reason string) error {
	args := m.Called(ctx, id, userID, reason)
	return args.Error(0)
}

func (m *MockTripRepository) UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error {
	args := m.Called(ctx, tripID, seatsDelta, expectedVersion)
	return args.Error(0)
}

// MockUsersClient is a mock implementation of UsersClient
type MockUsersClient struct {
	mock.Mock
}

func (m *MockUsersClient) GetUser(ctx context.Context, userID int64) (*clients.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.User), args.Error(1)
}

// MockPublisher is a mock implementation of Publisher
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishTripCreated(ctx context.Context, trip *domain.Trip) {
	m.Called(ctx, trip)
}

func (m *MockPublisher) PublishTripUpdated(ctx context.Context, trip *domain.Trip) {
	m.Called(ctx, trip)
}

func (m *MockPublisher) PublishTripCancelled(ctx context.Context, trip *domain.Trip, cancelledBy int64, reason string) {
	m.Called(ctx, trip, cancelledBy, reason)
}

func (m *MockPublisher) PublishReservationFailure(ctx context.Context, reservationID, tripID, reason string, availableSeats int) {
	m.Called(ctx, reservationID, tripID, reason, availableSeats)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ============================================================================
// CREATE TRIP TESTS
// ============================================================================

func TestCreateTrip_Success(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	driverID := int64(123)
	request := testutil.NewTestCreateTripRequest()

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository) // Not used in CreateTrip
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Driver exists
	mockUsersClient.On("GetUser", ctx, driverID).Return(&clients.User{ID: driverID}, nil)

	// Mock: Create trip succeeds
	mockRepo.On("Create", ctx, mock.MatchedBy(func(trip *domain.Trip) bool {
		return trip.DriverID == driverID &&
			trip.AvailableSeats == request.TotalSeats &&
			trip.ReservedSeats == 0 &&
			trip.Status == "published" &&
			trip.AvailabilityVersion == 1
	})).Return(nil)

	// Mock: Publish event succeeds (void method)
	mockPublisher.On("PublishTripCreated", ctx, mock.AnythingOfType("*domain.Trip"))

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	trip, err := service.CreateTrip(ctx, driverID, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, trip)
	assert.Equal(t, driverID, trip.DriverID)
	assert.Equal(t, request.TotalSeats, trip.AvailableSeats)
	assert.Equal(t, 0, trip.ReservedSeats)
	assert.Equal(t, "published", trip.Status)
	assert.Equal(t, 1, trip.AvailabilityVersion)

	mockRepo.AssertExpectations(t)
	mockUsersClient.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCreateTrip_PastDepartureDate(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	driverID := int64(123)
	request := testutil.NewTestCreateTripRequest()

	// Set departure date in the past
	pastDate := time.Now().Add(-24 * time.Hour)
	request.DepartureDatetime = pastDate.Format(time.RFC3339)

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	trip, err := service.CreateTrip(ctx, driverID, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, trip)
	assert.Equal(t, domain.ErrPastDeparture, err)

	// Verify no calls to repo or clients
	mockRepo.AssertNotCalled(t, "Create")
	mockUsersClient.AssertNotCalled(t, "GetUser")
}

func TestCreateTrip_InvalidSeats(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	driverID := int64(123)
	request := testutil.NewTestCreateTripRequest()

	// Set invalid seats (> 8)
	request.TotalSeats = 10

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	trip, err := service.CreateTrip(ctx, driverID, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, trip)
	assert.Contains(t, err.Error(), "total_seats must be between 1 and 8")

	mockRepo.AssertNotCalled(t, "Create")
	mockUsersClient.AssertNotCalled(t, "GetUser")
}

func TestCreateTrip_DriverNotFound(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	driverID := int64(999)
	request := testutil.NewTestCreateTripRequest()

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Driver not found
	mockUsersClient.On("GetUser", ctx, driverID).Return(nil, domain.ErrDriverNotFound)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	trip, err := service.CreateTrip(ctx, driverID, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, trip)
	assert.Contains(t, err.Error(), "failed to validate driver")

	mockRepo.AssertNotCalled(t, "Create")
	mockUsersClient.AssertExpectations(t)
}

// ============================================================================
// UPDATE TRIP TESTS
// ============================================================================

func TestUpdateTrip_Unauthorized(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	tripID := "trip-123"
	userID := int64(999) // Not the driver
	trip := testutil.NewTestTrip(123) // Driver is 123

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip exists
	mockRepo.On("FindByID", ctx, tripID).Return(trip, nil)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	description := "New description"
	request := domain.UpdateTripRequest{Description: &description}
	updatedTrip, err := service.UpdateTrip(ctx, tripID, userID, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedTrip)
	assert.Equal(t, domain.ErrUnauthorized, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateTrip_HasReservations(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	tripID := "trip-123"
	userID := int64(123) // Is the driver
	trip := testutil.WithAvailableSeats(testutil.NewTestTrip(123), 2) // Has reservations (2 seats reserved)

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip exists with reservations
	mockRepo.On("FindByID", ctx, tripID).Return(trip, nil)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	description := "New description"
	request := domain.UpdateTripRequest{Description: &description}
	updatedTrip, err := service.UpdateTrip(ctx, tripID, userID, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedTrip)
	assert.Equal(t, domain.ErrHasReservations, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

// ============================================================================
// EVENT PROCESSING TESTS
// ============================================================================

func TestProcessReservationCreated_Success(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	trip := testutil.WithAvailableSeats(testutil.NewTestTrip(123), 4)
	tripID := trip.ID.Hex()

	event := messaging.ReservationCreatedEvent{
		EventID:       "event-001",
		EventType:     "reservation.created",
		TripID:        tripID,
		SeatsReserved: 2,
		ReservationID: "reservation-001",
		Timestamp:     time.Now(),
	}

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip exists
	mockRepo.On("FindByID", ctx, tripID).Return(trip, nil).Times(2) // Called twice: before update and after

	// Mock: UpdateAvailability succeeds (decrease by 2 seats)
	mockRepo.On("UpdateAvailability", ctx, tripID, -2, trip.AvailabilityVersion).Return(nil)

	// Mock: Publish trip updated
	mockPublisher.On("PublishTripUpdated", ctx, trip)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	err := service.ProcessReservationCreated(ctx, event)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)

	// Verify compensation event was NOT published
	mockPublisher.AssertNotCalled(t, "PublishReservationFailure")
}

func TestProcessReservationCreated_OptimisticLockFailed(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	trip := testutil.WithAvailableSeats(testutil.NewTestTrip(123), 1) // Only 1 seat available
	tripID := trip.ID.Hex()

	event := messaging.ReservationCreatedEvent{
		EventID:       "event-002",
		EventType:     "reservation.created",
		TripID:        tripID,
		SeatsReserved: 2, // Requesting 2 seats, but only 1 available
		ReservationID: "reservation-002",
		Timestamp:     time.Now(),
	}

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip exists
	mockRepo.On("FindByID", ctx, tripID).Return(trip, nil)

	// Mock: UpdateAvailability fails (optimistic lock)
	mockRepo.On("UpdateAvailability", ctx, tripID, -2, trip.AvailabilityVersion).Return(domain.ErrOptimisticLockFailed)

	// Mock: Publish compensation event
	mockPublisher.On("PublishReservationFailure", ctx, "reservation-002", tripID, mock.Anything, 1).Return(nil)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	err := service.ProcessReservationCreated(ctx, event)

	// Assert
	assert.NoError(t, err) // Should ACK (no system error)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)

	// Verify success event was NOT published
	mockPublisher.AssertNotCalled(t, "PublishTripUpdated")
}

func TestProcessReservationCreated_TripNotFound(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	event := messaging.ReservationCreatedEvent{
		EventID:       "event-003",
		EventType:     "reservation.created",
		TripID:        "nonexistent-trip",
		SeatsReserved: 2,
		ReservationID: "reservation-003",
		Timestamp:     time.Now(),
	}

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip not found
	mockRepo.On("FindByID", ctx, "nonexistent-trip").Return(nil, domain.ErrTripNotFound)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	err := service.ProcessReservationCreated(ctx, event)

	// Assert
	assert.NoError(t, err) // Should ACK (not a system error)
	mockRepo.AssertExpectations(t)

	// No updates or events should be published
	mockRepo.AssertNotCalled(t, "UpdateAvailability")
	mockPublisher.AssertNotCalled(t, "PublishReservationFailure")
	mockPublisher.AssertNotCalled(t, "PublishTripUpdated")
}

func TestProcessReservationCancelled_Success(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	trip := testutil.WithAvailableSeats(testutil.NewTestTrip(123), 2) // 2 available, 2 reserved
	tripID := trip.ID.Hex()

	event := messaging.ReservationCancelledEvent{
		EventID:       "event-004",
		EventType:     "reservation.cancelled",
		TripID:        tripID,
		SeatsReleased: 2,
		ReservationID: "reservation-004",
		Timestamp:     time.Now(),
	}

	mockRepo := new(MockTripRepository)
	mockIdempotency := new(MockEventRepository)
	mockUsersClient := new(MockUsersClient)
	mockPublisher := new(MockPublisher)

	// Mock: Trip exists
	mockRepo.On("FindByID", ctx, tripID).Return(trip, nil).Times(2) // Called twice

	// Mock: UpdateAvailability succeeds (increase by 2 seats)
	mockRepo.On("UpdateAvailability", ctx, tripID, 2, trip.AvailabilityVersion).Return(nil)

	// Mock: Publish trip updated
	mockPublisher.On("PublishTripUpdated", ctx, trip)

	idempotencyService := NewIdempotencyService(mockIdempotency)
	service := NewTripService(mockRepo, idempotencyService, mockUsersClient, mockPublisher)

	// Act
	err := service.ProcessReservationCancelled(ctx, event)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
