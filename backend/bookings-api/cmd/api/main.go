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
	"bookings-api/internal/messaging"
	"bookings-api/internal/publisher"
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
	// RABBITMQ PUBLISHER INITIALIZATION
	// ============================================================================
	// Initialize RabbitMQ publisher for reservation events
	// Publisher emits events to trips-api when:
	//   - reservation.created: New booking created (trips-api decrements seats)
	//   - reservation.cancelled: Booking cancelled (trips-api increments seats)
	//
	// Publisher features:
	//   - Persistent messages (survive RabbitMQ restart)
	//   - Topic exchange: "bookings.events"
	//   - Structured logging with zerolog
	//   - Graceful error handling (no panics)
	reservationPublisher, err := publisher.NewReservationPublisher(cfg, log.Logger)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Failed to initialize RabbitMQ publisher")
	}
	log.Info().
		Str("exchange", "bookings.events").
		Msg("✅ RabbitMQ publisher initialized")

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

	// IdempotencyService: Used by RabbitMQ consumer to prevent duplicate event processing
	idempotencyService := service.NewIdempotencyService(eventRepo)

	// BookingService: Handles business logic for booking operations
	// Injected dependencies: repository, trips-api client, RabbitMQ publisher
	bookingService := service.NewBookingService(bookingRepo, tripsClient, reservationPublisher)

	log.Info().Msg("✅ Services initialized (ready for controllers and consumers)")

	// ============================================================================
	// RABBITMQ CONSUMER INITIALIZATION
	// ============================================================================
	// Initialize RabbitMQ consumer for trip events
	// Consumer handles:
	//   - trip.cancelled: Cancels all confirmed bookings for cancelled trips
	//   - reservation.failed: Marks bookings as failed when seat reservation fails
	//
	// Consumer features:
	//   - Idempotency: Prevents duplicate event processing using event_id
	//   - Manual ACK: Only acknowledges after successful processing
	//   - Prefetch: Processes 10 messages concurrently for better throughput
	//   - Graceful shutdown: Stops cleanly on SIGINT/SIGTERM
	consumer, err := messaging.NewTripsConsumer(
		cfg.RabbitMQURL,
		bookingRepo,
		idempotencyService,
	)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Failed to initialize RabbitMQ consumer")
	}
	log.Info().
		Str("rabbitmq_url", cfg.RabbitMQURL).
		Msg("✅ RabbitMQ consumer initialized")

	// Start consumer in background goroutine
	// This allows the consumer to process messages concurrently with HTTP requests
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	go func() {
		log.Info().Msg("🐰 Starting RabbitMQ consumer...")
		if err := consumer.Start(consumerCtx); err != nil {
			log.Error().
				Err(err).
				Msg("❌ RabbitMQ consumer stopped with error")
		}
	}()

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
	bookingController := controller.NewBookingController(bookingService)
	log.Info().Msg("✅ Controllers initialized")

	// ============================================================================
	// ROUTE REGISTRATION
	// ============================================================================
	// Register all HTTP routes with the router
	// Routes define the API endpoints and map them to controller methods
	// This includes:
	//   - Health check endpoint (GET /health)
	//   - Booking management endpoints (protected by JWT authentication)
	routes.SetupRoutes(router, healthController, bookingController, authService)
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

	// Attempt graceful shutdown of HTTP server
	// This will wait for active connections to finish (up to 15 seconds)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().
			Err(err).
			Msg("❌ Server forced to shutdown")
	}
	log.Info().Msg("✅ HTTP server stopped")

	// Stop RabbitMQ consumer
	// Cancel consumer context to stop processing new messages
	// Wait for in-flight messages to complete
	log.Info().Msg("⏳ Stopping RabbitMQ consumer...")
	consumerCancel()

	// Give consumer time to finish processing current messages
	time.Sleep(2 * time.Second)

	// Close RabbitMQ consumer connection
	if err := consumer.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("⚠️  Error closing RabbitMQ consumer connection")
	} else {
		log.Info().Msg("✅ RabbitMQ consumer stopped")
	}

	// Close RabbitMQ publisher connection
	// This closes the publishing channel and connection gracefully
	log.Info().Msg("⏳ Stopping RabbitMQ publisher...")
	if err := reservationPublisher.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("⚠️  Error closing RabbitMQ publisher connection")
	} else {
		log.Info().Msg("✅ RabbitMQ publisher stopped")
	}

	// Close database connection
	// Release all database connections in the pool
	// This prevents connection leaks and ensures clean shutdown
	if err := database.CloseDB(db); err != nil {
		log.Error().
			Err(err).
			Msg("⚠️  Error closing database connection")
	} else {
		log.Info().Msg("✅ Database connection closed")
	}

	log.Info().Msg("✅ Server gracefully stopped")
}
