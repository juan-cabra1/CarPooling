---
name: Bookings API - Phase 7
about: RabbitMQ Consumer with Idempotency (CRITICAL PHASE)
title: 'Phase 7: RabbitMQ Consumer - bookings-api [CRITICAL]'
labels: enhancement, bookings-api, phase-7, critical, idempotency
assignees: ''
---

## Phase 7: RabbitMQ Consumer (IDEMPOTENCY) ⚠️ CRITICAL

**Service:** bookings-api
**Priority:** CRITICAL 🔴
**Estimated Time:** 6-8 hours
**Branch:** `feature/bookings-api-phase-7-consumer`

---

### Description
Implement RabbitMQ consumer to handle trip events from trips-api. **THIS IS THE MOST CRITICAL PHASE** as it implements idempotency to prevent double-booking when RabbitMQ retries failed messages.

**Why Critical:**
- Prevents duplicate event processing
- Ensures data consistency across services
- Handles RabbitMQ retries safely
- Core of distributed system reliability

---

### Tasks
- [ ] Create `internal/messaging/rabbitmq.go` - RabbitMQ connection setup
- [ ] Create `internal/messaging/trip_consumer.go` with:
  - [ ] `TripEvent` struct matching event schema
  - [ ] `HandleTripEvent(event)` with idempotency check
  - [ ] `handleTripUpdated(event)` - confirm pending bookings
  - [ ] `handleTripCancelled(event)` - cancel all bookings for trip
  - [ ] `handleReservationFailed(event)` - mark booking as failed (compensation)
- [ ] Implement manual ACK/NACK logic:
  - [ ] ACK after successful processing
  - [ ] NACK with requeue on transient errors
  - [ ] NACK without requeue on permanent errors
- [ ] Add structured logging with zerolog for all events
- [ ] Setup consumer in main.go (run in goroutine)
- [ ] Declare queue `trips.events` with durable=true

---

### Success Criteria
- [ ] Consumer connects to RabbitMQ successfully
- [ ] **Duplicate events are skipped (idempotency works)**
- [ ] Events are ACKed only after successful processing
- [ ] Retries work on failures (DB temporarily down)
- [ ] trip.updated → pending bookings become confirmed
- [ ] trip.cancelled → all bookings for that trip are cancelled
- [ ] reservation.failed → booking marked as failed
- [ ] All event processing logged with event_id
- [ ] No duplicate bookings created from retried events
- [ ] Consumer runs in background without blocking HTTP server

---

### Dependencies
**Requires:**
- Phase 4 completed (IdempotencyService)
- Phase 3 completed (BookingRepository with FindByTripID)
- RabbitMQ running and accessible
- `processed_events` table with UNIQUE constraint on event_id

**Blocks:**
- Full end-to-end testing
- Integration with trips-api

---

### Files to Create/Modify
```
backend/bookings-api/
├── internal/messaging/
│   ├── rabbitmq.go          # Connection setup
│   └── trip_consumer.go     # Consumer implementation
└── cmd/api/main.go          # Start consumer in goroutine
```

---

### Event Schema (from trips-api)
```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "trip.updated",
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 2,
  "status": "published",
  "timestamp": "2024-11-10T10:00:00Z"
}
```

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440001",
  "event_type": "trip.cancelled",
  "trip_id": "507f1f77bcf86cd799439011",
  "cancelled_by": 123,
  "reason": "Car broke down",
  "timestamp": "2024-11-10T11:00:00Z"
}
```

---

### Idempotency Implementation (CRITICAL)
```go
func (c *tripConsumer) HandleTripEvent(delivery amqp.Delivery) error {
    var event TripEvent
    if err := json.Unmarshal(delivery.Body, &event); err != nil {
        c.logger.Error().Err(err).Msg("Failed to unmarshal event")
        return delivery.Nack(false, false) // Don't requeue invalid JSON
    }

    // 1. CHECK IDEMPOTENCY FIRST
    shouldProcess, err := c.idempotencyService.CheckAndMarkEvent(
        event.EventID,
        event.EventType,
    )
    if err != nil {
        c.logger.Error().Err(err).Str("event_id", event.EventID).Msg("Idempotency check failed")
        return delivery.Nack(false, true) // Requeue - transient error
    }

    if !shouldProcess {
        // Event already processed - skip
        c.logger.Info().Str("event_id", event.EventID).Msg("Event already processed, skipping")
        return delivery.Ack(false) // ACK without processing
    }

    // 2. PROCESS EVENT
    var processingErr error
    switch event.EventType {
    case "trip.updated":
        processingErr = c.handleTripUpdated(event)
    case "trip.cancelled":
        processingErr = c.handleTripCancelled(event)
    case "reservation.failed":
        processingErr = c.handleReservationFailed(event)
    default:
        c.logger.Warn().Str("event_type", event.EventType).Msg("Unknown event type")
        return delivery.Ack(false) // ACK unknown events
    }

    if processingErr != nil {
        c.logger.Error().Err(processingErr).Str("event_id", event.EventID).Msg("Event processing failed")
        return delivery.Nack(false, true) // Requeue - will retry
    }

    c.logger.Info().Str("event_id", event.EventID).Str("event_type", event.EventType).Msg("Event processed successfully")
    return delivery.Ack(false)
}
```

---

### Testing Steps

#### 1. Unit Test Idempotency (Most Important)
```bash
cd backend/bookings-api
go test ./internal/service/idempotency_service_test.go -v

# Should verify:
# - First event is processed
# - Duplicate event is skipped
# - Concurrent duplicates only process once
```

#### 2. Manual Testing with RabbitMQ
```bash
# Terminal 1: Start bookings-api
cd backend/bookings-api
go run cmd/api/main.go

# Terminal 2: Publish test event to RabbitMQ
# (Need rabbitmqadmin or similar tool)
rabbitmqadmin publish exchange=trips.events routing_key="" payload='
{
  "event_id": "test-event-123",
  "event_type": "trip.updated",
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 2,
  "status": "published",
  "timestamp": "2024-11-10T10:00:00Z"
}'

# Check logs:
# Should see: "Event processed successfully" with event_id=test-event-123

# Publish SAME event again
rabbitmqadmin publish exchange=trips.events routing_key="" payload='...'

# Check logs:
# Should see: "Event already processed, skipping" with event_id=test-event-123

# Verify in MySQL:
mysql -u root -p
USE carpooling_bookings;
SELECT * FROM processed_events WHERE event_id='test-event-123';
# Should show exactly 1 row
```

#### 3. Integration Test with trips-api
```bash
# Once trips-api is implemented:
# 1. Create a trip
# 2. Create a booking for that trip
# 3. trips-api publishes trip.updated
# 4. Verify booking status changes to 'confirmed'
# 5. Cancel trip in trips-api
# 6. Verify all bookings for that trip are cancelled
```

---

### Common Pitfalls to Avoid

❌ **DON'T:**
- Process event before checking idempotency
- ACK before successful processing
- Requeue permanent errors (invalid JSON)
- Skip logging event_id
- Use auto-ACK mode

✅ **DO:**
- Check idempotency FIRST
- ACK only after success
- NACK with requeue for transient errors
- Log every event_id
- Use manual ACK mode
- Test duplicate event handling

---

### Implementation Guide
1. Create branch: `git checkout -b feature/bookings-api-phase-7-consumer`
2. Use Claude Code plan mode:
   ```
   @CONTEXT_BOOKINGS_API.md Implement Phase 7 (RabbitMQ Consumer) - CRITICAL PHASE

   This is the MOST IMPORTANT phase. Implement:

   1. internal/messaging/rabbitmq.go - connection setup
   2. internal/messaging/trip_consumer.go with:
      - HandleTripEvent with IDEMPOTENCY check
      - handleTripUpdated, handleTripCancelled, handleReservationFailed
      - Manual ACK/NACK logic
      - Structured logging

   Critical requirements:
   - Check idempotency BEFORE processing (CheckAndMarkEvent)
   - Skip if already processed
   - ACK only after success
   - NACK with requeue on transient errors

   Event format and idempotency logic in CONTEXT_BOOKINGS_API.md
   ```
3. **REVIEW PLAN CAREFULLY** - this is critical code
4. Implement
5. **TEST THOROUGHLY** - especially duplicate events
6. Verify all success criteria
7. Commit: `git commit -m "feat(bookings): implement RabbitMQ consumer with idempotency"`
8. Push and create PR

---

### References
- Implementation Plan: `GITFLOW.md` - Phase 7
- Context: `CONTEXT_BOOKINGS_API.md` - "RabbitMQ Consumer" section
- How to Use: `HOW_TO_USE_PLAN_MODE.md`
- RabbitMQ Docs: https://www.rabbitmq.com/tutorials/tutorial-two-go.html

---

### Notes
- **Take your time with this phase** - it's the foundation of system reliability
- Test idempotency thoroughly - use the race condition test
- Make sure `processed_events` table has UNIQUE constraint on `event_id`
- Log with structured fields: `Str("event_id", event.EventID)`
- Run consumer in goroutine in main.go: `go consumer.Start("trips.events")`
- Handle graceful shutdown (defer channel close)

---

### Acceptance Criteria for PR Review
Before approving this PR, verify:
- [ ] Duplicate events are not processed twice (test with same event_id)
- [ ] Events are ACKed only after successful DB update
- [ ] Idempotency test passes (including race condition test)
- [ ] All event types handled (trip.updated, trip.cancelled, reservation.failed)
- [ ] Structured logging shows event_id for every event
- [ ] Consumer runs without blocking HTTP server
- [ ] Code handles DB connection failures gracefully (NACK with requeue)
