package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookings-api/internal/clients"
	"bookings-api/internal/config"
	"bookings-api/internal/controller"
	"bookings-api/internal/database"
	"bookings-api/internal/repository"
	"bookings-api/internal/routes"
	"bookings-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// ============================================================================
	// LOGGING SETUP
	// ============================================================================
	// Configure zerolog for structured logging
	// - Uses console writer for human-readable output in development
	// - Uses JSON output in production for log aggregation
	// - Includes timestamps for all log entries
	// - Supports configurable log level via LOG_LEVEL env variable
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level from environment variable (default: INFO)
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Configure output format based on environment
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		// Production: JSON output for log aggregation tools
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Development: Console writer for human-readable logs
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	}

	log.Info().
		Str("log_level", zerolog.GlobalLevel().String()).
		Str("environment", environment).
		Msg("🚀 Starting Bookings API service...")

	// ============================================================================
	// CONFIGURATION LOADING
	// ============================================================================
	// Load configuration from environment variables
	// This must be done FIRST before any other initialization
	// Configuration includes:
	//   - Server port (default: 8003)
	//   - Database connection URL (MySQL DSN)
	//   - JWT secret for authentication
	//   - RabbitMQ URL for message queue
	//   - Environment (development/production)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Failed to load configuration")
	}
	log.Info().
		Str("port", cfg.ServerPort).
		Str("environment", cfg.Environment).
		Msg("✅ Configuration loaded successfully")

	// ============================================================================
	// DATABASE INITIALIZATION
	// ============================================================================
	// Initialize MySQL database connection using GORM
	// This must be done BEFORE initializing repositories and services
	// Database is used for:
	//   - Storing booking records (bookings table)
	//   - Event idempotency tracking (processed_events table)
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Failed to initialize database connection")
	}
	log.Info().Msg("✅ Database connection established")

	// Run auto-migrations to create/update database tables
	// This creates tables if they don't exist and adds new columns
	// Safe to run on every startup (won't delete existing data)
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Failed to run database migrations")
	}
	log.Info().Msg("✅ Database migrations completed")

	// ============================================================================
	// REPOSITORY INITIALIZATION (Data Layer)
	// ============================================================================
	// Create repository instances for data access
	// Repositories abstract database operations and provide a clean interface
	bookingRepo := repository.NewBookingRepository(db)
	eventRepo := repository.NewEventRepository(db)
	log.Info().Msg("✅ Repositories initialized")

	// ============================================================================
	// HTTP CLIENT INITIALIZATION (External Dependencies)
	// ============================================================================
	// Create HTTP clients for calling other microservices
	tripsClient := clients.NewTripsClient(cfg.TripsAPIURL)
	log.Info().
		Str("trips_api_url", cfg.TripsAPIURL).
		Msg("✅ HTTP clients initialized")

	// ============================================================================
	// SERVICE INITIALIZATION (Business Logic Layer)
	// ============================================================================
	// Create service instances with dependency injection
	// Services contain business logic and orchestrate repositories and clients

	// AuthService: JWT token validation for authentication middleware
	authService := service.NewAuthService(cfg.JWTSecret)

	// IdempotencyService: Used by RabbitMQ consumer (Issue #5) to prevent duplicate event processing
	_ = service.NewIdempotencyService(eventRepo)

	// BookingService: Will be used by BookingController (Issue #6) for HTTP endpoints
	_ = service.NewBookingService(bookingRepo, tripsClient)

	log.Info().Msg("✅ Services initialized (ready for controllers and consumers)")

	// ============================================================================
	// GIN ROUTER INITIALIZATION
	// ============================================================================
	// Initialize Gin HTTP framework
	// Using gin.New() with custom middleware for better control:
	//   - Custom zerolog-based logging middleware (consistent with app logging)
	//   - Recovery middleware (recovers from panics)
	//   - Error handling middleware (standardized error responses)
	//
	// Note: We use gin.New() instead of gin.Default() to avoid duplicate logging
	// and to integrate with our zerolog setup

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Msg("🔒 Running in PRODUCTION mode")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Info().Msg("🔧 Running in DEVELOPMENT mode")
	}

	router := gin.New()

	// Add built-in middleware
	router.Use(gin.Recovery()) // Recover from panics

	// ============================================================================
	// CONTROLLER INITIALIZATION
	// ============================================================================
	// Create controller instances
	// Controllers handle HTTP requests and responses
	// Each controller is responsible for a specific domain (health, bookings, etc.)
	healthController := controller.NewHealthController("bookings-api", cfg.ServerPort)
	log.Info().Msg("✅ Controllers initialized")

	// ============================================================================
	// ROUTE REGISTRATION
	// ============================================================================
	// Register all HTTP routes with the router
	// Routes define the API endpoints and map them to controller methods
	// This includes:
	//   - Health check endpoint (GET /health)
	//   - Booking management endpoints (protected by JWT authentication)
	routes.SetupRoutes(router, healthController, authService)
	log.Info().Msg("✅ Routes registered")

	// ============================================================================
	// HTTP SERVER CONFIGURATION
	// ============================================================================
	// Configure HTTP server with production-ready timeouts
	// Timeouts prevent slow clients from holding connections indefinitely
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,

		// ReadTimeout: Maximum duration for reading the entire request (headers + body)
		// Prevents slow-read attacks
		ReadTimeout: 10 * time.Second,

		// WriteTimeout: Maximum duration for writing the response
		// Should be longer than ReadTimeout to account for processing time
		WriteTimeout: 30 * time.Second,

		// IdleTimeout: Maximum time to wait for the next request when keep-alives are enabled
		// Prevents idle connections from consuming resources
		IdleTimeout: 60 * time.Second,

		// ReadHeaderTimeout: Maximum time to read request headers
		// First line of defense against slow-read attacks
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ============================================================================
	// SERVER STARTUP (in goroutine for graceful shutdown)
	// ============================================================================
	// Start HTTP server in a separate goroutine
	// This allows us to listen for shutdown signals in the main goroutine
	go func() {
		log.Info().
			Str("port", cfg.ServerPort).
			Str("health_endpoint", "http://localhost:"+cfg.ServerPort+"/health").
			Msg("🌐 HTTP server listening")

		// ListenAndServe blocks until the server is shut down
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Msg("❌ Failed to start HTTP server")
		}
	}()

	// ============================================================================
	// GRACEFUL SHUTDOWN HANDLING
	// ============================================================================
	// Set up signal handling for graceful shutdown
	// This ensures the server shuts down cleanly when receiving:
	//   - SIGINT (Ctrl+C)
	//   - SIGTERM (Docker/Kubernetes shutdown signal)
	//
	// Graceful shutdown process:
	//   1. Stop accepting new requests
	//   2. Wait for in-flight requests to complete (up to 15 seconds)
	//   3. Close all connections
	//   4. Exit cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a shutdown signal
	sig := <-quit
	log.Info().
		Str("signal", sig.String()).
		Msg("⚠️  Shutdown signal received, starting graceful shutdown...")

	// Create a context with timeout for graceful shutdown
	// If shutdown takes longer than 15 seconds, force exit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	// This will wait for active connections to finish (up to 15 seconds)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Server forced to shutdown")
	}

	// Close database connection
	// Release all database connections in the pool
	// This prevents connection leaks and ensures clean shutdown
	if err := database.CloseDB(db); err != nil {
		log.Error().
			Err(err).
			Msg("⚠️  Error closing database connection")
	}

	log.Info().Msg("✅ Server gracefully stopped")
}
