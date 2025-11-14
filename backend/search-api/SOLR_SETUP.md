# Apache Solr Integration - Setup Guide

## Overview

The search-api now integrates Apache Solr for advanced search capabilities including:
- ✅ Full-text search on trip descriptions
- ✅ Complex filtering (city, price, seats, dates, preferences)
- ✅ Faceted search for aggregations
- ✅ High-performance sorting and ranking
- ⚠️  Geospatial search is handled by MongoDB (not Solr)

## Architecture

```
┌─────────────┐
│   User      │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  search-api │
└──────┬──────┘
       │
       ├─────────────┐
       │             │
       ▼             ▼
┌─────────────┐  ┌─────────────┐
│  Apache     │  │   MongoDB   │
│  Solr       │  │             │
│  (Search)   │  │  (Storage)  │
└─────────────┘  └─────────────┘
```

- **MongoDB**: Primary data storage + geospatial queries
- **Solr**: Advanced text search, facets, complex filtering
- **Memcached**: Query result caching
- **RabbitMQ**: Event-driven indexing

## Prerequisites

1. Docker and Docker Compose installed
2. Solr container running (included in docker-compose.yml)
3. MongoDB container running
4. Memcached container running (optional, for caching)

## Setup Instructions

### Step 1: Start Infrastructure

From the `search-api` directory:

```bash
cd backend/search-api
docker-compose up -d
```

This starts:
- MongoDB (port 27017)
- Apache Solr (port 8983)
- Memcached (port 11211)
- RabbitMQ (port 5672)

### Step 2: Verify Solr is Running

Open in your browser:
```
http://localhost:8983/solr
```

You should see the Solr Admin UI.

### Step 3: Create Solr Core and Schema

Run the setup script:

```bash
# Make script executable (Linux/Mac)
chmod +x scripts/setup_solr_schema.sh

# Run the script
./scripts/setup_solr_schema.sh
```

For Windows (Git Bash or WSL):
```bash
bash scripts/setup_solr_schema.sh
```

The script will:
1. Check if Solr is accessible
2. Delete existing core (if any)
3. Create new core named `carpooling_trips`
4. Add all required fields to the schema
5. Test the schema with a sample query

### Step 4: Start search-api

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/api/main.go
```

You should see:
```
✅ Configuration loaded successfully
✅ MongoDB indexes created successfully
✅ Connected to Apache Solr successfully
✅ Connected to Memcached successfully
✅ Repositories initialized successfully
✅ Controllers initialized successfully
✅ Routes configured successfully
✅ search-api server listening on port 8004
```

### Step 5: Verify Health Check

```bash
curl http://localhost:8004/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "search-api",
  "port": "8004",
  "services": {
    "mongodb": {
      "status": "healthy",
      "message": "Connected"
    },
    "solr": {
      "status": "healthy",
      "message": "Connected"
    },
    "memcached": {
      "status": "healthy",
      "message": "Connected"
    }
  }
}
```

## Solr Schema

### Indexed Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Trip ID (unique key) |
| `driver_id` | plong | Driver user ID |
| `driver_name` | text_general | Driver full name |
| `driver_rating` | pfloat | Driver rating (0.0-5.0) |
| `driver_total_trips` | pint | Total trips completed |
| `origin_city` | string | Origin city name |
| `origin_province` | string | Origin province |
| `destination_city` | string | Destination city name |
| `destination_province` | string | Destination province |
| `departure_datetime` | pdate | Departure date/time |
| `estimated_arrival_datetime` | pdate | Arrival date/time |
| `price_per_seat` | pfloat | Price per seat |
| `total_seats` | pint | Total seats in vehicle |
| `available_seats` | pint | Available seats |
| `car_brand` | text_general | Vehicle brand |
| `car_model` | text_general | Vehicle model |
| `car_year` | pint | Vehicle year |
| `car_color` | string | Vehicle color |
| `pets_allowed` | boolean | Pets allowed? |
| `smoking_allowed` | boolean | Smoking allowed? |
| `music_allowed` | boolean | Music allowed? |
| `status` | string | Trip status |
| `description` | text_general | Trip description |
| `search_text` | text_general | Concatenated search text |
| `popularity_score` | pfloat | Popularity ranking |
| `created_at` | pdate | Creation timestamp |
| `updated_at` | pdate | Last update timestamp |

### Field Types

- `string`: Exact matching, not tokenized
- `text_general`: Full-text search, tokenized
- `plong`: Long integer (point field for range queries)
- `pint`: Integer (point field for range queries)
| `pfloat`: Float (point field for range queries)
- `pdate`: Date/time (ISO 8601 format)
- `boolean`: True/false

## Usage Examples

### Example 1: Index a Trip

```go
import "search-api/internal/solr"

// Create Solr client
solrClient, err := solr.NewClient("http://localhost:8983/solr", "carpooling_trips")

// Index a trip
err = solrClient.Index(searchTrip)
```

### Example 2: Full-Text Search

```go
query := &solr.SearchQuery{
    Query: "viaje cómodo confiable", // Search in description and search_text
    Page:  1,
    Limit: 10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=(description:"viaje cómodo confiable" OR search_text:"viaje cómodo confiable")
```

### Example 3: City-to-City Search

```go
query := &solr.SearchQuery{
    OriginCity:      "Buenos Aires",
    DestinationCity: "Rosario",
    Page:            1,
    Limit:           20,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
fq=origin_city:"Buenos Aires"
fq=destination_city:"Rosario"
```

### Example 4: Price Range Filter

```go
query := &solr.SearchQuery{
    MinPrice: 1000.0,
    MaxPrice: 5000.0,
    Page:     1,
    Limit:    10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
fq=price_per_seat:[1000.0 TO 5000.0]
```

### Example 5: Available Seats Filter

```go
query := &solr.SearchQuery{
    MinSeats: 2,
    Page:     1,
    Limit:    10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
fq=available_seats:[2 TO *]
```

### Example 6: Preferences Filter

```go
petsAllowed := true
query := &solr.SearchQuery{
    PetsAllowed: &petsAllowed,
    Page:        1,
    Limit:       10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
fq=pets_allowed:true
```

### Example 7: Date Range Filter

```go
query := &solr.SearchQuery{
    DepartureFrom: time.Now(),
    DepartureTo:   time.Now().Add(7 * 24 * time.Hour),
    Page:          1,
    Limit:         10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
fq=departure_datetime:[2024-01-15T00:00:00Z TO 2024-01-22T00:00:00Z]
```

### Example 8: Sorting

```go
query := &solr.SearchQuery{
    SortBy:    "price_per_seat",
    SortOrder: "asc",
    Page:      1,
    Limit:     10,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=*:*
sort=price_per_seat asc
```

### Example 9: Complex Query (Multiple Filters)

```go
petsAllowed := false
query := &solr.SearchQuery{
    Query:           "viaje seguro",
    OriginCity:      "Buenos Aires",
    DestinationCity: "Córdoba",
    MinSeats:        2,
    MaxPrice:        6000.0,
    PetsAllowed:     &petsAllowed,
    Status:          "published",
    SortBy:          "departure_datetime",
    SortOrder:       "asc",
    Page:            1,
    Limit:           20,
}

results, err := solrClient.Search(query)
```

Solr query:
```
q=(description:"viaje seguro" OR search_text:"viaje seguro")
fq=origin_city:"Buenos Aires"
fq=destination_city:"Córdoba"
fq=available_seats:[2 TO *]
fq=price_per_seat:[* TO 6000.0]
fq=pets_allowed:false
fq=status:"published"
sort=departure_datetime asc
```

### Example 10: Faceted Search

```go
query := &solr.SearchQuery{
    EnableFacets: true,
    FacetFields:  []string{"origin_city", "destination_city", "status"},
    Page:         1,
    Limit:        10,
}

results, err := solrClient.Search(query)
// results.Facets contains aggregation data
```

## Testing Solr Directly

### Test Core Status

```bash
curl "http://localhost:8983/solr/admin/cores?action=STATUS&core=carpooling_trips"
```

### Test Query (all documents)

```bash
curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=10"
```

### Test Full-Text Search

```bash
curl "http://localhost:8983/solr/carpooling_trips/select?q=description:cómodo&rows=10"
```

### Test Filter Query

```bash
curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&fq=origin_city:\"Buenos Aires\"&rows=10"
```

### Test Facets

```bash
curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0&facet=true&facet.field=origin_city&facet.field=status"
```

## Graceful Degradation

If Solr is unavailable, the search-api will automatically fall back to MongoDB:

```
Health Check Response (Degraded Mode):
{
  "status": "degraded",
  "service": "search-api",
  "services": {
    "mongodb": { "status": "healthy" },
    "solr": { "status": "unhealthy", "message": "Connection refused" },
    "memcached": { "status": "healthy" }
  }
}
```

The API continues to work using MongoDB for search, but with limited functionality:
- ✅ Basic city-to-city search works
- ✅ Geospatial search works
- ❌ Full-text search on descriptions disabled
- ❌ Faceted search disabled
- ❌ Complex multi-field sorting disabled

## Troubleshooting

### Solr Not Accessible

```bash
docker ps | grep solr
docker logs solr
```

### Core Not Created

```bash
# Delete and recreate
curl "http://localhost:8983/solr/admin/cores?action=UNLOAD&core=carpooling_trips&deleteIndex=true"
./scripts/setup_solr_schema.sh
```

### Schema Issues

```bash
# Check schema
curl "http://localhost:8983/solr/carpooling_trips/schema"
```

### No Results in Solr

Check if documents are indexed:
```bash
curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0"
```

If `numFound` is 0, no documents are indexed yet. You need to:
1. Create trips via trips-api
2. Process events via RabbitMQ consumer (not implemented yet)
3. Or manually index trips using the Solr client

## Next Steps

1. ✅ Solr integration complete
2. ⏳ Implement RabbitMQ consumer to index trips automatically
3. ⏳ Implement search service with caching
4. ⏳ Add search endpoints (GET /api/v1/search/trips)
5. ⏳ Add autocomplete endpoint
6. ⏳ Add popular routes endpoint

## References

- [Apache Solr Documentation](https://solr.apache.org/guide/)
- [Go-Solr Library](https://github.com/rtt/Go-Solr)
- [Solr Query Syntax](https://solr.apache.org/guide/solr/latest/query-guide/standard-query-parser.html)
