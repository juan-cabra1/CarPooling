package routes

import (
	"net/http"
	"time"
	"trips-api/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(router *gin.Engine, tripController controller.TripController, chatController *controller.ChatController, jwtMiddleware gin.HandlerFunc) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Rutas públicas de trips (sin autenticación)
	router.GET("/trips", tripController.ListTrips)
	router.GET("/trips/:id", tripController.GetTrip)

	// Rutas protegidas de trips (requieren autenticación)
	protected := router.Group("/trips")
	protected.Use(jwtMiddleware)
	{
		protected.POST("", tripController.CreateTrip)
		protected.PUT("/:id", tripController.UpdateTrip)
		protected.PATCH("/:id", tripController.UpdateTrip)
		protected.DELETE("/:id", tripController.DeleteTrip)

		// Chat routes (protected - requires authentication)
		protected.POST("/:id/messages", chatController.SendMessage)
		protected.GET("/:id/messages", chatController.GetMessages)
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
