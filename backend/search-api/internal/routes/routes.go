package routes

import (
	"github.com/gin-gonic/gin"
	"search-api/internal/controller"
)

// SetupRoutes configures all HTTP routes for the search-api service
func SetupRoutes(router *gin.Engine, healthController *controller.HealthController) {
	// Health check endpoint (no authentication required)
	router.GET("/health", healthController.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes - accessible without authentication
		public := v1.Group("/search")
		{
			// TODO: Add public search endpoints here
			// Example: public.GET("/trips", searchController.SearchTrips)
			_ = public // Placeholder to avoid unused variable error
		}

		// Protected routes - require JWT authentication
		protected := v1.Group("")
		// TODO: Add JWT middleware when implementing authentication
		// protected.Use(middleware.JWTAuth())
		{
			// TODO: Add protected endpoints here
			// Example: protected.GET("/favorites", favoritesController.GetFavorites)
			_ = protected // Placeholder to avoid unused variable error
		}
	}

	// Additional route groups for future functionality
	// Admin routes (if needed)
	// admin := router.Group("/admin")
	// admin.Use(middleware.JWTAuth(), middleware.AdminOnly())
	// {
	//     // Admin-only endpoints
	// }
}
