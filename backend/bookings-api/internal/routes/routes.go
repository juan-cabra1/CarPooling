package routes

import (
	"github.com/carpooling-ucc/bookings-api/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, bookingCtrl *controller.BookingController, authMiddleware gin.HandlerFunc) {
	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "bookings-api is running",
		})
	})

	// Protected routes (require JWT authentication)
	protected := router.Group("/")
	protected.Use(authMiddleware)
	{
		// Booking operations
		protected.POST("/bookings", bookingCtrl.CreateBooking)
		protected.GET("/bookings/:id", bookingCtrl.GetBookingByID)
		protected.PUT("/bookings/:id/cancel", bookingCtrl.CancelBooking)
		protected.POST("/bookings/:id/confirm-arrival", bookingCtrl.ConfirmArrival)
		protected.GET("/bookings/user/:userId", bookingCtrl.GetUserBookings)
		protected.GET("/bookings/trip/:tripId", bookingCtrl.GetTripBookings)
	}

	// Internal routes (no authentication - for microservice communication)
	internal := router.Group("/internal")
	{
		internal.GET("/bookings/:id", bookingCtrl.GetBookingByIDInternal)
		internal.GET("/bookings/trip/:tripId", bookingCtrl.GetTripBookingsInternal)
		internal.PUT("/bookings/:id/complete", bookingCtrl.CompleteBookingInternal)
	}
}
