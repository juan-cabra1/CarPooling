package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"search-api/internal/config"
	"search-api/internal/database"
	"search-api/internal/repository"
	"search-api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	// Configure zerolog for structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Println("🚀 Starting search-api...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Printf("✅ Configuration loaded successfully")

	// Connect to MongoDB
	db, err := database.ConnectMongoDB(cfg.Mongo.URI, cfg.Mongo.DB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	// Create indexes - CRITICAL for performance and idempotency
	if err := database.CreateIndexes(db); err != nil {
		log.Fatalf("❌ Failed to create MongoDB indexes: %v", err)
	}

	// Initialize repositories
	_ = repository.NewSearchRepository(db)
	_ = repository.NewEventRepository(db)
	_ = repository.NewPopularRouteRepository(db)
	log.Println("✅ Repositories initialized successfully")

	// Future: Initialize services and controllers here
	// Example:
	// searchService := service.NewSearchService(searchRepo, popularRouteRepo)
	// searchController := controller.NewSearchController(searchService)

	// Setup Gin router
	router := gin.Default()
	routes.SetupRoutes(router)
	log.Println("✅ Routes configured successfully")

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
		log.Printf("✅ search-api server listening on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down search-api server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ search-api server exited gracefully")
}
