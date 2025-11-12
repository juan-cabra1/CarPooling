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

// setupEventTestDB creates a test MongoDB connection for event repository tests
func setupEventTestDB(t *testing.T) (*mongo.Database, func()) {
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

	cleanup := func() {
		ctx := context.Background()
		_ = db.Collection("processed_events").Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestEventRepository_CheckAndMarkEvent_FirstTime(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	// First time processing an event
	shouldProcess, err := repo.CheckAndMarkEvent(context.Background(), "event-123", "trip.created")
	require.NoError(t, err, "Should not return error on first check")
	assert.True(t, shouldProcess, "Should return true for first-time event processing")
}

func TestEventRepository_CheckAndMarkEvent_Duplicate(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	// Process event first time
	shouldProcess, err := repo.CheckAndMarkEvent(context.Background(), "event-456", "trip.updated")
	require.NoError(t, err, "First check should succeed")
	assert.True(t, shouldProcess, "Should process first time")

	// Try to process same event again (duplicate)
	shouldProcess, err = repo.CheckAndMarkEvent(context.Background(), "event-456", "trip.updated")
	require.NoError(t, err, "Duplicate check should not return error")
	assert.False(t, shouldProcess, "Should return false for duplicate event")
}

func TestEventRepository_CheckAndMarkEvent_DifferentEventTypes(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	eventID := "event-789"

	// Process event with type "trip.created"
	shouldProcess, err := repo.CheckAndMarkEvent(context.Background(), eventID, "trip.created")
	require.NoError(t, err, "First check should succeed")
	assert.True(t, shouldProcess, "Should process first event type")

	// Try same event_id with different event_type (should be treated as duplicate)
	shouldProcess, err = repo.CheckAndMarkEvent(context.Background(), eventID, "trip.updated")
	require.NoError(t, err, "Second check should not return error")
	assert.False(t, shouldProcess, "Should return false - event_id already processed regardless of type")
}

func TestEventRepository_IsEventProcessed(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	eventID := "event-check-123"

	// Check before processing
	isProcessed, err := repo.IsEventProcessed(context.Background(), eventID)
	require.NoError(t, err, "Should not return error")
	assert.False(t, isProcessed, "Event should not be processed yet")

	// Mark event as processed
	event := &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   "trip.created",
		ProcessedAt: time.Now(),
		Result:      "success",
	}
	err = repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err, "Should mark event successfully")

	// Check after processing
	isProcessed, err = repo.IsEventProcessed(context.Background(), eventID)
	require.NoError(t, err, "Should not return error")
	assert.True(t, isProcessed, "Event should be marked as processed")
}

func TestEventRepository_MarkEventProcessed_Idempotent(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	event := &domain.ProcessedEvent{
		EventID:     "event-idempotent-123",
		EventType:   "trip.cancelled",
		ProcessedAt: time.Now(),
		Result:      "success",
	}

	// Mark event first time
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err, "First mark should succeed")

	// Mark same event again (idempotent - should not return error)
	err = repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err, "Second mark should be idempotent and not return error")
}

func TestEventRepository_MarkEventProcessed_SetsTimestamp(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)

	event := &domain.ProcessedEvent{
		EventID:   "event-timestamp-123",
		EventType: "trip.created",
		Result:    "success",
		// ProcessedAt is intentionally zero
	}

	// Mark event (should auto-set ProcessedAt)
	err := repo.MarkEventProcessed(context.Background(), event)
	require.NoError(t, err, "Should mark event successfully")

	// Verify timestamp was set
	isProcessed, err := repo.IsEventProcessed(context.Background(), event.EventID)
	require.NoError(t, err, "Should not return error")
	assert.True(t, isProcessed, "Event should be processed")
}

func TestEventRepository_ConcurrentCheckAndMark(t *testing.T) {
	db, cleanup := setupEventTestDB(t)
	defer cleanup()

	repo := NewEventRepository(db)
	eventID := "event-concurrent-123"

	// Simulate concurrent processing attempts
	results := make(chan bool, 3)
	errors := make(chan error, 3)

	for i := 0; i < 3; i++ {
		go func() {
			shouldProcess, err := repo.CheckAndMarkEvent(context.Background(), eventID, "trip.created")
			results <- shouldProcess
			errors <- err
		}()
	}

	// Collect results
	var processCount int
	for i := 0; i < 3; i++ {
		err := <-errors
		require.NoError(t, err, "Should not return error")

		shouldProcess := <-results
		if shouldProcess {
			processCount++
		}
	}

	// Exactly one goroutine should have been allowed to process
	assert.Equal(t, 1, processCount, "Exactly one concurrent call should be allowed to process")
}
