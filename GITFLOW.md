# GitFlow Strategy - CarPooling Project

## Branch Strategy

### Main Branches
- `main` - Production-ready code (entrega final)
- `dev` - Integration branch (desarrollo activo)

### Feature Branches
- `feature/bookings-api-phase-{N}` - Para cada fase de implementación
- `feature/trips-api-phase-{N}`
- `feature/search-api-phase-{N}`
- Naming: `feature/{service}-{phase-name}`

### Workflow
```
main (protected)
  └── dev (base para development)
       ├── feature/bookings-api-phase-1-setup
       ├── feature/bookings-api-phase-2-database
       ├── feature/bookings-api-phase-3-repository
       └── ...
```

---

## Development Workflow

### 1. Start New Phase
```bash
# Always start from updated dev
git checkout dev
git pull origin dev

# Create feature branch
git checkout -b feature/bookings-api-phase-1-setup
```

### 2. Work on Phase
```bash
# Make changes following the plan
# Commit frequently with clear messages
git add .
git commit -m "feat(bookings): setup project structure and dependencies"
```

### 3. Push and Create PR
```bash
git push -u origin feature/bookings-api-phase-1-setup

# Create PR: feature/bookings-api-phase-1-setup → dev
# Request review from team
```

### 4. Merge and Next Phase
```bash
# After PR approved
git checkout dev
git pull origin dev
git checkout -b feature/bookings-api-phase-2-database
```

---

## Commit Message Convention

Format: `type(scope): description`

### Types
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `test`: Adding tests
- `docs`: Documentation
- `chore`: Build/config changes

### Examples
```
feat(bookings): add booking repository with GORM
feat(bookings): implement idempotency service
feat(bookings): add RabbitMQ consumer for trip events
test(bookings): add idempotency service tests
refactor(bookings): extract error types to domain package
docs(bookings): add API documentation
```

---

## Issues Structure

Cada fase del plan = 1 issue en GitHub

### Issue Template
```markdown
## Phase N: [Phase Name]

**Service:** bookings-api
**Priority:** High/Medium/Low
**Estimated Time:** X hours
**Branch:** feature/bookings-api-phase-N-{name}

### Description
[Brief description from implementation plan]

### Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

### Success Criteria
- [ ] Criteria 1
- [ ] Criteria 2

### Dependencies
- Requires: Phase X completed
- Blocks: Phase Y

### Files to Create/Modify
- `path/to/file1.go`
- `path/to/file2.go`

### Testing
- [ ] Unit tests pass
- [ ] Manual testing completed

### References
- Implementation Plan: Phase N section
```

---

## Phases Breakdown for bookings-api

### Week 1: Foundation

**Issue #1: Phase 1 - Project Setup**
```markdown
**Branch:** feature/bookings-api-phase-1-setup
**Tasks:**
- [ ] Initialize go.mod
- [ ] Create .env and .env.example
- [ ] Setup main.go skeleton
- [ ] Add all dependencies (gin, gorm, uuid, zerolog, amqp)
- [ ] Verify compilation

**Success Criteria:**
- [ ] `go mod tidy` succeeds
- [ ] `go build` succeeds
- [ ] Server starts on port 8003
```

**Issue #2: Phase 2 - Database Setup**
```markdown
**Branch:** feature/bookings-api-phase-2-database
**Tasks:**
- [ ] Create config package
- [ ] Create BookingDAO with GORM tags
- [ ] Create ProcessedEventDAO with GORM tags
- [ ] Create domain DTOs (BookingDTO, CreateBookingRequest, etc.)
- [ ] Setup MySQL connection in main.go
- [ ] Add auto-migration

**Success Criteria:**
- [ ] MySQL connection successful
- [ ] Tables created with proper schema
- [ ] UNIQUE index on processed_events.event_id
- [ ] All indexes created correctly

**Files:**
- internal/config/config.go
- internal/dao/booking.go
- internal/dao/processed_event.go
- internal/domain/booking.go
```

**Issue #3: Phase 3 - Repository Layer**
```markdown
**Branch:** feature/bookings-api-phase-3-repository
**Tasks:**
- [ ] Create BookingRepository interface
- [ ] Implement bookingRepository with GORM
- [ ] Create EventRepository interface
- [ ] Implement eventRepository with idempotency logic
- [ ] Add all CRUD methods

**Success Criteria:**
- [ ] All repository methods compile
- [ ] GORM queries are correct
- [ ] Idempotency check uses UNIQUE constraint

**Files:**
- internal/repository/booking_repository.go
- internal/repository/event_repository.go
```

**Issue #4: Phase 4 - Business Logic**
```markdown
**Branch:** feature/bookings-api-phase-4-services
**Tasks:**
- [ ] Create BookingService interface
- [ ] Implement CreateBooking with validations
- [ ] Implement GetBooking, ListBookings, CancelBooking
- [ ] Create IdempotencyService
- [ ] Implement CheckAndMarkEvent
- [ ] Add error handling

**Success Criteria:**
- [ ] All business validations implemented
- [ ] Duplicate booking prevention works
- [ ] Seat validation works
- [ ] Trip snapshot captured

**Files:**
- internal/service/booking_service.go
- internal/service/idempotency_service.go
```

**Issue #5: Phase 5 - HTTP Controllers**
```markdown
**Branch:** feature/bookings-api-phase-5-controllers
**Tasks:**
- [ ] Create BookingController interface
- [ ] Implement CreateBooking handler
- [ ] Implement GetBooking handler
- [ ] Implement ListBookings handler
- [ ] Implement CancelBooking handler
- [ ] Implement ConfirmBooking handler
- [ ] Add proper status codes and error responses

**Success Criteria:**
- [ ] All endpoints respond correctly
- [ ] JSON validation works
- [ ] Error messages are clear

**Files:**
- internal/controller/booking_controller.go
```

**Issue #6: Phase 6 - JWT Middleware**
```markdown
**Branch:** feature/bookings-api-phase-6-middleware
**Tasks:**
- [ ] Copy auth.go from users-api
- [ ] Copy cors.go from users-api
- [ ] Copy error.go from users-api
- [ ] Adjust imports for bookings-api
- [ ] Test JWT validation

**Success Criteria:**
- [ ] Protected routes require JWT
- [ ] user_id extracted from token
- [ ] Invalid tokens rejected

**Files:**
- internal/middleware/auth.go
- internal/middleware/cors.go
- internal/middleware/error.go
```

### Week 2: Integration (CRITICAL)

**Issue #7: Phase 7 - RabbitMQ Consumer (IDEMPOTENCY)**
```markdown
**Branch:** feature/bookings-api-phase-7-consumer
**Priority:** CRITICAL
**Tasks:**
- [ ] Create rabbitmq.go connection setup
- [ ] Create TripEvent struct
- [ ] Implement trip_consumer.go
- [ ] Add idempotency check in HandleTripEvent
- [ ] Implement handleTripUpdated
- [ ] Implement handleTripCancelled
- [ ] Implement handleReservationFailed
- [ ] Setup manual ACK/NACK
- [ ] Add structured logging

**Success Criteria:**
- [ ] Consumer connects to RabbitMQ
- [ ] Duplicate events are skipped
- [ ] Events are ACKed only after processing
- [ ] Retries work on failures
- [ ] Idempotency test passes

**Files:**
- internal/messaging/rabbitmq.go
- internal/messaging/trip_consumer.go
```

**Issue #8: Phase 8 - RabbitMQ Publisher**
```markdown
**Branch:** feature/bookings-api-phase-8-publisher
**Tasks:**
- [ ] Create ReservationEvent struct
- [ ] Implement reservation_publisher.go
- [ ] Generate UUID for event_id
- [ ] Implement PublishReservationCreated
- [ ] Implement PublishReservationCancelled
- [ ] Implement PublishReservationFailed
- [ ] Set persistent delivery mode

**Success Criteria:**
- [ ] Events published to correct exchange
- [ ] event_id is UUID v4
- [ ] Messages are persistent
- [ ] Logging on publish success/failure

**Files:**
- internal/messaging/reservation_publisher.go
```

**Issue #9: Phase 9 - HTTP Clients**
```markdown
**Branch:** feature/bookings-api-phase-9-clients
**Tasks:**
- [ ] Create TripDTO struct
- [ ] Implement trips_client.go
- [ ] Add timeout (5s)
- [ ] Add retry logic (2 attempts)
- [ ] Create UserDTO struct
- [ ] Implement users_client.go
- [ ] Add error handling

**Success Criteria:**
- [ ] Can fetch trip from trips-api
- [ ] Can fetch user from users-api
- [ ] Timeouts work correctly
- [ ] Retries work on transient errors

**Files:**
- internal/clients/trips_client.go
- internal/clients/users_client.go
```

**Issue #10: Phase 10 - Error Handling & Logging**
```markdown
**Branch:** feature/bookings-api-phase-10-errors
**Tasks:**
- [ ] Create AppError type
- [ ] Define error constants
- [ ] Setup zerolog in main.go
- [ ] Add logging to all services
- [ ] Add logging to consumer
- [ ] Add logging to publisher

**Success Criteria:**
- [ ] All errors have codes
- [ ] Structured logs in JSON format
- [ ] Key events logged (booking created, event processed, etc.)

**Files:**
- internal/domain/errors.go
```

**Issue #11: Phase 11 - Routes Setup**
```markdown
**Branch:** feature/bookings-api-phase-11-routes
**Tasks:**
- [ ] Create routes.go
- [ ] Setup public routes (GET)
- [ ] Setup protected routes (POST, PATCH)
- [ ] Add middleware to protected routes
- [ ] Add health check endpoint

**Success Criteria:**
- [ ] All routes registered
- [ ] Middleware applied correctly
- [ ] Health check returns 200

**Files:**
- internal/routes/routes.go
```

**Issue #12: Phase 12 - Main Assembly**
```markdown
**Branch:** feature/bookings-api-phase-12-main
**Tasks:**
- [ ] Wire all dependencies in main.go
- [ ] Initialize repositories
- [ ] Initialize services
- [ ] Initialize controllers
- [ ] Setup RabbitMQ consumer in goroutine
- [ ] Start HTTP server

**Success Criteria:**
- [ ] Server starts successfully
- [ ] All dependencies injected
- [ ] Consumer runs in background
- [ ] No dependency injection errors

**Files:**
- cmd/api/main.go
```

### Week 3: Production Ready

**Issue #13: Phase 13 - Docker & Config**
```markdown
**Branch:** feature/bookings-api-phase-13-docker
**Tasks:**
- [ ] Create Dockerfile (multi-stage)
- [ ] Create scripts/init_db.sql
- [ ] Add bookings-api to docker-compose.yml
- [ ] Add environment variables
- [ ] Test docker build
- [ ] Test docker-compose up

**Success Criteria:**
- [ ] Docker image builds
- [ ] Container starts on port 8003
- [ ] Can connect to MySQL
- [ ] Can connect to RabbitMQ

**Files:**
- Dockerfile
- scripts/init_db.sql
- docker-compose.yml
```

**Issue #14: Phase 14 - Testing**
```markdown
**Branch:** feature/bookings-api-phase-14-tests
**Tasks:**
- [ ] Create booking_service_test.go
- [ ] Create idempotency_service_test.go
- [ ] Test duplicate booking prevention
- [ ] Test idempotency (duplicate events)
- [ ] Test concurrent idempotency (race condition)
- [ ] Test seat validation
- [ ] Mock all dependencies

**Success Criteria:**
- [ ] All tests pass
- [ ] Idempotency race condition test passes
- [ ] Coverage > 70%

**Files:**
- internal/service/booking_service_test.go
- internal/service/idempotency_service_test.go
```

**Issue #15: Integration Testing**
```markdown
**Branch:** feature/bookings-api-integration-tests
**Tasks:**
- [ ] Test full flow: POST booking → RabbitMQ → trips-api
- [ ] Test idempotency end-to-end
- [ ] Test trip cancellation flow
- [ ] Test error scenarios

**Success Criteria:**
- [ ] Happy path works end-to-end
- [ ] Duplicate events don't cause issues
- [ ] Compensation events work
```

---

## PR Review Checklist

Before approving PR, verify:
- [ ] Code compiles without warnings
- [ ] All tests pass
- [ ] No hardcoded values (use config)
- [ ] Error handling is proper
- [ ] Logging is structured
- [ ] GORM tags are correct
- [ ] No sensitive data in commits
- [ ] README updated if needed
- [ ] Success criteria met

---

## Timeline Recommendation

**Week 1 (Issues #1-6):** Foundation
- Day 1-2: Setup, database, repository
- Day 3-4: Services, controllers
- Day 5: Middleware, routes

**Week 2 (Issues #7-12):** Integration
- Day 1-2: RabbitMQ consumer (CRITICAL - take your time)
- Day 3: RabbitMQ publisher
- Day 4: HTTP clients
- Day 5: Final assembly

**Week 3 (Issues #13-15):** Production
- Day 1-2: Docker, testing
- Day 3-4: Integration testing
- Day 5: Bug fixes, documentation

---

## Next Steps

1. **Create GitHub Issues** from this breakdown
2. **Start with Issue #1** on branch `feature/bookings-api-phase-1-setup`
3. **Use plan mode** for each phase (see CONTEXT.md)
4. **Create PR after each phase** for review
5. **Merge to dev** after approval
6. **Repeat** until all phases complete
