package routes

import (
	"bookings-api/internal/controller"
	"bookings-api/internal/middleware"
	"bookings-api/internal/service"

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
//   - bookingController: Controller for booking management endpoints
//   - authService: Service for JWT token validation
//
// Route structure:
//   GET  /health              - Service health check (public)
//   GET  /api/v1/bookings     - List all bookings (auth required)
//   GET  /api/v1/bookings/:id - Get specific booking (auth required)
//   POST /api/v1/bookings     - Create new booking (auth required)
//   PATCH /api/v1/bookings/:id/cancel - Cancel booking (auth required)
func SetupRoutes(
	router *gin.Engine,
	healthController *controller.HealthController,
	bookingController *controller.BookingController,
	authService service.AuthService,
) {
	// ============================================================================
	// MIDDLEWARE REGISTRATION
	// ============================================================================
	// Apply CORS middleware first to set headers on all responses
	// This allows frontend applications to make cross-origin requests
	router.Use(middleware.CORSMiddleware())

	// Register error handling middleware globally
	// This must be registered AFTER routes are defined to catch errors from handlers
	// The ErrorHandler middleware:
	//   - Captures errors added via c.Error()
	//   - Maps domain.AppError to appropriate HTTP status codes
	//   - Returns standardized JSON error responses
	router.Use(middleware.ErrorHandler())

	// ============================================================================
	// PUBLIC ROUTES (No authentication required)
	// ============================================================================

	// Health check endpoint
	// Used by load balancers, Kubernetes probes, and monitoring systems
	// Returns: {"status": "ok", "service": "bookings-api", "port": "8003"}
	router.GET("/health", healthController.HealthCheck)

	// ============================================================================
	// API v1 ROUTES (Authentication required)
	// ============================================================================
	// Protected routes that require JWT authentication
	// All booking-related endpoints require a valid JWT token
	// The AuthMiddleware extracts user_id, email, and role from the token

	v1 := router.Group("/api/v1")
	{
		// Booking routes - all protected by JWT authentication
		bookings := v1.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware(authService)) // JWT authentication
		{
			// Booking CRUD endpoints
			bookings.GET("", bookingController.ListBookings)           // List user's bookings
			bookings.GET("/:id", bookingController.GetBooking)         // Get specific booking
			bookings.POST("", bookingController.CreateBooking)         // Create new booking
			bookings.PATCH("/:id/cancel", bookingController.CancelBooking) // Cancel booking
		}
	}
}
