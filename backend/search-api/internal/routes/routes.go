package routes

import (
	"time"

	"search-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the search-api
func SetupRoutes(router *gin.Engine) {
	// Apply global middlewares
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Search endpoints will be added here in future phases
		// Example:
		// v1.GET("/search/trips", searchController.SearchTrips)
		// v1.GET("/search/popular-routes", searchController.GetPopularRoutes)
		_ = v1 // Prevent unused variable error
	}
}

// healthCheck returns the health status of the search-api service
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"status":    "ok",
			"service":   "search-api",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
