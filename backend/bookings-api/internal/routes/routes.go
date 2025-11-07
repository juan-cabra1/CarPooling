package routes

import (
	"bookings-api/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all HTTP routes for the bookings-api service
//
// This function registers all endpoints with the Gin router.
// Routes are organized by functionality:
//   - Health check routes (for monitoring and load balancers)
//   - Booking management routes (CRUD operations)
//   - User-specific routes (user's bookings)
//
// Parameters:
//   - router: The Gin engine instance to register routes on
//   - healthController: Controller for health check endpoints
//
// Route structure:
//   GET  /health              - Service health check (public)
//   GET  /api/v1/bookings     - List all bookings (auth required)
//   GET  /api/v1/bookings/:id - Get specific booking (auth required)
//   POST /api/v1/bookings     - Create new booking (auth required)
//   PUT  /api/v1/bookings/:id - Update booking (auth required)
//   DELETE /api/v1/bookings/:id - Cancel booking (auth required)
func SetupRoutes(router *gin.Engine, healthController *controller.HealthController) {
	// ============================================================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================================================

	// Health check endpoint
	// Used by load balancers, Kubernetes probes, and monitoring systems
	// Returns: {"status": "ok", "service": "bookings-api", "port": "8003"}
	router.GET("/health", healthController.HealthCheck)

	// ============================================================================
	// API v1 ROUTES (Authentication required - to be added later)
	// ============================================================================
	// Future routes will be added here:
	//
	// v1 := router.Group("/api/v1")
	// {
	//     // Booking routes
	//     bookings := v1.Group("/bookings")
	//     bookings.Use(authMiddleware) // JWT authentication
	//     {
	//         bookings.GET("", bookingController.ListBookings)
	//         bookings.GET("/:id", bookingController.GetBooking)
	//         bookings.POST("", bookingController.CreateBooking)
	//         bookings.PUT("/:id", bookingController.UpdateBooking)
	//         bookings.DELETE("/:id", bookingController.CancelBooking)
	//     }
	// }
}
