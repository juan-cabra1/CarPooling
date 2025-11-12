# Future Enhancements for Bookings API

This document outlines identified improvements and missing features that should be considered for production readiness and the final presentation.

## 🚨 Critical: Missing Trip Lifecycle Events

### Current Implementation Status

The bookings-api RabbitMQ consumer currently handles **only 2 events**:
- ✅ `trip.cancelled` - Implemented
- ✅ `reservation.failed` - Implemented

### Missing Events (Requires trips-api Changes)

#### 1. `trip.completed` Event

**Problem:**
- Trip status `completed` exists in trips-api domain model but is **never set**
- Booking status `completed` exists but **cannot be reached**
- No mechanism to transition bookings from `confirmed` → `completed`

**Impact:**
- Bookings remain in `confirmed` status forever, even after trip completion
- No way to query historical/completed bookings accurately
- Business analytics incomplete (can't distinguish active vs completed bookings)

**Recommended Solution:**
1. Add endpoint to trips-api: `PATCH /trips/:id/complete`
2. Implement automatic trip completion based on `estimated_arrival_datetime`
3. Publish `trip.completed` event when status changes
4. Update bookings-api consumer to handle event and transition all confirmed bookings to `completed`

**Event Structure (Proposed):**
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.completed",
  "trip_id": "mongodb-objectid",
  "driver_id": 123,
  "status": "completed",
  "completed_at": "2024-11-10T15:00:00Z",
  "source_service": "trips-api",
  "correlation_id": "uuid-v4",
  "timestamp": "2024-11-10T15:00:05Z"
}
```

**Handler Logic (Proposed):**
```go
func (c *TripsConsumer) HandleTripCompleted(body []byte) error {
    // 1. Check idempotency
    // 2. Find all bookings with trip_id and status='confirmed'
    // 3. Update ALL to status='completed'
    // 4. Log completion
    // 5. ACK message
}
```

---

#### 2. `trip.deleted` Event

**Problem:**
- Trips can be hard deleted via `DELETE /trips/:id` endpoint
- **NO event is published** when deletion occurs
- Bookings-api has no way to know a trip was deleted
- Results in **orphaned bookings** (trip_id references non-existent trip)

**Impact:**
- Data integrity violation (foreign key references dangling)
- Bookings reference non-existent trips
- Users can have confirmed bookings for deleted trips
- No audit trail for deleted trips affecting bookings

**Current Workaround:**
- Using trip cancellation instead of deletion (soft delete via status)
- This maintains referential integrity but doesn't solve the underlying issue

**Recommended Solution:**

**Option A: Implement trip.deleted Event (RECOMMENDED)**
1. Modify trips-api `DeleteTrip()` service method
2. Publish `trip.deleted` event BEFORE deleting from database
3. Update bookings-api consumer to:
   - Cancel all bookings for deleted trip
   - Publish `reservation.cancelled` events to free seats (if needed)
   - Log deletion for audit trail

**Option B: Soft Delete (PREFERRED FOR PRODUCTION)**
1. Replace hard delete with soft delete in trips-api
2. Add `deleted_at` timestamp field to Trip model
3. Filter deleted trips from queries using `deleted_at IS NULL`
4. Publish `trip.cancelled` event on soft delete
5. Maintain full audit trail and data integrity

**Event Structure (Proposed for Option A):**
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.deleted",
  "trip_id": "mongodb-objectid",
  "driver_id": 123,
  "deleted_by": 123,
  "deletion_reason": "Driver request",
  "source_service": "trips-api",
  "correlation_id": "uuid-v4",
  "timestamp": "2024-11-10T16:00:00Z"
}
```

**Why Soft Delete is Better:**
- ✅ Maintains referential integrity
- ✅ Preserves audit trail
- ✅ Allows data recovery if needed
- ✅ Supports compliance requirements (GDPR, data retention)
- ✅ Enables analytics on deleted trips
- ✅ Prevents accidental data loss

---

## 📝 Architecture Decisions

### Why `trip.updated` is NOT Consumed

**Decision:** Bookings-api does **not** consume `trip.updated` events.

**Rationale:**
- **Bookings are immutable contracts** between passenger and driver
- Booking captures trip snapshot at reservation time (price, destination, datetime)
- Allowing drivers to modify trip details after booking breaks passenger trust
- Price manipulation prevention: Driver cannot increase price after booking

**Business Rules:**
- If driver changes destination → Must cancel and recreate trip
- If driver changes datetime → Passengers must accept change (future feature)
- If driver changes price → Existing bookings maintain original price
- Seat availability updates → Handled via `reservation.created/cancelled` events

**Current Event Flow:**
```
1. User creates booking → Publishes reservation.created
2. Trips-api reserves seats → Publishes trip.updated (NOT consumed by bookings-api)
3. Booking status remains as created (pending/confirmed based on initial response)
```

**Alternative Considered:**
Consuming `trip.updated` to transition `pending` → `confirmed` was considered but rejected because:
- Bookings should be created as `confirmed` directly (atomic operation)
- Pending state only needed if reservation is asynchronous
- Current implementation may need review to ensure bookings are created as `confirmed`

---

## 🔧 Technical Debt & Improvements

### 1. Booking Creation Flow Review

**Current State:** Unclear if bookings are created as `pending` or `confirmed`

**Needs Investigation:**
- Review booking creation in `booking_service.go`
- Confirm if `reservation.created` is published synchronously
- Determine if pending state is necessary
- Consider making booking creation atomic with trip reservation

### 2. Bulk Update Optimization

**Current Implementation:**
```go
// Handler loops through bookings and cancels individually
for _, booking := range confirmedBookings {
    err := c.bookingRepo.CancelBooking(booking.BookingUUID, reason)
}
```

**Potential Improvement:**
Add bulk update method to repository:
```go
func (r *bookingRepository) CancelBookingsByTripID(tripID, reason string) error {
    // Single SQL UPDATE statement for all bookings
    // More efficient for trips with many bookings
}
```

**Trade-off:**
- ✅ Better performance for trips with many bookings
- ❌ Less granular error handling
- ❌ All-or-nothing approach (transaction required)

**Recommendation:** Implement for production if average bookings per trip > 10

### 3. Dead Letter Queue (DLQ)

**Current State:**
- Failed messages are NACKed and requeued
- No limit on retry attempts
- Malformed messages are ACKed and discarded

**Improvement:**
Configure RabbitMQ Dead Letter Exchange:
```go
// Queue arguments
args := amqp.Table{
    "x-dead-letter-exchange": "trips.dlx",
    "x-dead-letter-routing-key": "failed",
    "x-message-ttl": 86400000, // 24 hours
}
```

**Benefits:**
- Prevent infinite retry loops
- Analyze failed messages for debugging
- Manual reprocessing capability
- Monitoring and alerting on DLQ depth

### 4. Consumer Health Check

**Current State:**
- HTTP `/health` endpoint doesn't check consumer status
- No visibility into consumer connection health

**Improvement:**
Add consumer status to health check:
```go
type HealthResponse struct {
    Service   string `json:"service"`
    Status    string `json:"status"`
    Consumer  ConsumerHealth `json:"consumer"`
}

type ConsumerHealth struct {
    Connected bool   `json:"connected"`
    Queue     string `json:"queue"`
    Messages  int    `json:"pending_messages"`
}
```

### 5. Metrics & Monitoring

**Missing:**
- Event processing duration metrics
- Event type counters
- Error rate tracking
- Queue depth monitoring

**Recommended Tools:**
- Prometheus for metrics collection
- Grafana for visualization
- Alert on consumer errors > threshold

---

## 📊 Production Readiness Checklist

Before deploying to production:

### Events
- [ ] Implement `trip.completed` event in trips-api
- [ ] Implement `trip.completed` handler in bookings-api
- [ ] Implement `trip.deleted` event OR soft delete in trips-api
- [ ] Add comprehensive event integration tests

### Reliability
- [ ] Configure Dead Letter Queue
- [ ] Add retry limits and backoff strategy
- [ ] Implement circuit breaker for external dependencies
- [ ] Add consumer health check endpoint

### Observability
- [ ] Add Prometheus metrics
- [ ] Set up Grafana dashboards
- [ ] Configure alerts for consumer errors
- [ ] Add distributed tracing (OpenTelemetry)

### Data Integrity
- [ ] Review booking creation flow (pending vs confirmed)
- [ ] Add database constraints for referential integrity
- [ ] Implement soft delete for trips (recommended)
- [ ] Add data migration for existing bookings

### Testing
- [ ] Integration tests for all event handlers
- [ ] Load testing for concurrent message processing
- [ ] Chaos engineering (RabbitMQ failures, DB failures)
- [ ] Idempotency testing (duplicate events)

### Documentation
- [ ] Update API documentation
- [ ] Add event schema documentation
- [ ] Create runbook for common issues
- [ ] Document rollback procedures

---

## 🎯 Priority for Final Presentation

### High Priority (Must Have)
1. ✅ Document missing events and their impact
2. ✅ Explain booking immutability decision
3. ⚠️ Implement `trip.completed` lifecycle (if time permits)
4. ⚠️ Add basic consumer health check

### Medium Priority (Nice to Have)
1. Dead Letter Queue configuration
2. Metrics and monitoring setup
3. Integration tests for event handlers

### Low Priority (Future Work)
1. Bulk update optimization
2. Circuit breaker implementation
3. Soft delete migration

---

## 📚 References

- **CONTEXT_BOOKINGS_API.md** - Bookings API specification
- **trips-api/internal/messaging/** - Event publisher implementation
- **trips-api/internal/domain/trip.go** - Trip status definitions
- **RabbitMQ Best Practices** - https://www.rabbitmq.com/best-practices.html
- **Idempotency Patterns** - https://www.enterpriseintegrationpatterns.com/patterns/messaging/IdempotentReceiver.html

---

**Last Updated:** 2025-11-10
**Status:** Ready for review and presentation planning
