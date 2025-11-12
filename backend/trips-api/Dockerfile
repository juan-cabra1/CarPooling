# ============================================================================
# STAGE 1: BUILD STAGE
# ============================================================================
# Purpose: Compile the Go application with all optimizations
# Base image: golang:1.23-alpine for minimal build environment
FROM golang:1.24-alpine AS builder

# Metadata labels
LABEL maintainer="CarPooling Team"
LABEL description="Trips API for CarPooling system"
LABEL version="1.0"

# Set working directory for the build
WORKDIR /app

# Install build dependencies
# - git: required for some go modules
# - ca-certificates: for secure HTTPS connections during build
RUN apk add --no-cache git ca-certificates


# Copy dependency files first for optimal Docker layer caching
# This allows Docker to cache the dependency download layer
# Dependencies only re-download when go.mod or go.sum changes
COPY go.mod go.sum ./

# Download all dependencies
# This creates a separate cached layer for dependencies
RUN go mod download && go mod verify

# Copy the entire source code structure
# This includes:
# - cmd/: Application entry points
# - internal/: Internal packages (config, controller, database, middleware, repository, routes, service, models)
COPY . .

# Build the application binary with production optimizations
# Build flags explained:
# - CGO_ENABLED=0: Disable CGO for a fully static binary
# - GOOS=linux: Target Linux OS (for Alpine runtime)
# - -a: Force rebuilding of packages
# - -installsuffix cgo: Add suffix to package directory to avoid conflicts
# - -ldflags="-w -s": Strip debug information and symbol table (reduces binary size)
#   * -w: Omit DWARF symbol table
#   * -s: Omit symbol table and debug information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s" \
    -o trips-api \
    ./cmd/api/main.go

# Verify the binary was created successfully
# RUN ls -lh trips-api && file trips-api

# ============================================================================
# STAGE 2: RUNTIME STAGE
# ============================================================================
# Purpose: Create minimal runtime environment with only necessary components
# Base image: alpine:latest for smallest possible image size
FROM alpine:latest

# Install runtime dependencies only
# - ca-certificates: Required for HTTPS connections to MongoDB, RabbitMQ, etc.
# - tzdata: Timezone data for proper time handling
# Clean up apk cache in the same layer to minimize image size
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create a non-root user for security best practices
# UID 1000 is standard for application users
# This follows the principle of least privilege
RUN adduser -D -u 1000 -g '' appuser

# Set working directory for the application
WORKDIR /app

# Copy only the compiled binary from the builder stage
# This keeps the final image minimal (no source code, no build tools)
COPY --from=builder /app/trips-api .

# Change ownership of the binary to the non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user for runtime security
# Application runs as 'appuser' instead of root
USER appuser

# Expose the application port
# Default port for trips-api service
EXPOSE 8002

# Add health check for container orchestration
# Checks the /health endpoint every 30 seconds
# - interval: Time between health checks
# - timeout: Maximum time to wait for response
# - retries: Number of consecutive failures before marking unhealthy
# - start-period: Grace period before starting health checks
HEALTHCHECK --interval=30s --timeout=3s --retries=3 --start-period=40s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8002/health || exit 1

# ============================================================================
# ENVIRONMENT VARIABLES
# ============================================================================
# The following environment variables should be provided at runtime:
#
# Required:
# - MONGO_URI: MongoDB connection string (e.g., mongodb://mongo:27017/carpooling)
# - JWT_SECRET: Secret key for JWT token signing/verification
# - RABBITMQ_URL: RabbitMQ connection URL (e.g., amqp://guest:guest@rabbitmq:5672/)
#
# Optional:
# - PORT: Server port (default: 8002)
# - GIN_MODE: Gin framework mode (release/debug, default: release)
# - LOG_LEVEL: Logging level (debug/info/warn/error)
#
# Example docker run command:
# docker run -d \
#   -e MONGO_URI=mongodb://mongo:27017/carpooling \
#   -e JWT_SECRET=your-secret-key \
#   -e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
#   -e GIN_MODE=release \
#   -p 8002:8002 \
#   trips-api:latest
# ============================================================================

# Set default environment variables
ENV GIN_MODE=release \
    PORT=8002

# Run the compiled binary
# Using exec form (JSON array) to ensure proper signal handling
CMD ["./trips-api"]
