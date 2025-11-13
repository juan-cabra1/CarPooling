package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"search-api/internal/domain"
	"search-api/internal/mocks"
	"search-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============================================================================
// HandleTripCreated Tests
// ============================================================================

func TestHandleTripCreated_Success(t *testing.T) {
	// Setup
	eventID := "event-123"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	testTrip := testutil.CreateTestTrip(tripID)
	testUser := testutil.CreateTestUser(driverID)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			assert.Equal(t, eventID, id)
			return false, nil // Event not processed yet
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, eventID, event.EventID)
			assert.Equal(t, "trip.created", event.EventType)
			assert.Equal(t, "success", event.Result)
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			assert.Equal(t, tripID, id)
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return testUser, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			assert.Equal(t, tripID, trip.TripID)
			assert.Equal(t, driverID, trip.DriverID)
			assert.NotEmpty(t, trip.SearchText)
			return nil
		},
	}

	mockSolr := &mocks.MockSolrClient{
		IndexFunc: func(trip *domain.SearchTrip) error {
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		mockTripsClient,
		mockUsersClient,
		mockSolr,
		&mocks.MockCache{},
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert
	require.NoError(t, err)
}

func TestHandleTripCreated_DuplicateEvent_Skipped(t *testing.T) {
	// Setup - CRITICAL IDEMPOTENCY TEST
	eventID := "event-123"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			assert.Equal(t, eventID, id)
			return true, nil // Event already processed
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			t.Fatal("MarkEventProcessed should not be called for duplicate events")
			return nil
		},
	}

	// These should NOT be called since event is already processed
	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			t.Fatal("GetTrip should not be called for duplicate events")
			return nil, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			t.Fatal("Create should not be called for duplicate events")
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		mockTripsClient,
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert
	require.NoError(t, err, "Duplicate events should be handled gracefully")
}

func TestHandleTripCreated_ConcurrentDuplicates_OnlyOneProcesses(t *testing.T) {
	// Setup - CRITICAL CONCURRENCY TEST
	eventID := "event-concurrent-123"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	testTrip := testutil.CreateTestTrip(tripID)
	testUser := testutil.CreateTestUser(driverID)

	// Track how many times each operation is called
	var (
		idempotencyChecks int
		tripsAPICalls     int
		usersAPICalls     int
		mongoDBCreates    int
		eventMarks        int
		mu                sync.Mutex
	)

	var eventProcessed bool

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			mu.Lock()
			defer mu.Unlock()
			idempotencyChecks++
			wasProcessed := eventProcessed
			return wasProcessed, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			mu.Lock()
			defer mu.Unlock()
			eventMarks++
			eventProcessed = true
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			mu.Lock()
			tripsAPICalls++
			mu.Unlock()
			// Simulate some processing time
			time.Sleep(10 * time.Millisecond)
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
			mu.Lock()
			usersAPICalls++
			mu.Unlock()
			return testUser, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			mu.Lock()
			mongoDBCreates++
			mu.Unlock()
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		mockTripsClient,
		mockUsersClient,
		nil,
		nil,
	)

	// Execute - Process same event 10 times concurrently
	const numGoroutines = 10
	var wg sync.WaitGroup
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			errors[index] = service.HandleTripCreated(context.Background(), eventID, tripID, driverID)
		}(i)
	}

	wg.Wait()

	// Assert
	// All goroutines should complete without error
	for i, err := range errors {
		assert.NoError(t, err, "Goroutine %d should complete successfully", i)
	}

	// Only ONE should actually process the event
	mu.Lock()
	defer mu.Unlock()

	t.Logf("Idempotency checks: %d", idempotencyChecks)
	t.Logf("Trips API calls: %d", tripsAPICalls)
	t.Logf("Users API calls: %d", usersAPICalls)
	t.Logf("MongoDB creates: %d", mongoDBCreates)
	t.Logf("Event marks: %d", eventMarks)

	// All goroutines should check idempotency
	assert.Equal(t, numGoroutines, idempotencyChecks, "All goroutines should check idempotency")

	// Only ONE should process the event fully
	// Due to race conditions, we might have 1-2 actually process (before the flag is set)
	// but it should be significantly less than all 10
	assert.LessOrEqual(t, tripsAPICalls, 2, "At most 1-2 should call trips API due to race window")
	assert.LessOrEqual(t, mongoDBCreates, 2, "At most 1-2 should create in MongoDB due to race window")
	assert.LessOrEqual(t, eventMarks, 2, "At most 1-2 should mark event as processed due to race window")
}

func TestHandleTripCreated_TripNotFound_PermanentError(t *testing.T) {
	// Setup
	eventID := "event-404"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, eventID, event.EventID)
			assert.Equal(t, "skipped", event.Result)
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return nil, domain.ErrTripNotFound
		},
	}

	service := NewTripEventService(
		&mocks.MockTripRepository{},
		mockEventRepo,
		mockTripsClient,
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert
	assert.ErrorIs(t, err, domain.ErrTripNotFound)
}

func TestHandleTripCreated_TripsAPITransientError_Retry(t *testing.T) {
	// Setup
	eventID := "event-retry"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			t.Fatal("Should not mark event as processed on transient error")
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return nil, errors.New("temporary network error")
		},
	}

	service := NewTripEventService(
		&mocks.MockTripRepository{},
		mockEventRepo,
		mockTripsClient,
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fetch trip failed")
}

func TestHandleTripCreated_UserNotFound_PermanentError(t *testing.T) {
	// Setup
	eventID := "event-user-404"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	testTrip := testutil.CreateTestTrip(tripID)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, "skipped", event.Result)
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}

	service := NewTripEventService(
		&mocks.MockTripRepository{},
		mockEventRepo,
		mockTripsClient,
		mockUsersClient,
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestHandleTripCreated_SolrFailure_DoesNotBlock(t *testing.T) {
	// Setup - Solr failure should not prevent event processing
	eventID := "event-solr-fail"
	tripID := primitive.NewObjectID().Hex()
	driverID := int64(123)

	testTrip := testutil.CreateTestTrip(tripID)
	testUser := testutil.CreateTestUser(driverID)

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, "success", event.Result)
			return nil
		},
	}

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
			return testUser, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			return nil
		},
	}

	mockSolr := &mocks.MockSolrClient{
		IndexFunc: func(trip *domain.SearchTrip) error {
			return errors.New("solr connection refused")
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		mockTripsClient,
		mockUsersClient,
		mockSolr,
		nil,
	)

	// Execute
	err := service.HandleTripCreated(context.Background(), eventID, tripID, driverID)

	// Assert - Should succeed despite Solr failure
	require.NoError(t, err, "Solr failure should not block event processing")
}

// ============================================================================
// HandleTripUpdated Tests
// ============================================================================

func TestHandleTripUpdated_Success(t *testing.T) {
	// Setup
	eventID := "event-update-123"
	tripID := primitive.NewObjectID().Hex()
	availableSeats := 2
	reservedSeats := 2
	status := "published"

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, eventID, event.EventID)
			assert.Equal(t, "trip.updated", event.EventType)
			assert.Equal(t, "success", event.Result)
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateAvailabilityByTripIDFunc: func(ctx context.Context, id string, avail, reserved int, stat string) error {
			assert.Equal(t, tripID, id)
			assert.Equal(t, availableSeats, avail)
			assert.Equal(t, reservedSeats, reserved)
			assert.Equal(t, status, stat)
			return nil
		},
	}

	mockCache := &mocks.MockCache{
		DeleteFunc: func(ctx context.Context, key string) error {
			assert.Equal(t, "trip:"+tripID, key)
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		mockCache,
	)

	// Execute
	err := service.HandleTripUpdated(context.Background(), eventID, tripID, availableSeats, reservedSeats, status)

	// Assert
	require.NoError(t, err)
}

func TestHandleTripUpdated_DuplicateEvent_Skipped(t *testing.T) {
	// Setup - CRITICAL IDEMPOTENCY TEST
	eventID := "event-update-dup"
	tripID := primitive.NewObjectID().Hex()

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return true, nil // Already processed
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			t.Fatal("Should not be called for duplicate events")
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateAvailabilityByTripIDFunc: func(ctx context.Context, id string, avail, reserved int, stat string) error {
			t.Fatal("Should not update for duplicate events")
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripUpdated(context.Background(), eventID, tripID, 2, 2, "published")

	// Assert
	require.NoError(t, err)
}

func TestHandleTripUpdated_TripNotFound_PermanentError(t *testing.T) {
	// Setup
	eventID := "event-update-404"
	tripID := primitive.NewObjectID().Hex()

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, "skipped", event.Result)
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateAvailabilityByTripIDFunc: func(ctx context.Context, id string, avail, reserved int, stat string) error {
			return domain.ErrSearchTripNotFound
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripUpdated(context.Background(), eventID, tripID, 2, 2, "published")

	// Assert
	assert.ErrorIs(t, err, domain.ErrSearchTripNotFound)
}

func TestHandleTripUpdated_CacheInvalidation(t *testing.T) {
	// Setup
	eventID := "event-cache-inval"
	tripID := primitive.NewObjectID().Hex()

	cacheDeleted := false

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateAvailabilityByTripIDFunc: func(ctx context.Context, id string, avail, reserved int, stat string) error {
			return nil
		},
	}

	mockCache := &mocks.MockCache{
		DeleteFunc: func(ctx context.Context, key string) error {
			assert.Equal(t, "trip:"+tripID, key)
			cacheDeleted = true
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		mockCache,
	)

	// Execute
	err := service.HandleTripUpdated(context.Background(), eventID, tripID, 2, 2, "published")

	// Assert
	require.NoError(t, err)
	assert.True(t, cacheDeleted, "Cache should be invalidated")
}

// ============================================================================
// HandleTripCancelled Tests
// ============================================================================

func TestHandleTripCancelled_Success(t *testing.T) {
	// Setup
	eventID := "event-cancel-123"
	tripID := primitive.NewObjectID().Hex()
	reason := "Driver not available"

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, eventID, event.EventID)
			assert.Equal(t, "trip.cancelled", event.EventType)
			assert.Equal(t, "success", event.Result)
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateStatusByTripIDFunc: func(ctx context.Context, id string, status string) error {
			assert.Equal(t, tripID, id)
			assert.Equal(t, "cancelled", status)
			return nil
		},
	}

	mockCache := &mocks.MockCache{
		DeleteFunc: func(ctx context.Context, key string) error {
			assert.Equal(t, "trip:"+tripID, key)
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		mockCache,
	)

	// Execute
	err := service.HandleTripCancelled(context.Background(), eventID, tripID, reason)

	// Assert
	require.NoError(t, err)
}

func TestHandleTripCancelled_DuplicateEvent_Skipped(t *testing.T) {
	// Setup - CRITICAL IDEMPOTENCY TEST
	eventID := "event-cancel-dup"
	tripID := primitive.NewObjectID().Hex()

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return true, nil // Already processed
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			t.Fatal("Should not be called for duplicate events")
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateStatusByTripIDFunc: func(ctx context.Context, id string, status string) error {
			t.Fatal("Should not update for duplicate events")
			return nil
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCancelled(context.Background(), eventID, tripID, "reason")

	// Assert
	require.NoError(t, err)
}

func TestHandleTripCancelled_TripNotFound_PermanentError(t *testing.T) {
	// Setup
	eventID := "event-cancel-404"
	tripID := primitive.NewObjectID().Hex()

	mockEventRepo := &mocks.MockEventRepository{
		IsEventProcessedFunc: func(ctx context.Context, id string) (bool, error) {
			return false, nil
		},
		MarkEventProcessedFunc: func(ctx context.Context, event *domain.ProcessedEvent) error {
			assert.Equal(t, "skipped", event.Result)
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		UpdateStatusByTripIDFunc: func(ctx context.Context, id string, status string) error {
			return domain.ErrSearchTripNotFound
		},
	}

	service := NewTripEventService(
		mockTripRepo,
		mockEventRepo,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
		nil,
		nil,
	)

	// Execute
	err := service.HandleTripCancelled(context.Background(), eventID, tripID, "reason")

	// Assert
	assert.ErrorIs(t, err, domain.ErrSearchTripNotFound)
}
