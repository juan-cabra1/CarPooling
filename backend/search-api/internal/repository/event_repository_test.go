package repository

import (
	"context"
	"testing"
	"time"

	"search-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupEventTest creates a test MongoDB connection and repository for event tests
func setupEventTest(t *testing.T) (EventRepository, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err, "Failed to connect to MongoDB")

	db := client.Database("search_api_test_events")

	// Create unique index on event_id for idempotency
	_, err = db.Collection("processed_events").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "event_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	require.NoError(t, err, "Failed to create unique index on event_id")

	repo := NewEventRepository(db)

	cleanup := func() {
		ctx := context.Background()
		_ = db.Collection("processed_events").Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return repo, cleanup
}

func TestEventRepository_IsEventProcessed_NotFound(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	// Check if event is processed (should be false - not found)
	isProcessed, err := repo.IsEventProcessed(context.Background(), "event-not-exists")

	require.NoError(t, err)
	assert.False(t, isProcessed, "Event should not be processed yet")
}

func TestEventRepository_MarkEventProcessed_Success(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	eventID := "event-123"
	event := &domain.ProcessedEvent{
		EventID:   eventID,
		EventType: "trip.created",
		Result:    "success",
	}

	// Mark event as processed
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err)

	// Verify it's now marked as processed
	isProcessed, err := repo.IsEventProcessed(context.Background(), eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Event should be marked as processed")
}

func TestEventRepository_MarkEventProcessed_Duplicate(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	eventID := "event-456"
	event := &domain.ProcessedEvent{
		EventID:   eventID,
		EventType: "trip.updated",
		Result:    "success",
	}

	// Mark event first time
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err)

	// Try to mark same event again - should succeed (idempotent)
	err = repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err, "Marking duplicate event should not return error (idempotent)")

	// Verify it's still marked as processed
	isProcessed, err := repo.IsEventProcessed(context.Background(), eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed)
}

func TestEventRepository_MarkEventProcessed_SetsTimestamp(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	event := &domain.ProcessedEvent{
		EventID:   "event-789",
		EventType: "trip.cancelled",
		Result:    "success",
		// ProcessedAt is intentionally zero - should be auto-set
	}

	// Mark event (should auto-set ProcessedAt)
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err)

	// Verify event was processed
	isProcessed, err := repo.IsEventProcessed(context.Background(), event.EventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Event should be marked as processed")
}

func TestEventRepository_IsEventProcessed_MultipleEvents(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	events := []*domain.ProcessedEvent{
		{EventID: "event-001", EventType: "trip.created", Result: "success"},
		{EventID: "event-002", EventType: "trip.updated", Result: "success"},
		{EventID: "event-003", EventType: "trip.cancelled", Result: "failed"},
	}

	// Mark all events as processed
	for _, event := range events {
		err := repo.MarkEventProcessed(context.Background(), event)
		require.NoError(t, err)
	}

	// Verify all are marked as processed
	for _, event := range events {
		isProcessed, err := repo.IsEventProcessed(context.Background(), event.EventID)
		require.NoError(t, err)
		assert.True(t, isProcessed, "Event %s should be processed", event.EventID)
	}

	// Verify a non-existent event is not processed
	isProcessed, err := repo.IsEventProcessed(context.Background(), "event-999")
	require.NoError(t, err)
	assert.False(t, isProcessed, "Non-existent event should not be processed")
}

func TestEventRepository_ConcurrentMarkEvent(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	eventID := "event-concurrent"
	event := &domain.ProcessedEvent{
		EventID:   eventID,
		EventType: "trip.created",
		Result:    "success",
	}

	// Simulate concurrent processing attempts
	errors := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			err := repo.MarkEventProcessed(context.Background(), event)
			errors <- err
		}()
	}

	// Collect results - all should succeed (idempotent)
	for i := 0; i < 5; i++ {
		err := <-errors
		require.NoError(t, err, "Concurrent mark should not return error (idempotent)")
	}

	// Verify event is marked exactly once
	isProcessed, err := repo.IsEventProcessed(context.Background(), eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Event should be marked as processed")
}

func TestEventRepository_MarkEventProcessed_WithMetadata(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	event := &domain.ProcessedEvent{
		EventID:      "event-metadata",
		EventType:    "trip.created",
		Result:       "success",
		ProcessedAt:  time.Now(),
		ErrorMessage: "",
	}

	// Mark event with metadata
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err)

	// Verify it's processed
	isProcessed, err := repo.IsEventProcessed(context.Background(), event.EventID)
	require.NoError(t, err)
	assert.True(t, isProcessed)
}
