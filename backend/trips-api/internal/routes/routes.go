package routes

import (
	"net/http"
	"time"
	"trips-api/internal/controller"
	"trips-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	router *gin.Engine,
	tripController controller.TripController,
	chatController *controller.ChatController,
	trackingController controller.TrackingController,
	wsController *controller.WSController,
	websocketController controller.WebSocketController,
	routeController controller.RouteController,
	jwtMiddleware gin.HandlerFunc,
) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Endpoint público: ruta Google Maps (el API key se mantiene server-side)
	router.GET("/route", routeController.GetRouteEmbed)

	// Rutas públicas de trips (sin autenticación)
	router.GET("/trips", middleware.RateLimitMiddleware(60), tripController.ListTrips)
	router.GET("/trips/:id", middleware.RateLimitMiddleware(60), tripController.GetTrip)

	// Chat WebSocket — auth via ?token=<JWT> (browsers can't send Authorization headers in WS handshake)
	router.GET("/trips/:id/ws", wsController.HandleWS)

	// Rutas protegidas de trips (requieren autenticación)
	protected := router.Group("/trips")
	protected.Use(jwtMiddleware)
	protected.Use(middleware.RateLimitMiddleware(120))
	{
		protected.POST("", tripController.CreateTrip)
		protected.PUT("/:id", tripController.UpdateTrip)
		protected.PATCH("/:id", tripController.UpdateTrip)
		protected.DELETE("/:id", tripController.DeleteTrip)
		protected.PATCH("/:id/status", tripController.PatchTripStatus)

		// Chat routes
		protected.POST("/:id/messages", chatController.SendMessage)
		protected.GET("/:id/messages", chatController.GetMessages)

		// Tracking routes
		protected.POST("/:id/start", trackingController.StartTrip)
		protected.POST("/:id/location", trackingController.UpdateLocation)
		protected.POST("/:id/complete", trackingController.CompleteTrip)
		protected.GET("/:id/tracking", trackingController.GetTripTracking)
	}

	// Location tracking WebSocket — auth via JWT middleware, trip_id as query param
	wsProtected := router.Group("/ws")
	wsProtected.Use(jwtMiddleware)
	{
		wsProtected.GET("", websocketController.HandleWebSocket)
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
