# Repository Layer - Search API

This package contains MongoDB repository implementations for the search-api microservice.

## Overview

The repository layer provides:
- **Full-text search fallback** (when Solr is down)
- **Geospatial queries** using MongoDB's 2dsphere indexes
- **Idempotency checks** with UNIQUE constraint handling
- **Popular routes tracking** and autocomplete data

## Repositories

### 1. TripRepository ([trip_repository.go](trip_repository.go))

Handles all trip search operations with MongoDB.

**Key Methods:**
- `Create(trip *domain.SearchTrip) error` - Creates a new trip document
- `FindByID(id string) (*domain.SearchTrip, error)` - Retrieves a trip by MongoDB ObjectID
- `Update(trip *domain.SearchTrip) error` - Updates an entire trip document
- `UpdateStatus(id string, status string) error` - Updates only the trip status
- `UpdateAvailability(id string, availableSeats int) error` - Updates available seats
- `Search(filters map[string]interface{}, page, limit int) ([]*domain.SearchTrip, int64, error)` - Generic search with pagination
- `SearchByLocation(lat, lng float64, radiusKm int, additionalFilters map[string]interface{}) ([]*domain.SearchTrip, error)` - **Geospatial search** using MongoDB `$near` operator
- `SearchByRoute(originCity, destinationCity string, filters map[string]interface{}) ([]*domain.SearchTrip, error)` - City-to-city route search

**Geospatial Search Example:**
```go
// Search for trips within 50km of Buenos Aires
trips, err := tripRepo.SearchByLocation(
    -34.6037,  // latitude
    -58.3816,  // longitude
    50,        // radius in km
    map[string]interface{}{
        "departure_datetime": bson.M{"$gte": time.Now()},
    },
)
```

**Important Notes:**
- Uses **GeoJSON Point** format: `{type: "Point", coordinates: [lng, lat]}`
- MongoDB expects **[longitude, latitude]** order (lng first!)
- `$near` operator automatically sorts results by distance
- Requires **2dsphere index** on `origin.coordinates` field

### 2. EventRepository ([event_repository.go](event_repository.go))

Handles idempotency tracking for processed events from RabbitMQ.

**Key Methods:**
- `CheckAndMarkEvent(eventID, eventType string) (bool, error)` - **Atomic** check and mark operation
  - Returns `true` if event should be processed (first time)
  - Returns `false` if event is duplicate (already processed)
  - Handles MongoDB duplicate key error (code 11000)
- `IsEventProcessed(eventID string) (bool, error)` - Check if event was processed
- `MarkEventProcessed(event *domain.ProcessedEvent) error` - Mark event as processed (idempotent)

**Idempotency Example:**
```go
shouldProcess, err := eventRepo.CheckAndMarkEvent(ctx, "event-uuid-123", "trip.created")
if err != nil {
    return err
}

if !shouldProcess {
    log.Info().Msg("Event already processed, skipping")
    return nil
}

// Process the event...
```

**Important Notes:**
- Uses **UNIQUE index** on `event_id` field (critical for idempotency)
- `CheckAndMarkEvent` is **atomic** - prevents race conditions
- Safe for concurrent processing (multiple consumers)

### 3. PopularRouteRepository ([popular_route_repository.go](popular_route.go))

Tracks popular routes for trending and analytics.

**Key Methods:**
- `IncrementSearchCount(originCity, destinationCity string) error` - Increments search count (upsert)
- `GetPopularRoutes(limit int) ([]*domain.PopularRoute, error)` - Gets top N popular routes
- `GetAutocompleteSuggestions(prefix string, limit int) ([]string, error)` - City name autocomplete

**Popular Routes Example:**
```go
// Track a search
err := popularRouteRepo.IncrementSearchCount(ctx, "Buenos Aires", "Córdoba")

// Get top 10 popular routes
popularRoutes, err := popularRouteRepo.GetPopularRoutes(ctx, 10)

// Autocomplete cities starting with "Bue"
cities, err := popularRouteRepo.GetAutocompleteSuggestions(ctx, "Bue", 5)
// Returns: ["buenos aires", ...]
```

**Important Notes:**
- Uses **UNIQUE compound index** on `(origin_city, destination_city)`
- Upsert operation creates route if it doesn't exist
- Autocomplete returns **lowercase** city names for consistency
- Autocomplete searches both `origin_city` and `destination_city` fields

## MongoDB Indexes

All indexes are created in [internal/database/mongodb.go](../database/mongodb.go) via `CreateIndexes()`:

### Trips Collection
1. **Compound index**: `(status, departure_datetime)` - For filtering published trips by date
2. **Compound index**: `(origin.city, destination.city)` - For city-to-city searches
3. **2dsphere index**: `origin.coordinates` - **For geospatial queries** (required for `$near`)
4. **2dsphere index**: `destination.coordinates` - For destination geospatial queries

### Processed Events Collection
1. **UNIQUE index**: `event_id` - **Critical for idempotency**
2. **Index**: `event_type` - For filtering by event type
3. **Index**: `processed_at` - For timestamp queries

### Popular Routes Collection
1. **UNIQUE compound index**: `(origin_city, destination_city)` - Ensures one record per route

## GeoJSON Coordinates Format

The domain model uses **GeoJSON Point** format as required by MongoDB 2dsphere indexes:

```go
type GeoJSONPoint struct {
    Type        string    `bson:"type"`          // Always "Point"
    Coordinates []float64 `bson:"coordinates"`   // [longitude, latitude]
}
```

**Creating GeoJSON Points:**
```go
// Helper function
coordinates := domain.NewGeoJSONPoint(lat, lng)

// Manual creation
coordinates := domain.GeoJSONPoint{
    Type:        "Point",
    Coordinates: []float64{lng, lat}, // IMPORTANT: [lng, lat] order!
}
```

**Accessing Coordinates:**
```go
lat := location.Coordinates.Lat()  // Returns latitude
lng := location.Coordinates.Lng()  // Returns longitude
```

## Testing

All repositories have comprehensive unit tests:
- [trip_repository_test.go](trip_repository_test.go) - Trip repository tests including geospatial
- [event_repository_test.go](event_repository_test.go) - Idempotency and concurrency tests
- [popular_route_repository_test.go](popular_route_repository_test.go) - Popular routes and autocomplete tests

**Run tests:**
```bash
# Run all repository tests
go test ./internal/repository/...

# Run with verbose output
go test -v ./internal/repository/...

# Run specific test
go test -v -run TestTripRepository_SearchByLocation ./internal/repository/
```

**Test Requirements:**
- MongoDB running on `localhost:27017`
- Tests use separate test databases (`search_api_test`, `search_api_test_events`, `search_api_test_routes`)
- Tests automatically clean up data after execution

## Integration with Other Microservices

### Event Processing Flow
1. **trips-api** publishes event to RabbitMQ (e.g., `trip.created`)
2. **search-api** consumer receives event
3. **EventRepository** checks idempotency: `CheckAndMarkEvent(eventID, eventType)`
4. If `shouldProcess == true`:
   - Fetch driver data from **users-api** (denormalize)
   - Create/update trip in MongoDB via **TripRepository**
   - Index in Solr (if available)
5. If `shouldProcess == false`: Skip (duplicate event)

### Search Flow
1. User searches for trips (e.g., Buenos Aires → La Plata)
2. **PopularRouteRepository** tracks search: `IncrementSearchCount(origin, destination)`
3. **TripRepository** performs search:
   - Primary: Solr full-text search
   - Fallback: MongoDB `SearchByRoute()` or `SearchByLocation()`
4. Return results to client

## Performance Considerations

1. **Geospatial Queries**: Require 2dsphere index - ensure index exists before production
2. **Pagination**: Use `skip` and `limit` for large result sets
3. **Idempotency**: Unique index on `event_id` ensures O(1) duplicate detection
4. **Popular Routes**: Upsert with `$inc` is atomic and efficient

## References

- [MongoDB Geospatial Queries](https://www.mongodb.com/docs/manual/geospatial-queries/)
- [MongoDB 2dsphere Indexes](https://www.mongodb.com/docs/manual/core/2dsphere/)
- [GeoJSON Specification](https://www.mongodb.com/docs/manual/reference/geojson/)
- [MongoDB $near Operator](https://www.mongodb.com/docs/manual/reference/operator/query/near/)
