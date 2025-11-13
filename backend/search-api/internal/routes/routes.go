package routes

import (
	"search-api/internal/controller"
	"search-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the search-api
func SetupRoutes(router *gin.Engine, healthController *controller.HealthController) {
	// Apply global middlewares
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", healthController.HealthCheck)

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
