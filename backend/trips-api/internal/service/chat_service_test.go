package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"trips-api/internal/dao"
	"trips-api/internal/domain"
)

// Mock repositories for testing

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *dao.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByTripID(ctx context.Context, tripID string, limit int) ([]*dao.Message, error) {
	args := m.Called(ctx, tripID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dao.Message), args.Error(1)
}

type MockTripRepositoryForChat struct {
	mock.Mock
}

func (m *MockTripRepositoryForChat) FindByID(ctx context.Context, id string) (*domain.Trip, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trip), args.Error(1)
}

func (m *MockTripRepositoryForChat) Create(ctx context.Context, trip *domain.Trip) error {
	args := m.Called(ctx, trip)
	return args.Error(0)
}

func (m *MockTripRepositoryForChat) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error) {
	args := m.Called(ctx, filters, page, limit)
	return args.Get(0).([]domain.Trip), args.Get(1).(int64), args.Error(2)
}

func (m *MockTripRepositoryForChat) Update(ctx context.Context, id string, trip *domain.Trip) error {
	args := m.Called(ctx, id, trip)
	return args.Error(0)
}

func (m *MockTripRepositoryForChat) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTripRepositoryForChat) UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error {
	args := m.Called(ctx, tripID, seatsDelta, expectedVersion)
	return args.Error(0)
}

func (m *MockTripRepositoryForChat) Cancel(ctx context.Context, id string, cancelledBy int64, reason string) error {
	args := m.Called(ctx, id, cancelledBy, reason)
	return args.Error(0)
}

func (m *MockTripRepositoryForChat) UpdateLastActivity(ctx context.Context, tripID string, timestamp time.Time) error {
	args := m.Called(ctx, tripID, timestamp)
	return args.Error(0)
}

type MockPublisherForChat struct {
	mock.Mock
}

func (m *MockPublisherForChat) PublishChatMessage(tripID string, userID int64, message string) error {
	args := m.Called(tripID, userID, message)
	return args.Error(0)
}

// ═══════════════════════════════════════════════════════════════════════════
// TESTS: Concurrent Processing Verification
// ═══════════════════════════════════════════════════════════════════════════

// TestSendMessage_Success verifies that concurrent processing works correctly
// This test demonstrates the use of Goroutines + Channels + WaitGroup
func TestSendMessage_Success(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	// Mock all concurrent operations to succeed
	mockMessageRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockTripRepo.On("FindByID", mock.Anything, "trip-123").Return(&domain.Trip{}, nil)
	mockPublisher.On("PublishChatMessage", "trip-123", int64(1), "Hello").Return(nil)
	mockTripRepo.On("UpdateLastActivity", mock.Anything, "trip-123", mock.Anything).Return(nil)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	message, err := service.SendMessage(context.Background(), "trip-123", 1, "Test User", "Hello")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "trip-123", message.TripID)
	assert.Equal(t, int64(1), message.UserID)
	assert.Equal(t, "Test User", message.UserName)
	assert.Equal(t, "Hello", message.Message)

	// Verify all concurrent operations were called
	mockMessageRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

// TestSendMessage_EmptyMessage verifies validation
func TestSendMessage_EmptyMessage(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	message, err := service.SendMessage(context.Background(), "trip-123", 1, "Test User", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Equal(t, "message cannot be empty", err.Error())

	// No repository calls should have been made
	mockMessageRepo.AssertNotCalled(t, "Create")
	mockTripRepo.AssertNotCalled(t, "FindByID")
}

// TestSendMessage_TripNotFound verifies error handling for invalid trip
func TestSendMessage_TripNotFound(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	// Mock trip verification to fail (trip not found)
	mockMessageRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockTripRepo.On("FindByID", mock.Anything, "invalid-trip").Return(nil, errors.New("trip not found"))
	mockPublisher.On("PublishChatMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockTripRepo.On("UpdateLastActivity", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	message, err := service.SendMessage(context.Background(), "invalid-trip", 1, "Test User", "Hello")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Contains(t, err.Error(), "trip not found")

	// Verify trip verification was called
	mockTripRepo.AssertExpectations(t)
}

// TestSendMessage_DatabaseError verifies error handling for DB failures
func TestSendMessage_DatabaseError(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	// Mock message save to fail
	mockMessageRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
	mockTripRepo.On("FindByID", mock.Anything, "trip-123").Return(&domain.Trip{}, nil)
	mockPublisher.On("PublishChatMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockTripRepo.On("UpdateLastActivity", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	message, err := service.SendMessage(context.Background(), "trip-123", 1, "Test User", "Hello")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, message)

	// Verify critical operations were attempted
	mockMessageRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
}

// TestSendMessage_NonCriticalFailureDoesNotBreak verifies eventual consistency
func TestSendMessage_NonCriticalFailureDoesNotBreak(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	// Critical operations succeed, non-critical fails
	mockMessageRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockTripRepo.On("FindByID", mock.Anything, "trip-123").Return(&domain.Trip{}, nil)
	mockPublisher.On("PublishChatMessage", "trip-123", int64(1), "Hello").Return(errors.New("rabbitmq down"))
	mockTripRepo.On("UpdateLastActivity", mock.Anything, "trip-123", mock.Anything).Return(errors.New("update failed"))

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	message, err := service.SendMessage(context.Background(), "trip-123", 1, "Test User", "Hello")

	// Assert - should succeed despite non-critical failures
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "trip-123", message.TripID)

	// All operations were attempted
	mockMessageRepo.AssertExpectations(t)
	mockTripRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

// TestGetMessages_Success verifies message retrieval
func TestGetMessages_Success(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	expectedMessages := []*dao.Message{
		{TripID: "trip-123", UserID: 1, UserName: "User 1", Message: "Hello"},
		{TripID: "trip-123", UserID: 2, UserName: "User 2", Message: "Hi there"},
	}

	mockMessageRepo.On("FindByTripID", mock.Anything, "trip-123", 50).Return(expectedMessages, nil)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	messages, err := service.GetMessages(context.Background(), "trip-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, messages, 2)
	assert.Equal(t, "Hello", messages[0].Message)
	assert.Equal(t, "Hi there", messages[1].Message)

	mockMessageRepo.AssertExpectations(t)
}

// TestGetMessages_Empty verifies empty result handling
func TestGetMessages_Empty(t *testing.T) {
	// Arrange
	mockMessageRepo := new(MockMessageRepository)
	mockTripRepo := new(MockTripRepositoryForChat)
	mockPublisher := new(MockPublisherForChat)

	mockMessageRepo.On("FindByTripID", mock.Anything, "trip-123", 50).Return([]*dao.Message{}, nil)

	service := NewChatService(mockMessageRepo, mockTripRepo, mockPublisher)

	// Act
	messages, err := service.GetMessages(context.Background(), "trip-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, messages, 0)

	mockMessageRepo.AssertExpectations(t)
}
