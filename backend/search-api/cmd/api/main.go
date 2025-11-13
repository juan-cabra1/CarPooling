package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"search-api/internal/config"
	"search-api/internal/controller"
	"search-api/internal/database"
	"search-api/internal/repository"
	"search-api/internal/routes"
	"search-api/internal/solr"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Configure zerolog for structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting search-api")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	log.Info().Msg("Configuration loaded successfully")

	// Connect to MongoDB
	db, err := database.ConnectMongoDB(cfg.Mongo.URI, cfg.Mongo.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// Create indexes - CRITICAL for performance and idempotency
	if err := database.CreateIndexes(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to create MongoDB indexes")
	}
	log.Info().Msg("MongoDB indexes created successfully")

	// Connect to Apache Solr
	solrClient, err := solr.NewClient(cfg.Solr.URL, cfg.Solr.Core)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Solr (will fallback to MongoDB)")
		solrClient = nil // Continue without Solr - graceful degradation
	} else {
		// Test Solr connection
		if err := solrClient.Ping(); err != nil {
			log.Warn().Err(err).Msg("Solr ping failed (will fallback to MongoDB)")
			solrClient = nil
		} else {
			log.Info().Msg("Connected to Apache Solr successfully")
		}
	}

	// Connect to Memcached
	memcachedClient := memcache.New(cfg.Memcached.Servers...)
	if err := memcachedClient.Ping(); err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Memcached (caching disabled)")
		memcachedClient = nil
	} else {
		log.Info().Msg("Connected to Memcached successfully")
	}

	// Get MongoDB client for health checks
	mongoClient := getMongoClient(db)

	// Initialize repositories
	_ = repository.NewSearchRepository(db)
	_ = repository.NewEventRepository(db)
	_ = repository.NewPopularRouteRepository(db)
	log.Info().Msg("Repositories initialized successfully")

	// Initialize controllers
	healthController := controller.NewHealthController(
		mongoClient,
		solrClient,
		memcachedClient,
		cfg,
	)
	log.Info().Msg("Controllers initialized successfully")

	// Setup Gin router
	router := gin.Default()
	routes.SetupRoutes(router, healthController)
	log.Info().Msg("Routes configured successfully")

	// Configure HTTP server with timeouts
	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("search-api server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down search-api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("search-api server exited gracefully")
}

// getMongoClient extracts the *mongo.Client from *mongo.Database
func getMongoClient(db *mongo.Database) *mongo.Client {
	return db.Client()
}
