# Search API - Implementation Context

## Project Overview
**Academic Project** - CarPooling platform with microservices architecture.

### Requirements
- ✅ MongoDB for search-api (this service) - stores denormalized trip data
- ✅ Apache Solr for full-text and geospatial search
- ✅ Memcached for caching search results
- ✅ MySQL for users-api and bookings-api
- ✅ RabbitMQ for event-driven communication

### Architecture
```
users-api:8001 (MySQL) → ✅ DONE
trips-api:8002 (MongoDB) → ✅ DONE
bookings-api:8003 (MySQL) → ✅ DONE
search-api:8004 (MongoDB+Solr+Memcached) → TO IMPLEMENT (THIS SERVICE)
```

---

## Search API - Specifications

### Technology Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP)
- **Database:** MongoDB with native driver (denormalized trip data)
- **Search Engine:** Apache Solr 9.x (full-text + geospatial indexing)
- **Cache:** Memcached (search results caching)
- **Messaging:** RabbitMQ (streadway/amqp)
- **Auth:** JWT (same secret as other services)
- **Logging:** zerolog (structured JSON logs)
- **UUID:** google/uuid v4

### Purpose and Design Philosophy

**Why search-api exists:**
- **Performance:** Offload complex searches from trips-api
- **Scalability:** Dedicated indexing engine (Solr) for advanced queries
- **Caching:** Reduce database load with Memcached
- **Geospatial:** Efficient lat/lng radius searches with Solr
- **Full-text:** Search in trip descriptions, cities, addresses
- **Eventual Consistency:** Async updates via RabbitMQ don't block trip creation

**Data Flow:**
```
trips-api (source of truth)
  → publishes trip.created/updated/cancelled
  → search-api consumes events
  → updates MongoDB (backup) + Solr (indexing) + invalidate cache
```

---

## Database: MongoDB `carpooling_search`

**Collection: trips (Denormalized)**

```javascript
{
  _id: String,  // Same as trips-api ObjectId (for correlation)
  driver_id: Number,

  // Denormalized driver info (from users-api)
  driver: {
    name: String,
    email: String,
    rating: Number,
    total_trips: Number
  },

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
  available_seats: Number,
  reserved_seats: Number,

  car: {
    brand: String,
    model: String,
    year: Number,
    color: String
  },

  preferences: {
    pets_allowed: Boolean,
    smoking_allowed: Boolean,
    music_allowed: Boolean
  },

  status: String,  // 'published', 'full', 'in_progress', 'completed', 'cancelled'
  description: String,

  // Search-specific fields
  search_text: String,  // Concatenated: cities + description (for full-text)
  popularity_score: Number,  // Based on views, bookings (for ranking)

  created_at: Date,
  updated_at: Date,
  indexed_at: Date,  // Last time indexed in Solr
  solr_synced: Boolean  // Is Solr in sync?
}
```

**Indexes:**
```javascript
db.trips.createIndex({ status: 1, departure_datetime: 1 })
db.trips.createIndex({ "origin.city": 1, "destination.city": 1 })
db.trips.createIndex({ "origin.coordinates": "2dsphere" })
db.trips.createIndex({ "destination.coordinates": "2dsphere" })
db.trips.createIndex({ driver_id: 1 })
db.trips.createIndex({ solr_synced: 1 })
```

**Collection: processed_events (CRITICAL - Idempotency)**
```javascript
{
  _id: ObjectId(),
  event_id: String,        // UUID from RabbitMQ - UNIQUE INDEX
  event_type: String,      // 'trip.created', 'trip.updated', 'trip.cancelled'
  processed_at: Date,
  result: String,          // 'success', 'skipped', 'failed'
  error_message: String
}
```

**Indexes:**
```javascript
db.processed_events.createIndex({ event_id: 1 }, { unique: true })  // CRITICAL
```

**Collection: popular_routes (for autocomplete and trending)**
```javascript
{
  _id: ObjectId(),
  origin_city: String,
  destination_city: String,
  search_count: Number,
  trip_count: Number,
  last_searched: Date
}
```

**Indexes:**
```javascript
db.popular_routes.createIndex({ origin_city: 1, destination_city: 1 }, { unique: true })
db.popular_routes.createIndex({ search_count: -1 })
```

---

## Apache Solr Schema

**Core:** `carpooling_trips`

**Fields:**
```xml
<field name="id" type="string" indexed="true" stored="true" required="true" />
<field name="driver_id" type="plong" indexed="true" stored="true" />
<field name="driver_name" type="text_general" indexed="true" stored="true" />
<field name="driver_rating" type="pfloat" indexed="true" stored="true" />

<field name="origin_city" type="string" indexed="true" stored="true" />
<field name="origin_province" type="string" indexed="true" stored="true" />
<field name="origin_location" type="location" indexed="true" stored="true" />

<field name="destination_city" type="string" indexed="true" stored="true" />
<field name="destination_province" type="string" indexed="true" stored="true" />
<field name="destination_location" type="location" indexed="true" stored="true" />

<field name="departure_datetime" type="pdate" indexed="true" stored="true" />
<field name="price_per_seat" type="pfloat" indexed="true" stored="true" />
<field name="available_seats" type="pint" indexed="true" stored="true" />

<field name="pets_allowed" type="boolean" indexed="true" stored="true" />
<field name="smoking_allowed" type="boolean" indexed="true" stored="true" />
<field name="music_allowed" type="boolean" indexed="true" stored="true" />

<field name="status" type="string" indexed="true" stored="true" />
<field name="description" type="text_general" indexed="true" stored="true" />
<field name="search_text" type="text_general" indexed="true" stored="false" />

<field name="popularity_score" type="pfloat" indexed="true" stored="true" />
<field name="created_at" type="pdate" indexed="true" stored="true" />
```

**Geospatial Search:**
```
origin_location: "lat,lng"
destination_location: "lat,lng"
```

Query example: Find trips within 50km of Buenos Aires:
```
{!geofilt pt=-34.6037,-58.3816 sfield=origin_location d=50}
```

---

## REST Endpoints (Port 8004)

### Search Endpoints

```
GET    /search/trips              [Public]  Advanced search with filters
GET    /search/autocomplete       [Public]  City autocomplete
GET    /search/popular-routes     [Public]  Trending routes
GET    /trips/:id                 [Public]  Get trip by ID (from MongoDB)
GET    /health                    [Public]  Health check (includes Solr + Memcached status)
```

### Request/Response Examples

**GET /search/trips** (Advanced Search)
```
Query Parameters:
  - origin_city: string (e.g., "Buenos Aires")
  - destination_city: string (e.g., "Rosario")
  - origin_lat, origin_lng, origin_radius_km: geospatial search
  - destination_lat, destination_lng, destination_radius_km: geospatial search
  - date_from, date_to: departure date range (ISO8601)
  - min_seats: minimum available seats
  - max_price: maximum price per seat
  - pets_allowed: boolean
  - smoking_allowed: boolean
  - music_allowed: boolean
  - min_driver_rating: float (e.g., 4.0)
  - sort_by: "price_asc", "price_desc", "date_asc", "date_desc", "rating_desc", "distance_asc"
  - page: int (default 1)
  - limit: int (default 20, max 100)
```

Example Request:
```
GET /search/trips?origin_city=Buenos+Aires&destination_city=Rosario&date_from=2024-12-01&min_seats=2&max_price=6000&sort_by=price_asc&page=1&limit=10
```

Response (200):
```json
{
  "success": true,
  "data": {
    "trips": [
      {
        "id": "507f1f77bcf86cd799439011",
        "driver": {
          "id": 123,
          "name": "Juan Pérez",
          "rating": 4.8,
          "total_trips": 45
        },
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
        "price_per_seat": 5000,
        "available_seats": 3,
        "car": {
          "brand": "Toyota",
          "model": "Corolla",
          "year": 2020,
          "color": "Gris"
        },
        "preferences": {
          "pets_allowed": false,
          "smoking_allowed": false,
          "music_allowed": true
        },
        "description": "Viaje tranquilo a Rosario"
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 10,
    "cached": false,
    "query_time_ms": 45
  }
}
```

**GET /search/autocomplete?q=bue** (City Autocomplete)
```json
Response (200):
{
  "success": true,
  "data": {
    "suggestions": [
      { "city": "Buenos Aires", "province": "Buenos Aires", "trip_count": 234 },
      { "city": "Buena Esperanza", "province": "San Luis", "trip_count": 5 }
    ]
  }
}
```

**GET /search/popular-routes?limit=10** (Trending Routes)
```json
Response (200):
{
  "success": true,
  "data": {
    "routes": [
      {
        "origin_city": "Buenos Aires",
        "destination_city": "Rosario",
        "trip_count": 45,
        "search_count": 1200,
        "last_searched": "2024-11-10T12:00:00Z"
      },
      {
        "origin_city": "Córdoba",
        "destination_city": "Buenos Aires",
        "trip_count": 38,
        "search_count": 890,
        "last_searched": "2024-11-10T11:30:00Z"
      }
    ]
  }
}
```

**GET /trips/:id** (Get Single Trip)
```json
Response (200):
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "driver": { ... },
    "origin": { ... },
    "destination": { ... },
    "cached": true
  }
}
```

---

## RabbitMQ Events

### Consumed Events (from trips-api):

**trip.created** (New trip published)
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.created",
  "trip_id": "507f1f77bcf86cd799439011",
  "driver_id": 123,
  "origin_city": "Buenos Aires",
  "destination_city": "Rosario",
  "departure_datetime": "2024-12-01T10:00:00Z",
  "available_seats": 3,
  "status": "published",
  "timestamp": "2024-11-10T10:00:00Z"
}
```

**trip.updated** (Trip availability or details changed)
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.updated",
  "trip_id": "507f1f77bcf86cd799439011",
  "available_seats": 1,
  "reserved_seats": 2,
  "status": "published",
  "timestamp": "2024-11-10T10:05:00Z"
}
```

**trip.cancelled** (Trip cancelled by driver)
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.cancelled",
  "trip_id": "507f1f77bcf86cd799439011",
  "cancellation_reason": "Car breakdown",
  "timestamp": "2024-11-10T12:00:00Z"
}
```

### Published Events: None
search-api is a **consumer-only** service. It does not publish events.

---

## Event Flow & Data Synchronization

### NEW TRIP CREATED:
```
1. User creates trip in trips-api
2. trips-api publishes trip.created event
3. search-api consumes event (with idempotency check)
4. search-api fetches FULL trip details from trips-api (HTTP GET)
5. search-api fetches driver info from users-api (HTTP GET)
6. search-api stores denormalized document in MongoDB
7. search-api indexes document in Solr
8. search-api marks solr_synced = true
```

### TRIP UPDATED (availability changed):
```
1. trips-api publishes trip.updated event
2. search-api consumes event (with idempotency check)
3. search-api updates MongoDB (available_seats, status)
4. search-api updates Solr document
5. search-api invalidates cache for that trip_id
```

### TRIP CANCELLED:
```
1. trips-api publishes trip.cancelled event
2. search-api consumes event (with idempotency check)
3. search-api updates MongoDB (status = 'cancelled')
4. search-api updates Solr document (status = 'cancelled')
5. search-api invalidates cache for that trip_id
```

### IDEMPOTENCY SCENARIO:
```
1. search-api consumes trip.created (event_id: abc-123)
2. Processing succeeds, trip indexed in Solr
3. ACK fails, RabbitMQ retries same message
4. search-api receives SAME event (event_id: abc-123)
5. Idempotency check: event already in processed_events
6. SKIP processing, just ACK
7. No duplicate indexing ✅
```

---

## Caching Strategy (Memcached)

### Cache Keys:
```
trip:{trip_id}                         → Single trip (TTL: 5 minutes)
search:{hash}                          → Search results (TTL: 2 minutes)
autocomplete:{query}                   → Autocomplete suggestions (TTL: 1 hour)
popular_routes                         → Popular routes (TTL: 15 minutes)
```

### Cache Invalidation:
- **On trip.updated:** Invalidate `trip:{trip_id}` and all `search:*` keys for that city pair
- **On trip.cancelled:** Invalidate `trip:{trip_id}` and all `search:*` keys
- **On trip.created:** No invalidation needed (new data)

### Cache-Aside Pattern:
```go
func SearchTrips(query SearchQuery) ([]Trip, error) {
    cacheKey := fmt.Sprintf("search:%s", query.Hash())

    // 1. Try cache
    if cached, found := cache.Get(cacheKey); found {
        return cached, nil
    }

    // 2. Query Solr
    results, err := solr.Search(query)
    if err != nil {
        return nil, err
    }

    // 3. Store in cache
    cache.Set(cacheKey, results, 2*time.Minute)

    return results, nil
}
```

---

## Critical Requirements

### 1. IDEMPOTENCY (Prevent Duplicate Indexing)

**Problem:** RabbitMQ retries → duplicate events → duplicate documents in Solr

**Solution:**
```go
func HandleTripEvent(event TripEvent) error {
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
    case "trip.created":
        return handleTripCreated(event)
    case "trip.updated":
        return handleTripUpdated(event)
    case "trip.cancelled":
        return handleTripCancelled(event)
    }
}
```

### 2. DENORMALIZATION (Driver Info from users-api)

**Why:** Search results should include driver name, rating, total trips without N+1 queries

**Solution:**
```go
func HandleTripCreated(event TripCreatedEvent) error {
    // 1. Fetch full trip from trips-api
    trip, err := tripsClient.GetTrip(event.TripID)
    if err != nil {
        return err
    }

    // 2. Fetch driver info from users-api (IMPORTANT)
    driver, err := usersClient.GetUser(trip.DriverID)
    if err != nil {
        return err
    }

    // 3. Store denormalized document in MongoDB
    searchTrip := domain.SearchTrip{
        ID:       trip.ID,
        DriverID: trip.DriverID,
        Driver: domain.Driver{
            Name:       driver.Name,
            Email:      driver.Email,
            Rating:     driver.Rating,
            TotalTrips: driver.TotalTrips,
        },
        Origin:      trip.Origin,
        Destination: trip.Destination,
        // ... rest of fields
        SearchText:  buildSearchText(trip),
        SolrSynced:  false,
    }

    if err := repo.Create(searchTrip); err != nil {
        return err
    }

    // 4. Index in Solr
    if err := solrClient.Index(searchTrip); err != nil {
        logger.Error().Err(err).Msg("Solr indexing failed, will retry later")
        // Don't return error - will be indexed by background job
    }

    // 5. Mark as synced
    repo.MarkSolrSynced(searchTrip.ID, true)

    return nil
}
```

### 3. GEOSPATIAL SEARCH (Solr + MongoDB)

**Problem:** Find trips within X km of user location

**Solution (Solr):**
```go
func SearchByLocation(lat, lng float64, radiusKm int) ([]Trip, error) {
    query := fmt.Sprintf(
        "{!geofilt pt=%f,%f sfield=origin_location d=%d}",
        lat, lng, radiusKm,
    )

    return solr.Query(query)
}
```

**Solution (MongoDB 2dsphere fallback):**
```go
filter := bson.M{
    "origin.coordinates": bson.M{
        "$near": bson.M{
            "$geometry": bson.M{
                "type":        "Point",
                "coordinates": []float64{lng, lat},
            },
            "$maxDistance": radiusKm * 1000, // meters
        },
    },
}
```

### 4. POPULARITY RANKING

**Algorithm:**
```go
popularity_score = (bookings_count * 10) + (views_count * 1) + (driver_rating * 5)
```

Used for default sorting when no explicit sort order provided.

---

## Architecture Pattern

```
internal/
├── config/              # Environment variables
│   └── config.go
├── domain/              # Business models (DTOs)
│   ├── trip.go          # SearchTrip (denormalized)
│   ├── search_query.go
│   ├── driver.go
│   └── errors.go
├── repository/          # Data access interfaces + implementations
│   ├── trip_repository.go       # MongoDB
│   └── event_repository.go      # MongoDB (idempotency)
├── service/             # Business logic
│   ├── search_service.go        # Main search logic
│   ├── sync_service.go          # Solr sync background job
│   ├── idempotency_service.go
│   └── auth_service.go
├── controller/          # HTTP handlers
│   ├── search_controller.go
│   └── trip_controller.go
├── middleware/          # JWT, CORS, errors
│   ├── auth.go
│   ├── cors.go
│   └── error.go
├── messaging/           # RabbitMQ
│   ├── rabbitmq.go
│   └── trip_consumer.go
├── http/                # External HTTP clients
│   ├── trips_client.go
│   └── users_client.go
├── solr/                # Apache Solr client
│   ├── client.go
│   └── schema.go
├── cache/               # Memcached client
│   └── client.go
└── routes/
    └── routes.go
cmd/api/
└── main.go             # Dependency injection, server start
cmd/sync/               # Background Solr sync job
└── main.go
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
    github.com/rtt/Go-Solr v0.0.0-20190813145837-4c4c2e62a3fd  // Solr client
    github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874  // Memcached
)
```

---

## Environment Variables (.env)

```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=carpooling_search

# Solr
SOLR_URL=http://localhost:8983/solr/carpooling_trips

# Memcached
MEMCACHED_SERVERS=localhost:11211

# JWT (same secret as other services)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
SERVER_PORT=8004

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# External APIs
TRIPS_API_URL=http://localhost:8002
USERS_API_URL=http://localhost:8001

# Search Config
SEARCH_DEFAULT_LIMIT=20
SEARCH_MAX_LIMIT=100
CACHE_TTL_SEARCH=2m
CACHE_TTL_TRIP=5m
CACHE_TTL_AUTOCOMPLETE=1h
```

---

## Business Logic

### Search Priority:
1. **Try Memcached** - Return cached results if available
2. **Query Solr** - Use Solr for complex queries (geospatial, full-text, facets)
3. **Fallback to MongoDB** - If Solr is down, query MongoDB directly (slower)

### Filtering Rules:
- Only show trips with `status = 'published'`
- Only show trips with `available_seats > 0` (unless explicitly requested)
- Only show trips with `departure_datetime >= now()` (future trips)

### Sorting Options:
- `price_asc`: Price low to high
- `price_desc`: Price high to low
- `date_asc`: Soonest departure
- `date_desc`: Latest departure
- `rating_desc`: Highest driver rating
- `distance_asc`: Closest origin to user location (requires lat/lng)
- `popularity_desc`: Highest popularity score (default)

### Autocomplete Logic:
- Match prefix on `origin.city` or `destination.city`
- Return max 10 suggestions
- Order by `trip_count DESC`
- Cache for 1 hour

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
    ErrSolrUnavailable   = &AppError{Code: "SOLR_UNAVAILABLE", Message: "Search engine temporarily unavailable"}
    ErrCacheUnavailable  = &AppError{Code: "CACHE_UNAVAILABLE", Message: "Cache temporarily unavailable"}
    ErrInvalidQuery      = &AppError{Code: "INVALID_QUERY", Message: "Invalid search query"}
    ErrInvalidGeoCoords  = &AppError{Code: "INVALID_GEO_COORDS", Message: "Invalid coordinates"}
)
```

---

## Testing Focus

### Unit Tests (service layer)
```go
// Test: Search with filters
func TestSearch_WithFilters(t *testing.T) {
    query := SearchQuery{
        OriginCity: "Buenos Aires",
        DestinationCity: "Rosario",
        MinSeats: 2,
        MaxPrice: 6000,
    }

    results, err := service.Search(query)

    assert.NoError(t, err)
    assert.GreaterOrEqual(t, len(results), 1)
    for _, trip := range results {
        assert.Equal(t, "Buenos Aires", trip.Origin.City)
        assert.GreaterOrEqual(t, trip.AvailableSeats, 2)
        assert.LessOrEqual(t, trip.PricePerSeat, 6000.0)
    }
}

// Test: Geospatial search
func TestSearch_Geospatial(t *testing.T) {
    // Search within 50km of Buenos Aires center
    results, err := service.SearchByLocation(-34.6037, -58.3816, 50)

    assert.NoError(t, err)
    // Verify all results are within radius
}

// Test: Cache hit
func TestSearch_CacheHit(t *testing.T) {
    query := SearchQuery{OriginCity: "Buenos Aires"}

    // First call - cache miss
    results1, err := service.Search(query)
    assert.NoError(t, err)

    // Second call - cache hit
    results2, err := service.Search(query)
    assert.NoError(t, err)
    assert.Equal(t, results1, results2)
}
```

### Idempotency Tests (CRITICAL)
```go
// Test: Duplicate event skipped
func TestIdempotency_DuplicateEventSkipped(t *testing.T) {
    event := TripCreatedEvent{
        EventID: "event-123",
        TripID:  "trip-456",
    }

    // First processing
    err := handler.HandleTripCreated(event)
    assert.NoError(t, err)

    // Verify trip indexed in Solr
    trip, err := solr.Get("trip-456")
    assert.NoError(t, err)
    assert.NotNil(t, trip)

    // Second processing (duplicate)
    err = handler.HandleTripCreated(event)
    assert.NoError(t, err)

    // Verify NOT indexed twice
    count := solr.Count(`id:"trip-456"`)
    assert.Equal(t, 1, count)
}
```

### Integration Tests
```go
// Test: Full event flow
func TestEventFlow_TripCreatedToIndexed(t *testing.T) {
    // 1. Publish trip.created event
    // 2. Wait for consumer to process
    // 3. Verify trip in MongoDB
    // 4. Verify trip in Solr
    // 5. Verify searchable via API
}
```

---

## Background Jobs

### Solr Sync Job (cmd/sync/main.go)

**Purpose:** Re-index trips where `solr_synced = false` (failed indexing)

**Cron:** Every 5 minutes

```go
func main() {
    ticker := time.NewTicker(5 * time.Minute)

    for range ticker.C {
        trips, err := repo.FindUnsyncedTrips()
        if err != nil {
            logger.Error().Err(err).Msg("Failed to find unsynced trips")
            continue
        }

        for _, trip := range trips {
            if err := solr.Index(trip); err != nil {
                logger.Error().Err(err).Str("trip_id", trip.ID).Msg("Solr indexing failed")
            } else {
                repo.MarkSolrSynced(trip.ID, true)
                logger.Info().Str("trip_id", trip.ID).Msg("Trip re-indexed in Solr")
            }
        }
    }
}
```

---

## Apache Solr Configuration

### Core Creation:
```bash
# Create core
bin/solr create -c carpooling_trips

# Define schema
curl -X POST http://localhost:8983/solr/carpooling_trips/schema \
  -H 'Content-Type: application/json' \
  -d '{
    "add-field": [
      {"name": "driver_id", "type": "plong", "stored": true, "indexed": true},
      {"name": "origin_city", "type": "string", "stored": true, "indexed": true},
      {"name": "origin_location", "type": "location", "stored": true, "indexed": true},
      {"name": "destination_city", "type": "string", "stored": true, "indexed": true},
      {"name": "destination_location", "type": "location", "stored": true, "indexed": true},
      {"name": "departure_datetime", "type": "pdate", "stored": true, "indexed": true},
      {"name": "price_per_seat", "type": "pfloat", "stored": true, "indexed": true},
      {"name": "available_seats", "type": "pint", "stored": true, "indexed": true},
      {"name": "status", "type": "string", "stored": true, "indexed": true},
      {"name": "pets_allowed", "type": "boolean", "stored": true, "indexed": true},
      {"name": "smoking_allowed", "type": "boolean", "stored": true, "indexed": true},
      {"name": "search_text", "type": "text_general", "stored": false, "indexed": true}
    ]
  }'
```

---

## Success Criteria (Complete Implementation)

✅ Server runs on port 8004
✅ All 5 REST endpoints working
✅ MongoDB connection established
✅ Solr connection established and core created
✅ Memcached connection established
✅ RabbitMQ consumer running in background (consumes from trips-api)
✅ Idempotency working (duplicate events skipped)
✅ HTTP clients calling trips-api and users-api
✅ Search with filters returns correct results
✅ Geospatial search within radius working
✅ Autocomplete suggestions working
✅ Cache hit/miss working correctly
✅ Cache invalidation on trip updates
✅ Solr sync background job running
✅ Structured logging with zerolog
✅ Health check returns Solr + Memcached status
✅ Unit tests passing (especially idempotency)
✅ Docker container builds and runs with Solr + Memcached

---

## References

- **MongoDB Go Driver:** https://www.mongodb.com/docs/drivers/go/current/
- **Apache Solr:** https://solr.apache.org/guide/solr/latest/
- **Gomemcache:** https://github.com/bradfitz/gomemcache
- **Gin Framework:** https://gin-gonic.com/docs/
- **RabbitMQ Go:** https://www.rabbitmq.com/tutorials/tutorial-two-go.html
- **trips-api code:** `/home/user/CarPooling/backend/trips-api` (MongoDB + RabbitMQ consumer patterns)
- **GITFLOW:** See `GITFLOW.md` for workflow

---

## Important Notes

1. **Eventual Consistency:** search-api may be slightly behind trips-api (acceptable)
2. **Denormalization:** Store driver info to avoid N+1 queries
3. **Solr as Primary:** Use Solr for search, MongoDB as backup/storage
4. **Cache Aggressively:** Search results change slowly, cache for 2+ minutes
5. **Idempotency Critical:** Prevent duplicate Solr documents
6. **Graceful Degradation:** If Solr down, fallback to MongoDB (slower but works)
7. **Geospatial:** Solr's geofilt is faster than MongoDB's 2dsphere for large datasets
8. **Background Sync:** Re-index failed trips every 5 minutes
9. **No Events Published:** search-api is consumer-only
10. **Driver Info Freshness:** Consider periodic re-sync if driver ratings change frequently

---

## Quick Start for Implementation

When starting a feature:

1. **Checkout branch:** `git checkout -b feature/search-api/{issue-number}-{name}`
2. **Read this context:** Full specification here
3. **Copy patterns:** Reference trips-api for MongoDB + RabbitMQ, bookings-api for HTTP clients
4. **Use plan mode:** `@CONTEXT_SEARCH_API.md Implement feature X`
5. **Test:** Verify success criteria
6. **Commit:** `git commit -m "feat(search): {description}"`
7. **Push & PR:** Push to branch, create PR to dev