package service

import (
	"context"
	"errors"
	"testing"
	"trips-api/internal/domain"
	"trips-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventRepository) MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// TestCheckAndMarkEvent_FirstEventProcessed tests that a new event is marked and should be processed
func TestCheckAndMarkEvent_FirstEventProcessed(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	eventID := "event-123"
	eventType := "reservation.created"

	mockRepo := new(MockEventRepository)

	// El evento NO ha sido procesado antes
	mockRepo.On("IsEventProcessed", ctx, eventID).Return(false, nil)

	// Mock MarkEventProcessed - acepta cualquier ProcessedEvent con el eventID correcto
	mockRepo.On("MarkEventProcessed", ctx, mock.MatchedBy(func(event *domain.ProcessedEvent) bool {
		return event.EventID == eventID && event.EventType == eventType
	})).Return(nil)

	service := NewIdempotencyService(mockRepo)

	// Act
	shouldProcess, err := service.CheckAndMarkEvent(ctx, eventID, eventType)

	// Assert
	assert.NoError(t, err)
	assert.True(t, shouldProcess, "First event should be processed")
	mockRepo.AssertExpectations(t)
}

// TestCheckAndMarkEvent_DuplicateEventSkipped tests that a duplicate event is skipped
func TestCheckAndMarkEvent_DuplicateEventSkipped(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	eventID := "event-456"
	eventType := "reservation.cancelled"

	mockRepo := new(MockEventRepository)

	// El evento YA fue procesado
	mockRepo.On("IsEventProcessed", ctx, eventID).Return(true, nil)
	// MarkEventProcessed NO debe ser llamado porque el evento ya fue procesado

	service := NewIdempotencyService(mockRepo)

	// Act
	shouldProcess, err := service.CheckAndMarkEvent(ctx, eventID, eventType)

	// Assert
	assert.NoError(t, err)
	assert.False(t, shouldProcess, "Duplicate event should NOT be processed")
	mockRepo.AssertExpectations(t)

	// Verificar que MarkEventProcessed NO fue llamado
	mockRepo.AssertNotCalled(t, "MarkEventProcessed")
}

// TestCheckAndMarkEvent_IsEventProcessedError tests error handling when checking if event is processed
func TestCheckAndMarkEvent_IsEventProcessedError(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	eventID := "event-789"
	eventType := "reservation.created"

	mockRepo := new(MockEventRepository)

	// Simular error en IsEventProcessed (e.g., MongoDB connection error)
	expectedError := errors.New("mongodb connection error")
	mockRepo.On("IsEventProcessed", ctx, eventID).Return(false, expectedError)

	service := NewIdempotencyService(mockRepo)

	// Act
	shouldProcess, err := service.CheckAndMarkEvent(ctx, eventID, eventType)

	// Assert
	assert.Error(t, err)
	assert.False(t, shouldProcess)
	assert.Contains(t, err.Error(), "failed to check if event is processed")
	mockRepo.AssertExpectations(t)
}

// TestCheckAndMarkEvent_MarkEventProcessedError tests error handling when marking event as processed
func TestCheckAndMarkEvent_MarkEventProcessedError(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()
	eventID := "event-999"
	eventType := "reservation.created"

	mockRepo := new(MockEventRepository)

	// El evento NO ha sido procesado
	mockRepo.On("IsEventProcessed", ctx, eventID).Return(false, nil)

	// Simular error al marcar evento (error real de sistema, no duplicate key)
	expectedError := errors.New("mongodb write error")
	mockRepo.On("MarkEventProcessed", ctx, mock.MatchedBy(func(event *domain.ProcessedEvent) bool {
		return event.EventID == eventID
	})).Return(expectedError)

	service := NewIdempotencyService(mockRepo)

	// Act
	shouldProcess, err := service.CheckAndMarkEvent(ctx, eventID, eventType)

	// Assert
	assert.Error(t, err)
	assert.False(t, shouldProcess)
	assert.Contains(t, err.Error(), "failed to mark event as processed")
	mockRepo.AssertExpectations(t)
}

// TestCheckAndMarkEvent_ConcurrentDuplicates tests race condition handling
// CRITICAL: Este test debe ejecutarse con -race flag para detectar race conditions
//
// go test ./internal/service -run TestCheckAndMarkEvent_ConcurrentDuplicates -race -v
//
// NOTE: This is a simplified test for the basic plan. Full concurrency testing
// requires integration tests with real MongoDB to verify unique index behavior.
func TestCheckAndMarkEvent_ConcurrentDuplicates(t *testing.T) {
	t.Skip("Requires MongoDB integration test for realistic concurrency testing")

	// Expected behavior (documented for integration test):
	// 1. Create 10 concurrent goroutines
	// 2. Each calls CheckAndMarkEvent with same eventID
	// 3. MongoDB unique index ensures only ONE succeeds in marking
	// 4. All others receive duplicate key error (handled gracefully)
	// 5. Verify no data corruption (race detector)
	//
	// Run with: go test ./internal/service -race -v
}

// TestCheckAndMarkEvent_MultipleEventsSequential tests that different events are processed independently
func TestCheckAndMarkEvent_MultipleEventsSequential(t *testing.T) {
	// Arrange
	ctx := testutil.NewTestContext()

	testCases := []struct {
		eventID   string
		eventType string
	}{
		{"event-001", "reservation.created"},
		{"event-002", "reservation.cancelled"},
		{"event-003", "reservation.created"},
	}

	mockRepo := new(MockEventRepository)

	// Configurar mocks para cada evento
	for _, tc := range testCases {
		mockRepo.On("IsEventProcessed", ctx, tc.eventID).Return(false, nil)
		mockRepo.On("MarkEventProcessed", ctx, mock.MatchedBy(func(event *domain.ProcessedEvent) bool {
			return event.EventID == tc.eventID && event.EventType == tc.eventType
		})).Return(nil)
	}

	service := NewIdempotencyService(mockRepo)

	// Act & Assert
	for _, tc := range testCases {
		shouldProcess, err := service.CheckAndMarkEvent(ctx, tc.eventID, tc.eventType)

		assert.NoError(t, err, "Event %s should not have error", tc.eventID)
		assert.True(t, shouldProcess, "Event %s should be processed", tc.eventID)
	}

	mockRepo.AssertExpectations(t)
}
