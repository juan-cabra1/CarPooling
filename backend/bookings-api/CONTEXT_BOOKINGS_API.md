# Bookings API - Implementation Context

## Project Overview
**Academic Project** - CarPooling platform with microservices architecture.

### Requirements
- ✅ MySQL for users-api and bookings-api (this service)
- ✅ MongoDB for trips-api and search-api
- ✅ Apache Solr + Memcached for search-api
- ✅ RabbitMQ for event-driven communication

### Architecture
```
users-api:8001 (MySQL) → ✅ DONE
trips-api:8002 (MongoDB) → ✅ DONE
bookings-api:8003 (MySQL) → TO IMPLEMENT (THIS SERVICE)
search-api:8004 (MongoDB+Solr) → After this
```

---

## Bookings API - Specifications

### Technology Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP)
- **Database:** MySQL 8.0 with GORM
- **Messaging:** RabbitMQ (streadway/amqp)
- **Auth:** JWT (same secret as users-api and trips-api)
- **Logging:** zerolog (structured JSON logs)
- **UUID:** google/uuid v4

### Database: MySQL `carpooling_bookings`

**Table: bookings**
```sql
CREATE TABLE bookings (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  booking_uuid VARCHAR(36) NOT NULL UNIQUE,  -- UUID for external reference

  trip_id VARCHAR(36) NOT NULL,              -- MongoDB ObjectID from trips-api
  passenger_id BIGINT NOT NULL,              -- User ID from users-api

  seats_requested INT NOT NULL,
  total_price DECIMAL(10,2) NOT NULL,        -- seats_requested * price_per_seat

  status VARCHAR(20) NOT NULL,               -- 'pending', 'confirmed', 'cancelled', 'failed'

  cancelled_at DATETIME NULL,
  cancellation_reason TEXT NULL,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_trip_id (trip_id),
  INDEX idx_passenger_id (passenger_id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB;
```

**Table: processed_events (CRITICAL - Idempotency)**
```sql
CREATE TABLE processed_events (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  event_id VARCHAR(36) NOT NULL UNIQUE,      -- UUID from RabbitMQ - UNIQUE CONSTRAINT
  event_type VARCHAR(50) NOT NULL,            -- 'trip.updated', 'trip.cancelled', 'reservation.failed'
  processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  result VARCHAR(20) NOT NULL,                -- 'success', 'skipped', 'failed'
  error_message TEXT NULL,

  INDEX idx_event_type (event_type),
  INDEX idx_processed_at (processed_at)
) ENGINE=InnoDB;
```

---

## REST Endpoints (Port 8003)

```
POST   /bookings                 [Protected] Create booking
GET    /bookings/:id             [Protected] Get booking by ID (owner or trip driver only)
GET    /bookings                 [Protected] List user bookings (?trip_id=X&status=Y&page=1&limit=10)
PATCH  /bookings/:id/cancel      [Protected] Cancel booking (owner only, body: {reason: string})
GET    /health                   [Public]    Health check
```

### Request/Response Examples

**POST /bookings** (Create Booking)
```json
Request:
{
  "trip_id": "507f1f77bcf86cd799439011",
  "seats_requested": 2
}

Response (201):
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "trip_id": "507f1f77bcf86cd799439011",
    "passenger_id": 456,
    "seats_requested": 2,
    "total_price": 10000.00,
    "status": "pending",
    "created_at": "2024-11-10T10:00:00Z"
  }
}
```

**GET /bookings?passenger_id=456&status=confirmed&page=1&limit=10**
```json
Response (200):
{
  "success": true,
  "data": {
    "bookings": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "trip_id": "507f1f77bcf86cd799439011",
        "passenger_id": 456,
        "seats_requested": 2,
        "status": "confirmed",
        "created_at": "2024-11-10T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "limit": 10
  }
}
```

**PATCH /bookings/:id/cancel** (Cancel Booking)
```json
Request:
{
  "reason": "Change of plans"
}

Response (200):
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "cancelled",
    "cancelled_at": "2024-11-10T11:00:00Z",
    "cancellation_reason": "Change of plans"
  }
}
```

---

## RabbitMQ Events

### Published Events (to trips-api):

**reservation.created** (After creating booking)
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.created",
  "trip_id": "507f1f77bcf86cd799439011",
  "seats_reserved": 2,
  "reservation_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-11-10T10:00:00Z"
}
```

**reservation.cancelled** (After cancelling booking)
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.cancelled",
  "trip_id": "507f1f77bcf86cd799439011",
  "seats_released": 2,
  "reservation_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-11-10T11:00:00Z"
}
```

### Consumed Events (from trips-api):

**trip.updated** (Trip seats availability changed)
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.updated",
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 1,
  "reserved_seats": 2,
  "status": "published",
  "timestamp": "2024-11-10T10:00:01Z"
}
```

**trip.cancelled** (Trip was cancelled by driver)
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.cancelled",
  "trip_id": "507f1f77bcf86cd799439011",
  "cancellation_reason": "Car breakdown",
  "timestamp": "2024-11-10T12:00:00Z"
}
```

**reservation.failed** (Compensating event - no seats available)
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.failed",
  "reservation_id": "550e8400-e29b-41d4-a716-446655440000",
  "trip_id": "507f1f77bcf86cd799439011",
  "reason": "No seats available",
  "timestamp": "2024-11-10T10:00:05Z"
}
```

---

## Event Flow

### HAPPY PATH (Successful Booking):
```
1. User calls POST /bookings with trip_id and seats_requested
2. bookings-api validates trip exists (HTTP GET to trips-api)
3. bookings-api creates booking with status='pending'
4. bookings-api publishes reservation.created to RabbitMQ
5. trips-api consumes event, decreases available_seats
6. trips-api publishes trip.updated
7. bookings-api consumes trip.updated, updates booking status='confirmed'
```

### FAILURE SCENARIO (No seats available):
```
1. User calls POST /bookings
2. bookings-api creates booking with status='pending' (based on stale data)
3. bookings-api publishes reservation.created
4. trips-api tries to decrease seats → FAILS (no seats or optimistic lock)
5. trips-api publishes reservation.failed
6. bookings-api consumes reservation.failed → updates booking status='failed'
```

### TRIP CANCELLED SCENARIO:
```
1. Driver cancels trip in trips-api
2. trips-api publishes trip.cancelled
3. bookings-api consumes event (with idempotency check)
4. bookings-api finds ALL bookings for that trip_id with status='confirmed'
5. bookings-api updates ALL those bookings to status='cancelled'
6. bookings-api sends notifications (optional: future feature)
```

### IDEMPOTENCY SCENARIO:
```
1. bookings-api consumes trip.updated (event_id: abc-123)
2. Processing succeeds, booking confirmed
3. ACK fails, RabbitMQ retries same message
4. bookings-api receives SAME event (event_id: abc-123)
5. Idempotency check: event already in processed_events table
6. SKIP processing, just ACK
7. No duplicate status update ✅
```

---

## Critical Requirements

### 1. IDEMPOTENCY (MOST IMPORTANT)

**Problem:** RabbitMQ retries → duplicate events → incorrect booking status or double cancellations

**Solution:**
```go
func HandleTripEvent(event TripEvent) error {
    // 1. Check if already processed (MySQL UNIQUE constraint)
    shouldProcess, err := idempotencyService.CheckAndMarkEvent(event.EventID, event.EventType)
    if err != nil {
        return err // Will retry
    }
    if !shouldProcess {
        logger.Info().Str("event_id", event.EventID).Msg("Event already processed, skipping")
        return nil // ACK without processing
    }

    // 2. Process event
    switch event.EventType {
    case "trip.updated":
        return handleTripUpdated(event)
    case "trip.cancelled":
        return handleTripCancelled(event)
    case "reservation.failed":
        return handleReservationFailed(event)
    }
}
```

**Key:** MySQL `processed_events` table with UNIQUE constraint on `event_id`

### 2. TRIP VALIDATION (Before Creating Booking)

**Problem:** User tries to book a trip that doesn't exist or has no seats

**Solution:**
```go
func CreateBooking(userID int64, req CreateBookingRequest) (*Booking, error) {
    // 1. HTTP GET to trips-api to validate trip exists and has seats
    trip, err := httpClient.GetTrip(req.TripID)
    if err != nil {
        return nil, ErrTripNotFound
    }

    if trip.AvailableSeats < req.SeatsRequested {
        return nil, ErrInsufficientSeats
    }

    if trip.Status != "published" {
        return nil, ErrTripNotAvailable
    }

    // 2. Create booking with status='pending'
    booking := &Booking{
        BookingUUID:    uuid.New().String(),
        TripID:         req.TripID,
        PassengerID:    userID,
        SeatsRequested: req.SeatsRequested,
        TotalPrice:     trip.PricePerSeat * float64(req.SeatsRequested),
        Status:         "pending",
    }

    // 3. Save to database
    if err := repo.Create(booking); err != nil {
        return nil, err
    }

    // 4. Publish reservation.created event
    event := ReservationCreatedEvent{
        EventID:       uuid.New().String(),
        EventType:     "reservation.created",
        TripID:        booking.TripID,
        SeatsReserved: booking.SeatsRequested,
        ReservationID: booking.BookingUUID,
        Timestamp:     time.Now(),
    }
    publisher.PublishReservationCreated(event)

    return booking, nil
}
```

---

## Architecture Pattern

```
internal/
├── config/              # Environment variables
│   └── config.go
├── domain/              # Business models (DTOs, requests, responses)
│   ├── booking.go
│   ├── trip.go          # Trip DTO for external API calls
│   └── errors.go
├── dao/                 # GORM models (database entities)
│   ├── booking.go
│   └── processed_event.go
├── repository/          # Data access interfaces + GORM implementations
│   ├── booking_repository.go
│   └── event_repository.go
├── service/             # Business logic
│   ├── booking_service.go
│   ├── idempotency_service.go
│   └── auth_service.go
├── controller/          # HTTP handlers
│   └── booking_controller.go
├── middleware/          # JWT, CORS, errors
│   ├── auth.go
│   ├── cors.go
│   └── error.go
├── messaging/           # RabbitMQ
│   ├── rabbitmq.go
│   ├── trip_consumer.go
│   └── reservation_publisher.go
├── http/                # External HTTP clients
│   └── trips_client.go
└── routes/
    └── routes.go

cmd/api/
└── main.go             # Dependency injection, server start
```

---

## Dependencies (go.mod)

```go
require (
    github.com/gin-gonic/gin v1.10.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/joho/godotenv v1.5.1
    github.com/google/uuid v1.6.0
    github.com/rs/zerolog v1.33.0
    github.com/streadway/amqp v1.1.0
    gorm.io/driver/mysql v1.5.7
    gorm.io/gorm v1.25.10
)
```

---

## Environment Variables (.env)

```bash
# MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=carpooling_bookings

# JWT (same secret as users-api and trips-api)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
SERVER_PORT=8003

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# External APIs
TRIPS_API_URL=http://localhost:8002
USERS_API_URL=http://localhost:8001
```

---

## Business Validations

### Creating Booking:
- `trip_id` must exist (HTTP GET to trips-api)
- Trip `status` must be 'published'
- Trip must have `available_seats >= seats_requested`
- `seats_requested` must be between 1-8
- User cannot book their own trip (JWT `user_id` != trip `driver_id`)
- User cannot have multiple active bookings for same trip
- `status` = 'pending' initially (waiting for confirmation from trips-api)
- `booking_uuid` = UUID v4
- `total_price` = `trip.price_per_seat * seats_requested`

### Cancelling Booking:
- Only booking owner can cancel (JWT `user_id` == `passenger_id`)
- Cannot cancel if status is already 'cancelled' or 'failed'
- Set `status` = 'cancelled'
- Set `cancelled_at` = current timestamp
- Publish `reservation.cancelled` event to release seats
- trips-api will increase `available_seats`

### Consuming trip.cancelled:
- Find ALL bookings for `trip_id` with `status='confirmed'`
- Update ALL to `status='cancelled'`
- Set `cancellation_reason` = event reason
- DO NOT publish reservation.cancelled (seats already freed by trip cancellation)

---

## Error Codes (domain/errors.go)

```go
type AppError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

var (
    ErrBookingNotFound       = &AppError{Code: "BOOKING_NOT_FOUND", Message: "Booking not found"}
    ErrTripNotFound          = &AppError{Code: "TRIP_NOT_FOUND", Message: "Trip not found"}
    ErrTripNotAvailable      = &AppError{Code: "TRIP_NOT_AVAILABLE", Message: "Trip not available for booking"}
    ErrInsufficientSeats     = &AppError{Code: "INSUFFICIENT_SEATS", Message: "Not enough seats available"}
    ErrUnauthorized          = &AppError{Code: "UNAUTHORIZED", Message: "Not authorized"}
    ErrCannotBookOwnTrip     = &AppError{Code: "CANNOT_BOOK_OWN_TRIP", Message: "Cannot book your own trip"}
    ErrDuplicateBooking      = &AppError{Code: "DUPLICATE_BOOKING", Message: "Already have active booking for this trip"}
    ErrCannotCancelBooking   = &AppError{Code: "CANNOT_CANCEL_BOOKING", Message: "Cannot cancel booking in current status"}
    ErrInvalidSeats          = &AppError{Code: "INVALID_SEATS", Message: "Invalid number of seats"}
)
```

---

## Testing Focus

### Unit Tests (service layer)
```go
// Test: Booking creation validation
func TestCreateBooking_InsufficientSeats(t *testing.T) {
    // Mock trips-api returns trip with 1 available seat
    // Request 2 seats
    // Should return ErrInsufficientSeats
}

// Test: Cannot book own trip
func TestCreateBooking_CannotBookOwnTrip(t *testing.T) {
    // user_id == driver_id
    // Should return ErrCannotBookOwnTrip
}

// Test: Duplicate booking prevention
func TestCreateBooking_DuplicateBooking(t *testing.T) {
    // User already has active booking for trip
    // Should return ErrDuplicateBooking
}
```

### Idempotency Tests (CRITICAL)
```go
// Test: Duplicate event skipped
func TestIdempotency_DuplicateEventSkipped(t *testing.T) {
    eventID := "event-123"

    // First call - should process
    shouldProcess, _ := service.CheckAndMarkEvent(eventID, "trip.updated")
    assert.True(t, shouldProcess)

    // Second call - should skip
    shouldProcess, _ = service.CheckAndMarkEvent(eventID, "trip.updated")
    assert.False(t, shouldProcess)
}

// Test: Concurrent duplicate events
func TestIdempotency_ConcurrentDuplicates(t *testing.T) {
    var wg sync.WaitGroup
    processed := atomic.Int32{}

    // 10 goroutines try to process same event
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            shouldProcess, _ := service.CheckAndMarkEvent("event-race", "trip.updated")
            if shouldProcess {
                processed.Add(1)
            }
        }()
    }

    wg.Wait()
    assert.Equal(t, int32(1), processed.Load()) // Only ONE processed
}
```

### Event Handling Tests
```go
// Test: trip.cancelled updates all bookings
func TestHandleTripCancelled_UpdatesAllBookings(t *testing.T) {
    // Create 3 confirmed bookings for trip
    // Publish trip.cancelled event
    // All 3 bookings should be status='cancelled'
}

// Test: reservation.failed updates booking
func TestHandleReservationFailed_UpdatesBookingStatus(t *testing.T) {
    // Create pending booking
    // Publish reservation.failed event
    // Booking should be status='failed'
}
```

---

## Common Patterns from users-api

### Config Loading (Same as users-api)
```go
type Config struct {
    DBHost       string
    DBPort       string
    DBUser       string
    DBPassword   string
    DBName       string
    JWTSecret    string
    ServerPort   string
    RabbitMQURL  string
    TripsAPIURL  string
    UsersAPIURL  string
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    // Load from os.Getenv with defaults...
}
```

### GORM Model Pattern
```go
// dao/booking.go
type Booking struct {
    ID             uint      `gorm:"primaryKey;autoIncrement"`
    BookingUUID    string    `gorm:"type:varchar(36);uniqueIndex;not null"`
    TripID         string    `gorm:"type:varchar(36);index;not null"`
    PassengerID    int64     `gorm:"index;not null"`
    SeatsRequested int       `gorm:"not null"`
    TotalPrice     float64   `gorm:"type:decimal(10,2);not null"`
    Status         string    `gorm:"type:varchar(20);index;not null"`
    CancelledAt    *time.Time
    CancellationReason string `gorm:"type:text"`
    CreatedAt      time.Time `gorm:"autoCreateTime"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

// domain/booking.go (DTO)
type Booking struct {
    ID                 string    `json:"id"`
    TripID             string    `json:"trip_id"`
    PassengerID        int64     `json:"passenger_id"`
    SeatsRequested     int       `json:"seats_requested"`
    TotalPrice         float64   `json:"total_price"`
    Status             string    `json:"status"`
    CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
    CancellationReason string     `json:"cancellation_reason,omitempty"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
}
```

### Repository Pattern (Same as users-api)
```go
type BookingRepository interface {
    Create(booking *dao.Booking) error
    FindByID(id string) (*dao.Booking, error)
    FindByPassengerID(passengerID int64, filters map[string]interface{}) ([]*dao.Booking, error)
    FindByTripID(tripID string, status string) ([]*dao.Booking, error)
    Update(booking *dao.Booking) error
    // ...
}

type bookingRepository struct {
    db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
    return &bookingRepository{db: db}
}
```

---

## Key Differences from trips-api

| Aspect | trips-api | bookings-api |
|--------|-----------|--------------|
| **Database** | MongoDB with native driver | MySQL with GORM |
| **ORM** | Manual document mapping | GORM auto-migration |
| **Primary Key** | ObjectID (MongoDB) | Auto-increment BIGINT + UUID |
| **Concurrency** | Optimistic locking needed | No optimistic locking needed |
| **External calls** | users-api (validate driver) | trips-api (validate trip) |
| **Messaging Role** | Consumer + Publisher | Consumer + Publisher |
| **Idempotency** | MongoDB unique index | MySQL UNIQUE constraint |
| **Main Events Consumed** | reservation.created/cancelled | trip.updated/cancelled/reservation.failed |
| **Main Events Published** | trip.created/updated/cancelled/reservation.failed | reservation.created/cancelled |

---

## Success Criteria (Complete Implementation)

✅ Server runs on port 8003
✅ All 5 REST endpoints working
✅ JWT authentication on protected routes
✅ MySQL connection established with GORM
✅ Tables created with proper indexes and constraints
✅ RabbitMQ consumer running in background (consumes from trips-api)
✅ RabbitMQ publisher emitting events with UUIDs
✅ Idempotency working (duplicate events skipped via MySQL UNIQUE constraint)
✅ HTTP client calling trips-api for validation
✅ Structured logging with zerolog
✅ Health check returns service status
✅ Unit tests passing (especially idempotency tests)
✅ Event handlers working (trip.cancelled updates all bookings)
✅ Docker container builds and runs

---

## References

- **GORM Documentation:** https://gorm.io/docs/
- **Gin Framework:** https://gin-gonic.com/docs/
- **RabbitMQ Go:** https://www.rabbitmq.com/tutorials/tutorial-two-go.html
- **users-api code:** `/home/user/CarPooling/backend/users-api` (MySQL/GORM reference)
- **trips-api code:** `/home/user/CarPooling/backend/trips-api` (RabbitMQ patterns)
- **GITFLOW:** See `GITFLOW.md` for workflow

---

## Important Notes

1. **MySQL with GORM** - Use auto-migration, similar to users-api
2. **UUID for external reference** - booking_uuid is the public ID, internal id is BIGINT
3. **Idempotency with UNIQUE constraint** - MySQL will reject duplicate event_id
4. **Generate UUIDs before publishing** - Never reuse event IDs
5. **Manual ACK in RabbitMQ** - Only after successful processing
6. **Structured logging** - Log all key events with event_id and booking_uuid
7. **Error codes, not messages** - Use AppError with codes
8. **HTTP client for trips-api** - Validate trips before creating bookings
9. **Eventual consistency** - Booking starts as 'pending', becomes 'confirmed' after trip.updated event
10. **Compensating events** - Handle reservation.failed to update booking status

---

## Quick Start for Implementation

When starting a feature:

1. **Checkout branch:** `git checkout -b feature/bookings-api/{issue-number}-{name}`
2. **Read this context:** Full specification here
3. **Copy patterns:** Reference users-api for GORM patterns, trips-api for RabbitMQ patterns
4. **Use plan mode:** `@CONTEXT_BOOKINGS_API.md Implement feature X`
5. **Test:** Verify success criteria
6. **Commit:** `git commit -m "feat(bookings): {description}"`
7. **Push & PR:** Push to branch, create PR to dev
