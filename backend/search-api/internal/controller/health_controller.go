package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	solr "github.com/rtt/Go-Solr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"search-api/internal/config"
)

// HealthController handles health check endpoints
type HealthController struct {
	mongoClient      *mongo.Client
	solrClient       *solr.Connection
	memcachedClient  *memcache.Client
	logger           zerolog.Logger
	config           *config.Config
}

// ServiceHealthStatus represents the health status of a single service
type ServiceHealthStatus struct {
	Status  string `json:"status"`  // "healthy" or "unhealthy"
	Message string `json:"message"` // Connection status or error message
}

// HealthCheckResponse represents the complete health check response
type HealthCheckResponse struct {
	Status   string                         `json:"status"`  // "ok" or "degraded"
	Service  string                         `json:"service"` // "search-api"
	Port     string                         `json:"port"`    // "8004"
	Services map[string]ServiceHealthStatus `json:"services"`
}

// NewHealthController creates a new health controller instance
func NewHealthController(
	mongoClient *mongo.Client,
	solrClient *solr.Connection,
	memcachedClient *memcache.Client,
	logger zerolog.Logger,
	cfg *config.Config,
) *HealthController {
	return &HealthController{
		mongoClient:     mongoClient,
		solrClient:      solrClient,
		memcachedClient: memcachedClient,
		logger:          logger,
		config:          cfg,
	}
}

// HealthCheck handles the GET /health endpoint
// Returns the health status of the search-api service and all its dependencies
func (hc *HealthController) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response := HealthCheckResponse{
		Service:  "search-api",
		Port:     hc.config.ServerPort,
		Services: make(map[string]ServiceHealthStatus),
	}

	// Track overall health status
	allHealthy := true

	// Check MongoDB health
	mongoStatus := hc.checkMongoHealth(ctx)
	response.Services["mongodb"] = mongoStatus
	if mongoStatus.Status == "unhealthy" {
		allHealthy = false
		hc.logger.Warn().Msg("MongoDB health check failed")
	}

	// Check Apache Solr health
	solrStatus := hc.checkSolrHealth(ctx)
	response.Services["solr"] = solrStatus
	if solrStatus.Status == "unhealthy" {
		allHealthy = false
		hc.logger.Warn().Msg("Apache Solr health check failed")
	}

	// Check Memcached health
	memcachedStatus := hc.checkMemcachedHealth(ctx)
	response.Services["memcached"] = memcachedStatus
	if memcachedStatus.Status == "unhealthy" {
		allHealthy = false
		hc.logger.Warn().Msg("Memcached health check failed")
	}

	// Set overall status
	if allHealthy {
		response.Status = "ok"
		hc.logger.Debug().Msg("All services healthy")
	} else {
		response.Status = "degraded"
		hc.logger.Warn().Msg("Service running in degraded mode")
	}

	// Always return 200 OK, even in degraded mode
	// This allows the service to continue operating with reduced functionality
	c.JSON(http.StatusOK, response)
}

// checkMongoHealth tests MongoDB connectivity
func (hc *HealthController) checkMongoHealth(ctx context.Context) ServiceHealthStatus {
	if hc.mongoClient == nil {
		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "MongoDB client not initialized",
		}
	}

	// Ping MongoDB with timeout
	err := hc.mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		hc.logger.Error().
			Err(err).
			Msg("MongoDB ping failed")

		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "Failed to ping MongoDB: " + err.Error(),
		}
	}

	return ServiceHealthStatus{
		Status:  "healthy",
		Message: "Connected",
	}
}

// checkSolrHealth tests Apache Solr connectivity
func (hc *HealthController) checkSolrHealth(ctx context.Context) ServiceHealthStatus {
	if hc.solrClient == nil {
		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "Solr client not initialized",
		}
	}

	// Try a simple ping query to Solr
	// Since Go-Solr doesn't have a dedicated ping method, we perform a simple query
	query := solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{"*:*"},
			"rows": []string{"0"},
		},
	}

	// Execute query with timeout handling
	done := make(chan error, 1)
	go func() {
		_, err := hc.solrClient.Select(&query)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			hc.logger.Error().
				Err(err).
				Msg("Solr query failed")

			return ServiceHealthStatus{
				Status:  "unhealthy",
				Message: "Failed to connect to Solr: " + err.Error(),
			}
		}

		return ServiceHealthStatus{
			Status:  "healthy",
			Message: "Connected",
		}

	case <-ctx.Done():
		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "Solr health check timeout",
		}
	}
}

// checkMemcachedHealth tests Memcached connectivity
func (hc *HealthController) checkMemcachedHealth(ctx context.Context) ServiceHealthStatus {
	if hc.memcachedClient == nil {
		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "Memcached client not initialized",
		}
	}

	// Try to get a test key (it's OK if it doesn't exist)
	// We're just testing connectivity, not the actual value
	testKey := "health_check_ping"

	done := make(chan error, 1)
	go func() {
		_, err := hc.memcachedClient.Get(testKey)
		// ErrCacheMiss is actually a good sign - it means we connected successfully
		if err == memcache.ErrCacheMiss {
			done <- nil
		} else {
			done <- err
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			hc.logger.Error().
				Err(err).
				Msg("Memcached connection failed")

			return ServiceHealthStatus{
				Status:  "unhealthy",
				Message: "Failed to connect to Memcached: " + err.Error(),
			}
		}

		return ServiceHealthStatus{
			Status:  "healthy",
			Message: "Connected",
		}

	case <-ctx.Done():
		return ServiceHealthStatus{
			Status:  "unhealthy",
			Message: "Memcached health check timeout",
		}
	}
}
