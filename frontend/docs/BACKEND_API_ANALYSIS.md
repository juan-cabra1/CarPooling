# Backend API Analysis - CarPooling Microservices

## Overview

The CarPooling system consists of 4 microservices built with Go (Gin framework):

1. **users-api** (Port 8001) - MySQL - User management and authentication
2. **trips-api** (Port 8002) - MongoDB - Trip CRUD with optimistic locking
3. **bookings-api** (Port 8003) - MySQL - Reservation management
4. **search-api** (Port 8004) - MongoDB + Solr + Memcached - High-performance search

---

## 1. users-api (Port 8001) - MySQL

### User Model (Database DAO)

```go
type UserDAO struct {
    ID                     int64      `gorm:"primaryKey;autoIncrement"`
    Email                  string     `gorm:"type:varchar(255);unique;not null"`
    EmailVerified          bool       `gorm:"default:false;not null"`
    EmailVerificationToken *string    `gorm:"type:varchar(255)"`
    PasswordResetToken     *string    `gorm:"type:varchar(255)"`
    PasswordResetExpires   *time.Time
    Name                   string     `gorm:"type:varchar(100);not null"`
    Lastname               string     `gorm:"type:varchar(100);not null"`
    PasswordHash           string     `gorm:"type:varchar(255);not null"`
    Role                   string     `gorm:"type:enum('user','admin');default:'user'"`
    Phone                  string     `gorm:"type:varchar(20);not null"`
    Street                 string     `gorm:"type:varchar(255);not null"`
    Number                 int        `gorm:"not null"`
    PhotoURL               string     `gorm:"type:varchar(255)"`
    Sex                    string     `gorm:"type:enum('hombre','mujer','otro')"`
    AvgDriverRating        float64    `gorm:"type:decimal(3,2);default:0.00"`
    AvgPassengerRating     float64    `gorm:"type:decimal(3,2);default:0.00"`
    TotalTripsPassenger    int        `gorm:"default:0"`
    TotalTripsDriver       int        `gorm:"default:0"`
    Birthdate              time.Time  `gorm:"not null"`
    CreatedAt              time.Time  `gorm:"autoCreateTime"`
    UpdatedAt              time.Time  `gorm:"autoUpdateTime"`
}
```

### Enums

- **Role**: `'user'` | `'admin'`
- **Sex**: `'hombre'` | `'mujer'` | `'otro'`
- **RoleRated**: `'conductor'` | `'pasajero'`

### Key Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /users | No | Register new user |
| POST | /login | No | Login (returns JWT) |
| GET | /users/me | Yes | Get current user |
| GET | /users/:id | Yes | Get user by ID |
| PUT | /users/:id | Yes | Update user (owner only) |
| DELETE | /users/:id | Yes | Delete user (owner only) |
| POST | /change-password | Yes | Change password |
| GET | /users/:id/ratings | Yes | Get user ratings |
// Verificación de email y recuperación de contraseña
	router.GET("/verify-email", authController.VerifyEmail)
	router.POST("/resend-verification", authController.ResendVerificationEmail)
	router.POST("/forgot-password", authController.RequestPasswordReset)
	router.POST("/reset-password", authController.ResetPassword)


### Optional Fields (nullable in Go = `| undefined` in TS)

- EmailVerificationToken
- PasswordResetToken
- PasswordResetExpires
- PhotoURL (can be empty string)

---

## 2. trips-api (Port 8002) - MongoDB

### Trip Model

```go
type Trip struct {
    ID                       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    DriverID                 int64
    Origin                   Location
    Destination              Location
    DepartureDatetime        time.Time
    EstimatedArrivalDatetime time.Time
    PricePerSeat             float64
    TotalSeats               int
    ReservedSeats            int
    AvailableSeats           int
    AvailabilityVersion      int // Optimistic locking
    Car                      Car
    Preferences              Preferences
    Status                   string
    Description              string
    CancelledAt              *time.Time
    CancelledBy              *int64
    CancellationReason       string
    CreatedAt                time.Time
    UpdatedAt                time.Time
}

type Location struct {
    City        string
    Province    string
    Address     string
    Coordinates Coordinates
}

type Coordinates struct {
    Lat float64 // Latitude
    Lng float64 // Longitude
}

type Car struct {
    Brand string
    Model string
    Year  int
    Color string
    Plate string
}

type Preferences struct {
    PetsAllowed    bool
    SmokingAllowed bool
    MusicAllowed   bool
}
```

### Enums

- **TripStatus**: `'draft'` | `'published'` | `'full'` | `'in_progress'` | `'completed'` | `'cancelled'`

### Key Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /trips | No | List trips (paginated, filterable) |
| GET | /trips/:id | No | Get trip by ID |
| POST | /trips | Yes | Create trip |
| PUT/PATCH | /trips/:id | Yes | Update trip (owner only) |
| DELETE | /trips/:id | Yes | Delete trip (owner only) |

### Optional Fields

- CancelledAt
- CancelledBy
- CancellationReason

---

## 3. bookings-api (Port 8003) - MySQL

### Booking Model

```go
type Booking struct {
    ID                 uint   `json:"-"` // Internal, not exposed
    BookingUUID        string `json:"id"` // External UUID
    TripID             string // MongoDB ObjectID
    PassengerID        int64
    DriverID           int64
    SeatsRequested     int
    TotalPrice         float64
    Status             string
    CancelledAt        *time.Time
    CancellationReason string
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

### Enums

- **BookingStatus**: `'pending'` | `'confirmed'` | `'cancelled'` | `'completed'` | `'failed'`

### Key Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/v1/bookings | Yes | Create booking |
| GET | /api/v1/bookings | Yes | List user's bookings |
| GET | /api/v1/bookings/:id | Yes | Get booking (owner only) |
| PATCH | /api/v1/bookings/:id/cancel | Yes | Cancel booking |

### Optional Fields

- CancelledAt
- CancellationReason

---

## 4. search-api (Port 8004) - MongoDB + Solr

### SearchTrip Model (Denormalized)

```go
type SearchTrip struct {
    ID                       primitive.ObjectID
    TripID                   string
    DriverID                 int64
    Driver                   Driver // Denormalized from users-api
    Origin                   Location
    Destination              Location
    DepartureDatetime        time.Time
    EstimatedArrivalDatetime time.Time
    PricePerSeat             float64
    TotalSeats               int
    AvailableSeats           int
    Car                      Car
    Preferences              Preferences
    Status                   string
    Description              string
    SearchText               string  // Concatenated for text search
    PopularityScore          float64 // For ranking
    CreatedAt                time.Time
    UpdatedAt                time.Time
}

type Driver struct {
    ID         int64
    Name       string
    Email      string
    PhotoURL   string
    Rating     float64 // avg_driver_rating
    TotalTrips int     // total_trips_driver
}

// GeoJSON for MongoDB 2dsphere indexes
type GeoJSONPoint struct {
    Type        string    // Must be "Point"
    Coordinates []float64 // [longitude, latitude] - lng first!
}
```

### Key Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /api/v1/search/trips | No | Search trips with filters |
| GET | /api/v1/search/location | No | Geospatial search by coordinates |
| GET | /api/v1/trips/:id | No | Get trip details |
| GET | /api/v1/search/autocomplete | No | City autocomplete |
| GET | /api/v1/search/popular-routes | No | Trending routes |

### Search Query Parameters

- `origin_city`, `destination_city` (string)
- `q` (string) - Full-text search
- `min_seats` (int)
- `max_price` (float)
- `min_driver_rating` (float, 0-5)
- `pets_allowed`, `smoking_allowed`, `music_allowed` (bool)
- `date_from`, `date_to` (ISO8601)
- `sort_by`: `popularity` | `price_asc` | `price_desc` | `date_asc` | `date_desc`
- `page` (int, default: 1)
- `limit` (int, default: 20, max: 100)

---

## Global API Patterns

### Standard Response Format

**Success:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message"
  }
}
```

### Pagination Format

```json
{
  "success": true,
  "data": {
    "items": [...],
    "total": 150,
    "page": 1,
    "limit": 10,
    "total_pages": 15
  }
}
```

### Authentication (JWT)

**Header:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Token Claims:**
- `user_id` (int64)
- `email` (string)
- `role` (string)
- `exp` (int64)

### Common HTTP Status Codes

- `200` OK
- `201` Created
- `400` Bad Request (validation errors)
- `401` Unauthorized (missing/invalid token)
- `403` Forbidden (not resource owner)
- `404` Not Found
- `409` Conflict (duplicate, optimistic lock failure)
- `422` Unprocessable Entity (business rule violations)
- `500` Internal Server Error

### Date/Time Format

- Format: RFC3339 (ISO 8601)
- Example: `2025-11-13T10:30:00Z`
- Timezone: UTC recommended
- Go type: `time.Time`

### Coordinate Formats

**trips-api (Simple):**
```json
{
  "lat": 4.7110,
  "lng": -74.0721
}
```

**search-api (GeoJSON):**
```json
{
  "type": "Point",
  "coordinates": [-74.0721, 4.7110]  // [lng, lat] - lng first!
}
```

---

## Go → TypeScript Mapping Rules

| Go Type | TypeScript Type | Notes |
|---------|----------------|-------|
| `time.Time` | `string` | RFC3339/ISO 8601 format |
| `*string` | `string \| undefined` | Pointer = optional |
| `*int`, `*int64` | `number \| undefined` | Pointer = optional |
| `*float64` | `number \| undefined` | Pointer = optional |
| `*bool` | `boolean \| undefined` | Pointer = optional |
| `*time.Time` | `string \| undefined` | Pointer = optional |
| `[]Type` | `Type[]` | Slice to array |
| `map[string]interface{}` | `Record<string, any>` | Generic map |
| `int`, `int64`, `uint` | `number` | All numeric types |
| `float64` | `number` | Floating point |
| `bool` | `boolean` | Boolean |
| `primitive.ObjectID` | `string` | MongoDB ObjectID as hex string |
| `enum (const)` | `'val1' \| 'val2'` | Union of string literals |

---

## Event-Driven Architecture (RabbitMQ)

### trips-api Events

**Publishes:**
- `trip.created`
- `trip.updated`
- `trip.cancelled`
- `trip.completed`

**Consumes:**
- `reservation.created` → Decrement available seats
- `reservation.cancelled` → Increment available seats

### bookings-api Events

**Publishes:**
- `reservation.created`
- `reservation.cancelled`

**Consumes:**
- `trip.cancelled` → Cancel all bookings for trip

### search-api Events

**Consumes:**
- `trip.created` → Index in MongoDB + Solr
- `trip.updated` → Re-index
- `trip.cancelled` → Remove from indexes
- `reservation.confirmed` → Update available seats

---

## API Patterns & Configuration

### Base Paths by Service

- **users-api**: `/` (root - no prefix)
- **trips-api**: `/` (root - no prefix)
- **bookings-api**: `/api/v1` (all routes prefixed)
- **search-api**: `/api/v1` (all routes prefixed)

### Authentication

**Protected Routes (Require JWT):**
- users-api: `/users/me`, `/users/:id`, `/change-password`, `/users/:id/ratings`
- trips-api: `/trips` (POST), `/trips/:id` (PUT/PATCH/DELETE)
- bookings-api: ALL `/api/v1/bookings/*` routes
- search-api: NONE (all public)

**Token Format:**
```
Authorization: Bearer {JWT_TOKEN}
```

**Token Storage:** `localStorage.getItem('token')`

**Token Expiration:** 24 hours from generation

**Token Claims:**
```typescript
{
  user_id: number,  // int64 from backend
  email: string,
  role: string,     // 'user' | 'admin'
  exp: number       // Unix timestamp
}
```

### Response Format Variations

**users-api & trips-api:**
```json
// Success
{ "success": true, "data": { ... } }

// Error (simple string)
{ "success": false, "error": "Error message" }
```

**bookings-api & search-api:**
```json
// Success
{ "success": true, "data": { ... } }

// Error (structured object)
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": "Additional context (optional)"
  }
}
```

### HTTP Status Codes Summary

| Code | Usage |
|------|-------|
| 200 | Successful GET, PUT, PATCH, DELETE |
| 201 | Successful POST (created resource) |
| 204 | OPTIONS preflight (CORS) |
| 400 | Validation errors, business rule violations |
| 401 | Missing/invalid/expired JWT token |
| 403 | Valid token but not resource owner |
| 404 | Resource not found |
| 409 | Duplicate resource, optimistic lock failure |
| 422 | Unprocessable entity (semantic errors) |
| 500 | Internal server error |
| 503 | External service unavailable |

### Vite Proxy Configuration

**Current vite.config.ts proxies:**
```typescript
'/api/users': 'http://localhost:8001'
'/api/trips': 'http://localhost:8002'
'/api/bookings': 'http://localhost:8003'
'/api/search': 'http://localhost:8004'
```

**Note:** This creates a mismatch! Backend actual paths:
- users-api uses `/users`, `/login`, etc. (NO `/api` prefix)
- trips-api uses `/trips` (NO `/api` prefix)
- bookings-api uses `/api/v1/bookings`
- search-api uses `/api/v1/search`

**Recommended Proxy Fix:**

Since Vite proxy rewrites paths, the frontend should call:
- `/api/users/...` → proxied to `http://localhost:8001/...` (path rewrite removes `/api/users`)
- `/api/trips/...` → proxied to `http://localhost:8002/...` (path rewrite removes `/api/trips`)
- `/api/bookings/...` → proxied to `http://localhost:8003/api/v1/...` (rewrite logic needed)
- `/api/search/...` → proxied to `http://localhost:8004/api/v1/...` (rewrite logic needed)

Or, for simplicity, frontend services should use actual backend paths and proxy should pass through:
- Call `/users`, `/login` → proxy to port 8001
- Call `/trips` → proxy to port 8002
- Call `/api/v1/bookings` → proxy to port 8003
- Call `/api/v1/search` → proxy to port 8004

### CORS Configuration

All services configured with:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE
```

Preflight OPTIONS requests return 204.

---

## Summary

This CarPooling system demonstrates:

- **Microservices architecture** with clear service boundaries
- **Polyglot persistence** (MySQL + MongoDB)
- **Event-driven communication** (RabbitMQ)
- **High-performance search** (Solr + Memcached)
- **RESTful APIs** with consistent patterns
- **JWT authentication** across all services
- **Optimistic locking** for concurrency control
- **Denormalization** for performance (search-api)

All services follow Go best practices with repository, service, and controller layers.
