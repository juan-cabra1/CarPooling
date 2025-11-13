package routes

import (
	"search-api/internal/controller"
	"search-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the search-api
func SetupRoutes(
	router *gin.Engine,
	healthController *controller.HealthController,
	searchController *controller.SearchController,
) {
	// Apply global middlewares
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", healthController.HealthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Search endpoints (all public, no auth required)
		v1.GET("/search/trips", searchController.SearchTrips)
		v1.GET("/search/location", searchController.SearchByLocation)
		v1.GET("/search/autocomplete", searchController.GetAutocomplete)
		v1.GET("/search/popular-routes", searchController.GetPopularRoutes)

		// Trip detail endpoint
		v1.GET("/trips/:id", searchController.GetTrip)
	}
}
