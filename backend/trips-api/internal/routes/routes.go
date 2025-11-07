package routes

import (
	"net/http"
	"time"
	"trips-api/internal/controller"
	"trips-api/internal/middleware"
	"trips-api/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	router *gin.Engine,
	tripController controller.TripController,
	authService service.AuthService,
) {
	// Middlewares globales
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Rutas públicas de trips (sin autenticación)
		v1.GET("/trips/:id", tripController.GetTrip)
		v1.GET("/trips", tripController.ListTrips)

		// Rutas protegidas de trips (requieren autenticación)
		protectedTrips := v1.Group("/trips")
		protectedTrips.Use(middleware.AuthMiddleware(authService))
		{
			protectedTrips.POST("", tripController.CreateTrip)
			protectedTrips.PUT("/:id", tripController.UpdateTrip)
			protectedTrips.DELETE("/:id", tripController.DeleteTrip)
			protectedTrips.PATCH("/:id/cancel", tripController.CancelTrip)
		}
	}
}

// healthCheck maneja el endpoint de health check
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":    "ok",
			"service":   "trips-api",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
