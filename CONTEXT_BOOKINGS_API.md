# Bookings API - Implementation Context

## Project Overview
**Academic Project** - CarPooling platform with microservices architecture.

### Requirements (Faculty)
- ✅ MongoDB for trips-api (main API)
- ✅ MySQL for users-api and bookings-api
- ✅ Apache Solr + Memcached for search-api
- ✅ RabbitMQ for event-driven communication

### Architecture
```
users-api:8001 (MySQL) → DONE
bookings-api:8003 (MySQL) → TO IMPLEMENT
trips-api:8002 (MongoDB) → TODO
search-api:8004 (MongoDB+Solr+Memcached) → TODO
```

---

## Bookings API - Specifications

### Technology Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP)
- **Database:** MySQL 5.7+ with GORM
- **Messaging:** RabbitMQ (streadway/amqp)
- **Auth:** JWT (same secret as users-api)
- **Logging:** zerolog (structured JSON logs)
- **UUID:** google/uuid v4

### Database: `carpooling_bookings`

**Table: bookings**
```sql
id BIGINT PK AUTO_INCREMENT
reservation_id VARCHAR(36) UNIQUE NOT NULL  -- UUID for idempotency
user_id BIGINT NOT NULL                     -- FK to users (logical)
trip_id VARCHAR(24) NOT NULL                -- MongoDB ObjectID from trips-api
seats_reserved INT NOT NULL
total_price DECIMAL(10,2) NOT NULL
status ENUM('pending','confirmed','cancelled','failed') DEFAULT 'pending'
payment_status ENUM('pending','completed','refunded') DEFAULT 'pending'
passenger_name VARCHAR(200) NOT NULL
passenger_phone VARCHAR(20) NOT NULL
passenger_email VARCHAR(255) NOT NULL
cancellation_reason TEXT NULL
cancelled_at DATETIME NULL
trip_driver_id BIGINT NULL                  -- snapshot
trip_origin_city VARCHAR(100) NULL
trip_destination_city VARCHAR(100) NULL
trip_departure_datetime DATETIME NULL
trip_price_per_seat DECIMAL(10,2) NULL
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

INDEX: user_id, trip_id, status, reservation_id, created_at
```

**Table: processed_events (CRITICAL - Idempotency)**
```sql
id BIGINT PK AUTO_INCREMENT
event_id VARCHAR(36) UNIQUE NOT NULL        -- From RabbitMQ messages
event_type VARCHAR(50) NOT NULL
processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
result ENUM('success','skipped','failed') DEFAULT 'success'
error_message TEXT NULL

UNIQUE KEY: event_id (prevents duplicate processing)
INDEX: event_type, processed_at
```

### REST Endpoints (Port 8003)
```
POST   /bookings                 [Protected] Create booking
GET    /bookings/:id             [Public]    Get booking by ID
GET    /bookings                 [Public]    List bookings (?user_id=X&status=Y)
PATCH  /bookings/:id/confirm     [Protected] Confirm booking (owner only)
PATCH  /bookings/:id/cancel      [Protected] Cancel booking (owner only)
GET    /health                   [Public]    Health check
```

### RabbitMQ Events

**Consumes from `trips.events` queue:**
```json
{
  "event_id": "uuid-v4",              // CRITICAL for idempotency
  "event_type": "trip.updated",       // or "trip.cancelled"
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 2,
  "status": "published",
  "timestamp": "2024-11-10T10:00:00Z"
}
```

**Publishes to `reservations.events` exchange:**
```json
{
  "event_id": "uuid-v4",              // Generate before publishing
  "event_type": "reservation.created", // or "reservation.cancelled"/"reservation.failed"
  "reservation_id": "uuid",
  "trip_id": "507f1f77bcf86cd799439011",
  "user_id": 456,
  "seats_reserved": 2,
  "timestamp": "2024-11-10T10:00:01Z"
}
```

---

## Critical Requirements

### 1. IDEMPOTENCY (MOST IMPORTANT)
**Problem:** RabbitMQ retries failed messages → duplicate processing → double bookings

**Solution:**
```go
func HandleTripEvent(event TripEvent) error {
    // 1. Check if already processed
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
    }
}
```

**Key Points:**
- UNIQUE constraint on `processed_events.event_id`
- Check BEFORE processing
- Generate UUID for published events
- Manual ACK after successful processing
- NACK with requeue on transient errors

### 2. Event Flow (Happy Path)
```
1. POST /bookings → validate trip exists (HTTP GET trips-api)
2. Create booking (status: 'pending')
3. Publish reservation.created (with UUID event_id)
4. trips-api consumes → decreases available_seats
5. trips-api publishes trip.updated (with UUID event_id)
6. bookings-api consumes → checks idempotency → updates booking (status: 'confirmed')
7. ACK message
```

### 3. Compensation Flow (No Seats)
```
1. POST /bookings → creates booking (based on stale data)
2. Publish reservation.created
3. trips-api tries to decrease seats → FAILS (optimistic lock or no seats)
4. trips-api publishes reservation.failed
5. bookings-api consumes → updates booking (status: 'failed')
6. User sees error
```

### 4. Business Validations
- User cannot book same trip twice (check active bookings)
- seats_reserved must be >= 1 and <= trip.available_seats
- User must exist (HTTP GET users-api)
- Trip must exist and have seats (HTTP GET trips-api)
- Calculate total_price = trip.price_per_seat * seats_reserved
- Capture trip snapshot for historical reference

---

## Architecture Pattern (Same as users-api)

```
internal/
├── config/              # Environment variables
│   └── config.go
├── dao/                 # GORM models
│   ├── booking.go
│   └── processed_event.go
├── domain/              # DTOs, requests, responses
│   ├── booking.go
│   └── errors.go
├── repository/          # Data access interfaces + implementations
│   ├── booking_repository.go
│   └── event_repository.go
├── service/             # Business logic
│   ├── booking_service.go
│   └── idempotency_service.go
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
├── clients/             # HTTP clients for other services
│   ├── trips_client.go
│   └── users_client.go
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
    gorm.io/gorm v1.25.12
    gorm.io/driver/mysql v1.5.7
)
```

---

## Environment Variables (.env)

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Prueba.9876
DB_NAME=carpooling_bookings

# JWT (same secret as users-api)
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

## Error Codes (domain/errors.go)

```go
type AppError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

var (
    ErrTripNotFound      = &AppError{Code: "TRIP_NOT_FOUND", Message: "Trip not found"}
    ErrUserNotFound      = &AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrNoSeatsAvailable  = &AppError{Code: "NO_SEATS_AVAILABLE", Message: "Not enough seats"}
    ErrDuplicateBooking  = &AppError{Code: "DUPLICATE_BOOKING", Message: "Already booked this trip"}
    ErrBookingNotFound   = &AppError{Code: "BOOKING_NOT_FOUND", Message: "Booking not found"}
    ErrUnauthorized      = &AppError{Code: "UNAUTHORIZED", Message: "Not authorized"}
)
```

---

## Testing Focus

### Unit Tests (service layer)
```go
// Test: Duplicate booking prevention
func TestCreateBooking_DuplicateBooking(t *testing.T) {
    // User already has active booking for trip
    // Should return ErrDuplicateBooking
}

// Test: Seat validation
func TestCreateBooking_NoSeatsAvailable(t *testing.T) {
    // Trip has 2 seats, user requests 3
    // Should return ErrNoSeatsAvailable
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

// Test: Concurrent duplicate events (race condition)
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

---

## Common Patterns to Copy from users-api

### 1. Config Loading
```go
// Copy internal/config/config.go structure
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

### 2. GORM Models
```go
// Copy GORM tag patterns
type BookingDAO struct {
    ID   int64  `gorm:"primaryKey;autoIncrement;column:id"`
    Name string `gorm:"type:varchar(100);not null;column:name"`
    // ...
}

func (BookingDAO) TableName() string {
    return "bookings"
}
```

### 3. Repository Pattern
```go
// Copy interface-based repository
type BookingRepository interface {
    Create(booking *dao.BookingDAO) error
    FindByID(id int64) (*dao.BookingDAO, error)
    // ...
}

type bookingRepository struct {
    db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
    return &bookingRepository{db: db}
}
```

### 4. Controller Pattern
```go
// Copy controller structure
type BookingController interface {
    CreateBooking(c *gin.Context)
    // ...
}

type bookingController struct {
    bookingService service.BookingService
}

func (ctrl *bookingController) CreateBooking(c *gin.Context) {
    var req domain.CreateBookingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"success": false, "error": err.Error()})
        return
    }
    // ...
    c.JSON(201, gin.H{"success": true, "data": booking})
}
```

### 5. Middleware
```go
// Copy auth middleware exactly (same JWT secret)
// Copy CORS middleware
// Copy error handler middleware
```

### 6. Main.go Initialization
```go
func main() {
    // 1. Load config
    // 2. Connect to MySQL
    // 3. Auto-migrate
    // 4. Connect to RabbitMQ
    // 5. Initialize repositories
    // 6. Initialize services
    // 7. Initialize controllers
    // 8. Setup routes
    // 9. Start consumer (goroutine)
    // 10. Start HTTP server
}
```

---

## Key Differences from users-api

| Aspect | users-api | bookings-api |
|--------|-----------|--------------|
| **External calls** | None | trips-api, users-api |
| **Messaging** | None | RabbitMQ consumer + publisher |
| **Idempotency** | N/A | Critical (processed_events table) |
| **UUIDs** | No | Yes (reservation_id, event_id) |
| **Snapshots** | No | Yes (trip data for history) |
| **Async processing** | No | Yes (RabbitMQ consumer in goroutine) |

---

## Success Criteria (All Phases Complete)

✅ Server runs on port 8003
✅ All 5 REST endpoints working
✅ JWT authentication on protected routes
✅ MySQL tables created with proper indexes
✅ RabbitMQ consumer running in background
✅ RabbitMQ publisher emitting events with UUIDs
✅ Idempotency working (duplicate events skipped)
✅ HTTP clients calling trips-api and users-api
✅ Structured logging with zerolog
✅ Health check returns service status
✅ Unit tests passing (especially idempotency tests)
✅ Docker container builds and runs
✅ Can create booking end-to-end
✅ Trip cancellation cancels all bookings

---

## References

- **Implementation Plan:** See GITFLOW.md for detailed phases
- **users-api code:** `/home/user/CarPooling/backend/users-api`
- **Architecture analysis:** Project review document
- **GitHub Issues:** Follow issue templates for each phase

---

## Important Notes

1. **Copy patterns from users-api** - Don't reinvent the wheel
2. **Idempotency is CRITICAL** - Take extra time on Phase 7
3. **Generate UUIDs before publishing** - Never reuse event IDs
4. **Manual ACK in RabbitMQ** - Only after successful processing
5. **Structured logging** - Log all key events (booking created, event processed, etc.)
6. **Error codes, not messages** - Use AppError with codes
7. **Test idempotency thoroughly** - Include race condition tests

---

## Quick Start for Each Phase

When starting a new phase:

1. **Checkout branch:** `git checkout -b feature/bookings-api-phase-N-{name}`
2. **Read phase details:** See GITFLOW.md Issue #N
3. **Copy patterns:** Reference users-api for similar code
4. **Implement:** Follow task checklist
5. **Test:** Verify success criteria
6. **Commit:** `git commit -m "feat(bookings): {description}"`
7. **Push & PR:** Push to branch, create PR to dev
8. **Review:** Wait for approval
9. **Merge:** Merge to dev
10. **Next phase:** Start from step 1

---

## Contact for Questions

- Review project analysis document for architecture decisions
- Check users-api implementation for patterns
- Reference this CONTEXT.md for specifications
- See GITFLOW.md for workflow and issues
