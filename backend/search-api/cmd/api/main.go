package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	solr "github.com/rtt/Go-Solr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"search-api/internal/config"
	"search-api/internal/controller"
	"search-api/internal/routes"
)

func main() {
	// ========================================================================
	// LOGGING SETUP
	// ========================================================================
	// Configure zerolog for structured logging with pretty console output
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting search-api service...")

	// ========================================================================
	// CONFIGURATION LOADING
	// ========================================================================
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to load configuration")
	}

	log.Info().
		Str("environment", cfg.Environment).
		Str("port", cfg.ServerPort).
		Msg("Configuration loaded successfully")

	// ========================================================================
	// MONGODB INITIALIZATION
	// ========================================================================
	log.Info().
		Str("uri", cfg.MongoURI).
		Str("database", cfg.MongoDB).
		Msg("Connecting to MongoDB...")

	mongoClient, err := initMongoDB(cfg)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to initialize MongoDB - cannot operate without storage")
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Msg("Error disconnecting from MongoDB")
		}
	}()

	log.Info().Msg("MongoDB connected successfully")

	// ========================================================================
	// APACHE SOLR INITIALIZATION
	// ========================================================================
	log.Info().
		Str("url", cfg.GetSolrFullURL()).
		Msg("Connecting to Apache Solr...")

	solrClient, err := initSolr(cfg)
	if err != nil {
		// Solr is not critical - service can operate in degraded mode
		log.Warn().
			Err(err).
			Msg("Failed to connect to Apache Solr - search features will be disabled (degraded mode)")
		solrClient = nil // Continue without Solr
	} else {
		log.Info().Msg("Apache Solr connected successfully")
	}

	// ========================================================================
	// MEMCACHED INITIALIZATION
	// ========================================================================
	log.Info().
		Strs("servers", cfg.MemcachedServers).
		Msg("Connecting to Memcached...")

	memcachedClient, err := initMemcached(cfg)
	if err != nil {
		// Memcached is not critical - service can operate without caching
		log.Warn().
			Err(err).
			Msg("Failed to connect to Memcached - caching will be disabled (degraded mode)")
		memcachedClient = nil // Continue without cache
	} else {
		log.Info().Msg("Memcached connected successfully")
	}

	// ========================================================================
	// GIN ROUTER INITIALIZATION
	// ========================================================================
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Msg("Running in PRODUCTION mode")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Info().Msg("Running in DEVELOPMENT mode")
	}

	router := gin.Default()

	// Add custom logger middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Str("ip", param.ClientIP).
			Dur("latency", param.Latency).
			Msg("HTTP request")
		return ""
	}))

	// ========================================================================
	// CONTROLLERS SETUP
	// ========================================================================
	healthController := controller.NewHealthController(
		mongoClient,
		solrClient,
		memcachedClient,
		log.Logger,
		cfg,
	)

	// ========================================================================
	// ROUTES SETUP
	// ========================================================================
	routes.SetupRoutes(router, healthController)

	log.Info().Msg("Routes configured successfully")

	// ========================================================================
	// HTTP SERVER CONFIGURATION
	// ========================================================================
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ========================================================================
	// GRACEFUL SHUTDOWN SETUP
	// ========================================================================
	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Info().
			Str("port", cfg.ServerPort).
			Str("address", "http://localhost:"+cfg.ServerPort).
			Msg("Server is starting...")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Msg("Failed to start HTTP server")
		}
	}()

	log.Info().
		Str("health_check", "http://localhost:"+cfg.ServerPort+"/health").
		Msg("search-api is ready to accept requests")

	// ========================================================================
	// WAIT FOR INTERRUPT SIGNAL
	// ========================================================================
	<-quit
	log.Info().Msg("Shutting down server gracefully...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Server forced to shutdown")
	}

	// Close MongoDB connection
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Error disconnecting from MongoDB")
	} else {
		log.Info().Msg("MongoDB connection closed")
	}

	// Close Memcached connection
	if memcachedClient != nil {
		log.Info().Msg("Memcached connection closed")
	}

	log.Info().Msg("search-api service stopped successfully")
}

// ============================================================================
// SERVICE INITIALIZATION FUNCTIONS
// ============================================================================

// initMongoDB initializes and tests the MongoDB connection
func initMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure MongoDB client options
	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	// Create MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	// Test the connection with ping
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// initSolr initializes and tests the Apache Solr connection
func initSolr(cfg *config.Config) (*solr.Connection, error) {
	// Create Solr connection
	// Go-Solr expects: Init(host string, port int, core string)
	solrConn, err := solr.Init(cfg.SolrURL, 8983, cfg.SolrCore)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Solr connection: %w", err)
	}

	// Test connection with a simple query
	query := solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{"*:*"},
			"rows": []string{"0"},
		},
	}

	_, err = solrConn.Select(&query)
	if err != nil {
		return nil, fmt.Errorf("failed to ping Solr: %w", err)
	}

	return solrConn, nil
}

// initMemcached initializes and tests the Memcached connection
func initMemcached(cfg *config.Config) (*memcache.Client, error) {
	// Create Memcached client
	mc := memcache.New(cfg.MemcachedServers...)

	// Configure client timeouts
	mc.Timeout = 100 * time.Millisecond
	mc.MaxIdleConns = 10

	// Test connection with a ping operation
	testKey := "health_check_init"
	testValue := []byte("ping")

	// Try to set a test value
	err := mc.Set(&memcache.Item{
		Key:        testKey,
		Value:      testValue,
		Expiration: 10, // 10 seconds
	})
	if err != nil {
		return nil, fmt.Errorf("failed to test Memcached connection: %w", err)
	}

	// Try to get the test value back
	_, err = mc.Get(testKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test value from Memcached: %w", err)
	}

	// Delete the test key
	_ = mc.Delete(testKey)

	return mc, nil
}
