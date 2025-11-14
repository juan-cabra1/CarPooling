# Trips API - Implementation Context

## Project Overview
**Academic Project** - CarPooling platform with microservices architecture.

### Requirements (Faculty)
- ✅ **MongoDB for trips-api** (main API - this service)
- ✅ MySQL for users-api and bookings-api
- ✅ Apache Solr + Memcached for search-api
- ✅ RabbitMQ for event-driven communication

### Architecture
```
users-api:8001 (MySQL) → ✅ DONE
trips-api:8002 (MongoDB) → TO IMPLEMENT (THIS SERVICE)
bookings-api:8003 (MySQL) → After this
search-api:8004 (MongoDB+Solr) → After this
```

---

## Trips API - Specifications

### Technology Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP)
- **Database:** MongoDB with native driver
- **Messaging:** RabbitMQ (streadway/amqp)
- **Auth:** JWT (same secret as users-api)
- **Logging:** zerolog (structured JSON logs)
- **UUID:** google/uuid v4

### Database: MongoDB `carpooling_trips`

**Collection: trips**
```javascript
{
  _id: ObjectId(),
  driver_id: Number,

  origin: {
    city: String,
    province: String,
    address: String,
    coordinates: { lat: Number, lng: Number }
  },
  destination: {
    city: String,
    province: String,
    address: String,
    coordinates: { lat: Number, lng: Number }
  },

  departure_datetime: Date,
  estimated_arrival_datetime: Date,

  price_per_seat: Number,
  total_seats: Number,
  reserved_seats: Number,
  available_seats: Number,
  availability_version: Number,  // For optimistic locking

  car: {
    brand: String,
    model: String,
    year: Number,
    color: String,
    plate: String
  },

  preferences: {
    pets_allowed: Boolean,
    smoking_allowed: Boolean,
    music_allowed: Boolean
  },

  status: String,  // 'draft', 'published', 'full', 'in_progress', 'completed', 'cancelled'
  description: String,

  cancelled_at: Date,
  cancelled_by: Number,
  cancellation_reason: String,

  created_at: Date,
  updated_at: Date
}
```

**Indexes:**
```javascript
db.trips.createIndex({ driver_id: 1 })
db.trips.createIndex({ status: 1 })
db.trips.createIndex({ departure_datetime: 1 })
db.trips.createIndex({ "origin.city": 1, "destination.city": 1 })
```

**Collection: processed_events (CRITICAL - Idempotency)**
```javascript
{
  _id: ObjectId(),
  event_id: String,        // UUID from RabbitMQ - UNIQUE INDEX
  event_type: String,      // 'reservation.created', 'reservation.cancelled'
  processed_at: Date,
  result: String,          // 'success', 'skipped', 'failed'
  error_message: String
}
```

**Indexes:**
```javascript
db.processed_events.createIndex({ event_id: 1 }, { unique: true })  // CRITICAL
db.processed_events.createIndex({ event_type: 1 })
db.processed_events.createIndex({ processed_at: 1 })
```

---

## REST Endpoints (Port 8002)

```
POST   /trips                    [Protected] Create trip
GET    /trips/:id                [Public]    Get trip by ID
GET    /trips                    [Public]    List trips (?driver_id=X&status=Y&page=1&limit=10)
PUT    /trips/:id                [Protected] Update trip (owner only)
DELETE /trips/:id                [Protected] Delete trip (owner only)
PATCH  /trips/:id/cancel         [Protected] Cancel trip (owner only, body: {reason: string})
GET    /health                   [Public]    Health check
```

### Request/Response Examples

**POST /trips** (Create Trip)
```json
Request:
{
  "origin": {
    "city": "Buenos Aires",
    "province": "Buenos Aires",
    "address": "Av. 9 de Julio 1234",
    "coordinates": { "lat": -34.6037, "lng": -58.3816 }
  },
  "destination": {
    "city": "Rosario",
    "province": "Santa Fe",
    "address": "Pellegrini 1234",
    "coordinates": { "lat": -32.9468, "lng": -60.6393 }
  },
  "departure_datetime": "2024-12-01T10:00:00Z",
  "estimated_arrival_datetime": "2024-12-01T14:00:00Z",
  "price_per_seat": 5000,
  "total_seats": 3,
  "car": {
    "brand": "Toyota",
    "model": "Corolla",
    "year": 2020,
    "color": "Gris",
    "plate": "ABC123"
  },
  "preferences": {
    "pets_allowed": false,
    "smoking_allowed": false,
    "music_allowed": true
  },
  "description": "Viaje tranquilo a Rosario"
}

Response (201):
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "driver_id": 123,
    "origin": { ... },
    "destination": { ... },
    "available_seats": 3,
    "reserved_seats": 0,
    "status": "published",
    "created_at": "2024-11-10T10:00:00Z"
  }
}
```

**GET /trips?origin_city=Buenos Aires&destination_city=Rosario&page=1&limit=10**
```json
Response (200):
{
  "success": true,
  "data": {
    "trips": [ {...}, {...} ],
    "total": 25,
    "page": 1,
    "limit": 10
  }
}
```

---

## RabbitMQ Events

### Consumed Events (from bookings-api):

**reservation.created**
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.created",
  "trip_id": "507f1f77bcf86cd799439011",
  "seats_reserved": 2,
  "reservation_id": "uuid",
  "timestamp": "2024-11-10T10:00:00Z"
}
```

**reservation.cancelled**
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.cancelled",
  "trip_id": "507f1f77bcf86cd799439011",
  "seats_released": 2,
  "reservation_id": "uuid",
  "timestamp": "2024-11-10T11:00:00Z"
}
```

### Published Events (to bookings-api and search-api):

**trip.created, trip.updated, trip.cancelled**
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.updated",
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 2,
  "reserved_seats": 1,
  "status": "published",
  "timestamp": "2024-11-10T10:00:01Z"
}
```

**reservation.failed** (Compensating event)
```json
{
  "event_id": "uuid-v4",
  "event_type": "reservation.failed",
  "reservation_id": "uuid",
  "trip_id": "507f1f77bcf86cd799439011",
  "reason": "No seats available",
  "timestamp": "2024-11-10T10:00:05Z"
}
```

---

## Event Flow with Failure Handling

### HAPPY PATH:
```
1. bookings-api creates reservation
2. bookings-api publishes reservation.created
3. trips-api consumes event (with idempotency check)
4. trips-api decreases available_seats (optimistic locking)
5. trips-api publishes trip.updated
6. bookings-api/search-api consume trip.updated
```

### FAILURE SCENARIO (No seats available):
```
1. bookings-api creates reservation (based on stale data)
2. bookings-api publishes reservation.created
3. trips-api consumes event
4. trips-api tries UpdateAvailability with optimistic lock
5. Update FAILS (no seats or version conflict)
6. trips-api publishes reservation.failed
7. bookings-api consumes → marks reservation as failed
```

### IDEMPOTENCY SCENARIO:
```
1. trips-api consumes reservation.created (event_id: abc-123)
2. Processing succeeds, seats decreased
3. ACK fails, RabbitMQ retries same message
4. trips-api receives SAME event (event_id: abc-123)
5. Idempotency check: event already in processed_events
6. SKIP processing, just ACK
7. No duplicate seat decrease ✅
```

---

## Critical Requirements

### 1. IDEMPOTENCY (MOST IMPORTANT)

**Problem:** RabbitMQ retries → duplicate events → incorrect seat count

**Solution:**
```go
func HandleReservationEvent(event ReservationEvent) error {
    // 1. Check if already processed (MongoDB unique index)
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
    case "reservation.created":
        return handleReservationCreated(event)
    case "reservation.cancelled":
        return handleReservationCancelled(event)
    }
}
```

**Key:** MongoDB `processed_events` collection with UNIQUE index on `event_id`

### 2. OPTIMISTIC LOCKING (Prevent Race Conditions)

**Problem:** Two reservations try to book last seat simultaneously

**Solution:**
```go
func UpdateAvailability(tripID string, seatsDelta int, expectedVersion int) error {
    filter := bson.M{
        "_id": tripID,
        "availability_version": expectedVersion,
        "available_seats": bson.M{"$gte": -seatsDelta}, // Has enough seats
    }

    update := bson.M{
        "$inc": bson.M{
            "available_seats": seatsDelta,
            "reserved_seats": -seatsDelta,
            "availability_version": 1,
        },
    }

    result, err := collection.UpdateOne(ctx, filter, update)
    if result.MatchedCount == 0 {
        return ErrOptimisticLockFailed // No seats or version mismatch
    }
    return nil
}
```

**Key:** Increment `availability_version` on every update, check version in filter

---

## Architecture Pattern

```
internal/
├── config/              # Environment variables
│   └── config.go
├── domain/              # Business models (DTOs, requests, responses)
│   ├── trip.go
│   └── errors.go
├── repository/          # Data access interfaces + MongoDB implementations
│   ├── trip_repository.go
│   └── event_repository.go
├── service/             # Business logic
│   ├── trip_service.go
│   └── idempotency_service.go
├── controller/          # HTTP handlers
│   └── trip_controller.go
├── middleware/          # JWT, CORS, errors
│   ├── auth.go
│   ├── cors.go
│   └── error.go
├── messaging/           # RabbitMQ
│   ├── rabbitmq.go
│   ├── reservation_consumer.go
│   └── trip_publisher.go
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
    go.mongodb.org/mongo-driver v1.13.1
)
```

---

## Environment Variables (.env)

```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=carpooling_trips

# JWT (same secret as users-api)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
SERVER_PORT=8002

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# External APIs
USERS_API_URL=http://localhost:8001
```

---

## Business Validations

### Creating Trip:
- `departure_datetime` must be in the future
- `total_seats` must be between 1-8
- `driver_id` must exist (HTTP GET to users-api)
- `available_seats` = `total_seats` initially
- `reserved_seats` = 0 initially
- `status` = 'published' initially
- `availability_version` = 1 initially

### Updating Trip:
- Only trip owner can update (JWT `user_id` == `driver_id`)
- Cannot update if `reserved_seats` > 0
- Cannot change `total_seats` to less than `reserved_seats`

### Cancelling Trip:
- Only trip owner can cancel
- Set `status` = 'cancelled'
- Publish `trip.cancelled` event
- bookings-api will cancel all reservations

---

## Error Codes (domain/errors.go)

```go
type AppError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

var (
    ErrTripNotFound         = &AppError{Code: "TRIP_NOT_FOUND", Message: "Trip not found"}
    ErrDriverNotFound       = &AppError{Code: "DRIVER_NOT_FOUND", Message: "Driver not found"}
    ErrNoSeatsAvailable     = &AppError{Code: "NO_SEATS_AVAILABLE", Message: "No seats available"}
    ErrOptimisticLockFailed = &AppError{Code: "OPTIMISTIC_LOCK_FAILED", Message: "Version conflict"}
    ErrUnauthorized         = &AppError{Code: "UNAUTHORIZED", Message: "Not authorized"}
    ErrPastDeparture        = &AppError{Code: "PAST_DEPARTURE", Message: "Departure must be in future"}
    ErrHasReservations      = &AppError{Code: "HAS_RESERVATIONS", Message: "Cannot modify trip with reservations"}
)
```

---

## Testing Focus

### Unit Tests (service layer)
```go
// Test: Trip creation validation
func TestCreateTrip_PastDeparture(t *testing.T) {
    // departure_datetime in the past
    // Should return ErrPastDeparture
}

// Test: Seat validation
func TestCreateTrip_InvalidSeats(t *testing.T) {
    // total_seats = 10 (> 8)
    // Should return error
}
```

### Idempotency Tests (CRITICAL)
```go
// Test: Duplicate event skipped
func TestIdempotency_DuplicateEventSkipped(t *testing.T) {
    eventID := "event-123"

    // First call - should process
    shouldProcess, _ := service.CheckAndMarkEvent(eventID, "reservation.created")
    assert.True(t, shouldProcess)

    // Second call - should skip
    shouldProcess, _ = service.CheckAndMarkEvent(eventID, "reservation.created")
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
            shouldProcess, _ := service.CheckAndMarkEvent("event-race", "reservation.created")
            if shouldProcess {
                processed.Add(1)
            }
        }()
    }

    wg.Wait()
    assert.Equal(t, int32(1), processed.Load()) // Only ONE processed
}
```

### Optimistic Locking Tests
```go
// Test: Concurrent seat updates
func TestUpdateAvailability_OptimisticLock(t *testing.T) {
    // Two goroutines try to reserve last seat
    // One should succeed, other should fail with ErrOptimisticLockFailed
}
```

---

## Common Patterns from users-api

### Config Loading
```go
type Config struct {
    MongoURI     string
    MongoDB      string
    JWTSecret    string
    ServerPort   string
    RabbitMQURL  string
    UsersAPIURL  string
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    // Load from os.Getenv with defaults...
}
```

### Repository Pattern
```go
type TripRepository interface {
    Create(trip *domain.Trip) error
    FindByID(id string) (*domain.Trip, error)
    // ...
}

type tripRepository struct {
    collection *mongo.Collection
}

func NewTripRepository(db *mongo.Database) TripRepository {
    return &tripRepository{
        collection: db.Collection("trips"),
    }
}
```

### Controller Pattern
```go
func (ctrl *tripController) CreateTrip(c *gin.Context) {
    var req domain.CreateTripRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"success": false, "error": err.Error()})
        return
    }

    // Get user_id from JWT (set by AuthMiddleware)
    userID, _ := c.Get("user_id")

    trip, err := ctrl.tripService.CreateTrip(userID.(int64), req)
    if err != nil {
        // Handle specific error codes
        switch err.(*domain.AppError).Code {
        case "TRIP_NOT_FOUND":
            c.JSON(404, gin.H{"success": false, "error": err.Error()})
        default:
            c.JSON(500, gin.H{"success": false, "error": err.Error()})
        }
        return
    }

    c.JSON(201, gin.H{"success": true, "data": trip})
}
```

---

## Key Differences from users-api

| Aspect | users-api | trips-api |
|--------|-----------|-----------|
| **Database** | MySQL with GORM | MongoDB with native driver |
| **ORM** | GORM auto-migration | MongoDB collections + indexes |
| **Primary Key** | Auto-increment BIGINT | ObjectID (MongoDB default) |
| **Concurrency** | N/A | Optimistic locking with version field |
| **External calls** | None | users-api (validate driver) |
| **Messaging** | None | RabbitMQ consumer + publisher |
| **Idempotency** | N/A | Critical (processed_events collection) |

---

## Success Criteria (Complete Implementation)

✅ Server runs on port 8002
✅ All 6 REST endpoints working
✅ JWT authentication on protected routes
✅ MongoDB connection established
✅ Collections created with proper indexes
✅ RabbitMQ consumer running in background
✅ RabbitMQ publisher emitting events with UUIDs
✅ Idempotency working (duplicate events skipped)
✅ Optimistic locking prevents race conditions
✅ HTTP client calling users-api
✅ Structured logging with zerolog
✅ Health check returns service status
✅ Unit tests passing (especially idempotency and optimistic lock tests)
✅ Docker container builds and runs

---

## References

- **MongoDB Go Driver:** https://www.mongodb.com/docs/drivers/go/current/
- **Gin Framework:** https://gin-gonic.com/docs/
- **RabbitMQ Go:** https://www.rabbitmq.com/tutorials/tutorial-two-go.html
- **users-api code:** `/home/user/CarPooling/backend/users-api`
- **GITFLOW:** See `GITFLOW.md` for workflow

---

## Important Notes

1. **MongoDB is document-based** - No migrations like SQL, just create collections and indexes
2. **ObjectID** is MongoDB's default _id type, converts to string in JSON
3. **Optimistic locking is critical** - Use `availability_version` field
4. **Idempotency with unique index** - MongoDB will reject duplicate `event_id`
5. **Generate UUIDs before publishing** - Never reuse event IDs
6. **Manual ACK in RabbitMQ** - Only after successful processing
7. **Structured logging** - Log all key events with event_id
8. **Error codes, not messages** - Use AppError with codes

---

## Quick Start for Implementation

When starting a feature:

1. **Checkout branch:** `git checkout -b feature/trips-api/{issue-number}-{name}`
2. **Read this context:** Full specification here
3. **Copy patterns:** Reference users-api for similar code (but adapt for MongoDB)
4. **Use plan mode:** `@CONTEXT_TRIPS_API.md Implement feature X`
5. **Test:** Verify success criteria
6. **Commit:** `git commit -m "feat(trips): {description}"`
7. **Push & PR:** Push to branch, create PR to dev
