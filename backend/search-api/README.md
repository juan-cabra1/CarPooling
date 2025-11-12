# Search API Service

## Overview

The **Search API** is a consumer-only microservice within the CarPooling platform that provides advanced search capabilities for trips. It operates as an event-driven service that listens to events from the trips-api, denormalizes data by fetching driver information from users-api, and provides fast search operations through multiple data stores optimized for different query patterns.

### Key Features

- **Event-Driven Architecture**: Consumes events from trips-api via RabbitMQ
- **Data Denormalization**: Automatically fetches and stores driver information for fast reads
- **Full-Text Search**: Apache Solr integration for complex text-based queries
- **Geospatial Search**: MongoDB geospatial indexes for location-based queries
- **High-Performance Caching**: Memcached layer for frequently accessed data
- **Read-Optimized**: Designed specifically for fast search and retrieval operations

### Service Purpose

The search-api acts as a specialized read layer that:
1. Listens to trip events (created, updated, deleted) from trips-api
2. Denormalizes trip data with driver information from users-api
3. Indexes data in Solr for full-text search capabilities
4. Stores data in MongoDB with geospatial indexes for location queries
5. Caches frequently accessed results in Memcached
6. Exposes REST endpoints for various search operations

## Technology Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.21+ | Primary programming language |
| **Gin** | Latest | HTTP web framework |
| **MongoDB** | 7.0+ | Primary data store with geospatial support |
| **Apache Solr** | 9.0+ | Full-text search engine |
| **Memcached** | Latest | In-memory caching layer |
| **RabbitMQ** | 3.12+ | Message broker for event consumption |
| **JWT** | v5 | Authentication token validation |
| **Zerolog** | Latest | Structured logging |

## Architecture

### System Architecture

```
┌─────────────┐         ┌─────────────┐
│  Trips API  │────────>│  RabbitMQ   │
└─────────────┘         └──────┬──────┘
                              │ Events
                              ▼
┌─────────────────────────────────────────────┐
│            Search API (This Service)        │
│  ┌──────────────┐  ┌──────────────┐        │
│  │   Consumer   │  │  HTTP API    │        │
│  └──────┬───────┘  └──────▲───────┘        │
│         │                  │                │
│         ▼                  │                │
│  ┌──────────────┐         │                │
│  │   Service    │─────────┘                │
│  │    Layer     │                          │
│  └──────┬───────┘                          │
│         │                                   │
│    ┌────┴────┬────────┬─────────┐         │
│    ▼         ▼        ▼         ▼         │
│ ┌──────┐ ┌──────┐ ┌───────┐ ┌────────┐   │
│ │ Solr │ │ Mongo│ │ Cache │ │External│   │
│ └──────┘ └──────┘ └───────┘ │  APIs  │   │
│                              └────────┘   │
└─────────────────────────────────────────────┘
```

### Clean Architecture Layers

```
cmd/api/
  └── main.go                    # Application entry point

internal/
  ├── config/                    # Configuration management
  ├── domain/                    # Domain entities and models
  ├── dao/                       # Database Access Objects
  ├── repository/                # Data access layer (MongoDB)
  ├── search/                    # Search engine logic (Solr)
  ├── cache/                     # Caching layer (Memcached)
  ├── service/                   # Business logic layer
  ├── consumer/                  # RabbitMQ event consumers
  ├── controller/                # HTTP handlers
  ├── middleware/                # HTTP middleware (JWT, CORS, etc.)
  └── routes/                    # Route definitions
```

## Prerequisites

Before running the search-api service, ensure you have the following installed and running:

### Required Services

1. **Go** (version 1.21 or higher)
   ```bash
   go version
   ```

2. **MongoDB** (version 7.0 or higher)
   ```bash
   # Local installation
   mongod --version

   # Or via Docker
   docker run -d -p 27017:27017 --name mongodb mongo:7.0
   ```

3. **Apache Solr** (version 9.0 or higher)
   ```bash
   # Via Docker (recommended)
   docker run -d -p 8983:8983 --name solr solr:9.0

   # Create the trips core
   docker exec solr bin/solr create_core -c trips
   ```

4. **Memcached**
   ```bash
   # Local installation
   memcached -V

   # Or via Docker
   docker run -d -p 11211:11211 --name memcached memcached:latest
   ```

5. **RabbitMQ** (version 3.12 or higher)
   ```bash
   # Via Docker (recommended)
   docker run -d -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3.12-management
   ```

### External Dependencies

- **users-api**: Must be running on port 8001 (for driver data)
- **trips-api**: Must be running on port 8002 (for trip events)

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd backend/search-api
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment Variables

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your actual configuration
nano .env  # or use your preferred editor
```

**Required Environment Variables:**

```env
# Server
SERVER_PORT=8004

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=search_db

# Apache Solr
SOLR_URL=http://localhost:8983/solr
SOLR_CORE=trips

# Memcached
MEMCACHED_SERVERS=localhost:11211
CACHE_TTL=300

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=search.events

# External APIs
TRIPS_API_URL=http://localhost:8002
USERS_API_URL=http://localhost:8001

# JWT
JWT_SECRET=your-secret-key-here

# Environment
ENVIRONMENT=development
```

### 4. Setup Apache Solr

Create the required Solr core for trip indexing:

```bash
# If running locally
curl "http://localhost:8983/solr/admin/cores?action=CREATE&name=trips&configSet=_default"

# If using Docker
docker exec solr bin/solr create_core -c trips
```

### 5. Setup RabbitMQ Queue

The queue will be automatically created by the consumer, but you can create it manually:

```bash
# Access RabbitMQ management UI
http://localhost:15672
# Default credentials: guest/guest

# Create queue named: search.events
# Ensure trips-api is configured to publish to this queue
```

## Running the Service

### Development Mode

```bash
# Navigate to the service directory
cd backend/search-api

# Run with automatic reload (if using air)
air

# Or run directly
go run cmd/api/main.go
```

### Production Mode

```bash
# Build the binary
go build -o search-api cmd/api/main.go

# Run the binary
./search-api
```

### Using Docker

```bash
# Build the Docker image
docker build -t search-api:latest .

# Run the container
docker run -p 8004:8004 --env-file .env search-api:latest
```

### Using Docker Compose

```bash
# Start all services including dependencies
docker-compose up -d

# View logs
docker-compose logs -f search-api

# Stop all services
docker-compose down
```

## API Endpoints

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "search-api",
  "version": "1.0.0",
  "timestamp": "2025-11-12T10:00:00Z"
}
```

### Search Endpoints (Planned)

#### Search Trips by Text

```http
GET /api/v1/search/trips?q={query}&limit={limit}&offset={offset}
Authorization: Bearer {jwt_token}
```

**Query Parameters:**
- `q` (required): Search query string
- `limit` (optional): Number of results (default: 20)
- `offset` (optional): Pagination offset (default: 0)

#### Search Trips by Location

```http
GET /api/v1/search/nearby?lat={latitude}&lng={longitude}&radius={radius}
Authorization: Bearer {jwt_token}
```

**Query Parameters:**
- `lat` (required): Latitude coordinate
- `lng` (required): Longitude coordinate
- `radius` (optional): Search radius in meters (default: 5000)

#### Advanced Search

```http
POST /api/v1/search/advanced
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "origin": "Buenos Aires",
  "destination": "Córdoba",
  "departure_date": "2025-11-15",
  "seats_available": 2,
  "max_price": 5000
}
```

## Event Consumption

The service listens to the following events from trips-api:

### Event Types

| Event Type | Description | Action |
|------------|-------------|--------|
| `trip.created` | New trip created | Index in Solr, store in MongoDB, cache result |
| `trip.updated` | Trip details updated | Update Solr index, update MongoDB, invalidate cache |
| `trip.deleted` | Trip removed | Remove from Solr, delete from MongoDB, clear cache |
| `trip.status_changed` | Trip status changed | Update indexes and cache |

### Event Payload Example

```json
{
  "event_type": "trip.created",
  "timestamp": "2025-11-12T10:00:00Z",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "driver_id": "789e4567-e89b-12d3-a456-426614174001",
    "origin": "Buenos Aires",
    "destination": "Córdoba",
    "departure_time": "2025-11-15T08:00:00Z",
    "seats_available": 3,
    "price": 4500
  }
}
```

## Development Guidelines

### Code Structure

Follow clean architecture principles:

1. **Domain Layer**: Pure business entities with no external dependencies
2. **Data Access Layer**: Repositories and DAOs for data operations
3. **Service Layer**: Business logic and orchestration
4. **Presentation Layer**: HTTP handlers and middleware

### Adding New Features

1. Define domain models in `internal/domain/`
2. Create repository interfaces and implementations
3. Implement business logic in `internal/service/`
4. Add HTTP handlers in `internal/controller/`
5. Register routes in `internal/routes/`

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...
```

### Logging

Use zerolog for structured logging:

```go
log.Info().
    Str("trip_id", tripID).
    Str("event_type", "trip.created").
    Msg("Processing trip event")
```

### Error Handling

Always use proper error wrapping:

```go
if err != nil {
    return fmt.Errorf("failed to index trip: %w", err)
}
```

## Performance Optimization

### Caching Strategy

1. **Hot Data**: Cache search results for 5 minutes (CACHE_TTL=300)
2. **Warm Data**: Cache driver information for 15 minutes
3. **Cold Data**: Fetch from MongoDB/Solr as needed

### MongoDB Indexes

Ensure the following indexes are created:

```javascript
// Geospatial index for location queries
db.trips.createIndex({ "origin_location": "2dsphere" })
db.trips.createIndex({ "destination_location": "2dsphere" })

// Compound indexes for common queries
db.trips.createIndex({ "departure_time": 1, "seats_available": 1 })
db.trips.createIndex({ "driver_id": 1, "status": 1 })
```

### Solr Schema Optimization

Configure appropriate field types in Solr schema for optimal text search.

## Monitoring and Observability

### Metrics to Monitor

- RabbitMQ consumer lag
- Search query latency
- Cache hit/miss ratio
- MongoDB query performance
- Solr indexing speed

### Health Checks

- MongoDB connection status
- Solr availability
- Memcached connectivity
- RabbitMQ connection status

## Troubleshooting

### Common Issues

**Issue**: Service fails to start

- Check that all required environment variables are set
- Verify MongoDB, Solr, Memcached, and RabbitMQ are running
- Check logs for specific error messages

**Issue**: Events not being consumed

- Verify RabbitMQ connection string is correct
- Check that the queue exists and has messages
- Ensure trips-api is publishing to the correct queue

**Issue**: Search queries return empty results

- Verify Solr core is created and accessible
- Check that events are being processed and indexed
- Review Solr logs for indexing errors

**Issue**: High memory usage

- Adjust Memcached memory limits
- Review cache TTL settings
- Monitor MongoDB connection pooling

## Security Considerations

1. **JWT Validation**: All search endpoints require valid JWT tokens
2. **Rate Limiting**: Implement rate limiting for search queries
3. **Input Validation**: Sanitize all search query inputs
4. **Secure Connections**: Use TLS for production deployments
5. **Secrets Management**: Use secret management tools for sensitive data

## Contributing

1. Follow Go coding standards and best practices
2. Write comprehensive tests for new features
3. Update documentation for API changes
4. Use meaningful commit messages
5. Create pull requests for review

## License

[Specify your license here]

## Contact

For questions or support, contact the development team.

---

**Last Updated**: 2025-11-12
**Version**: 1.0.0
**Service Port**: 8004
