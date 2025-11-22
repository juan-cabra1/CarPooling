package controller

import (
	"net/http"
	"strconv"

	"bookings-api/internal/domain"
	"bookings-api/internal/service"

	"github.com/gin-gonic/gin"
)

// BookingController handles HTTP requests for booking management
type BookingController struct {
	bookingService service.BookingService
}

// NewBookingController creates a new instance of BookingController
func NewBookingController(bookingService service.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// CreateBooking handles POST /api/v1/bookings
// Creates a new booking for a passenger
func (bc *BookingController) CreateBooking(c *gin.Context) {
	// Extract authenticated user ID from JWT context
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		c.Error(domain.ErrUnauthorized)
		return
	}

	// Bind and validate request body
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Gin automatically returns 400 with validation errors
		c.Error(domain.NewAppError("VALIDATION_ERROR", "Invalid request body", err.Error()))
		return
	}

	// Authorization: user can only create bookings for themselves
	if req.PassengerID != userID {
		c.Error(domain.ErrUnauthorized.WithMessage("You can only create bookings for yourself"))
		return
	}

	// Call service to create booking
	booking, err := bc.bookingService.CreateBooking(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	// Return success response with 201 Created
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    booking,
	})
}

// GetBooking handles GET /api/v1/bookings/:id
// Retrieves a specific booking by ID
// Authorization: Only booking owner or trip driver can access
func (bc *BookingController) GetBooking(c *gin.Context) {
	// Extract authenticated user ID from JWT context
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		c.Error(domain.ErrUnauthorized)
		return
	}

	// Extract booking ID from URL path
	bookingID := c.Param("id")
	if bookingID == "" {
		c.Error(domain.NewAppError("INVALID_BOOKING_ID", "Booking ID is required", nil))
		return
	}

	// Call service to get booking
	booking, err := bc.bookingService.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		c.Error(err)
		return
	}

	// Authorization: verify user is booking owner or trip driver
	// Note: Service layer should handle driver check via trips-api
	// For now, we only check if user is the passenger
	if booking.PassengerID != userID {
		// TODO: Check if user is trip driver (requires trips-api call)
		// For now, only allow passenger to view their own booking
		c.Error(domain.ErrUnauthorized.WithMessage("You can only view your own bookings"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}

// ListBookings handles GET /api/v1/bookings
// Lists all bookings for the authenticated user with pagination
func (bc *BookingController) ListBookings(c *gin.Context) {
	// Extract authenticated user ID from JWT context
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		c.Error(domain.ErrUnauthorized)
		return
	}

	// Parse pagination query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Service layer will validate and correct invalid pagination values
	bookings, err := bc.bookingService.GetPassengerBookings(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	// Return success response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
	})
}

// CancelBooking handles PATCH /api/v1/bookings/:id/cancel
// Cancels a booking
// Authorization: Only booking owner or trip driver can cancel
func (bc *BookingController) CancelBooking(c *gin.Context) {
	// Extract authenticated user ID from JWT context
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		c.Error(domain.ErrUnauthorized)
		return
	}

	// Extract booking ID from URL path
	bookingID := c.Param("id")
	if bookingID == "" {
		c.Error(domain.NewAppError("INVALID_BOOKING_ID", "Booking ID is required", nil))
		return
	}

	// Bind optional cancellation reason from request body
	var req domain.CancelBookingRequest
	// Ignore binding errors since reason is optional
	_ = c.ShouldBindJSON(&req)

	// Call service to cancel booking
	// Service layer handles authorization (passenger OR driver check)
	err = bc.bookingService.CancelBooking(c.Request.Context(), bookingID, userID, req.Reason)
	if err != nil {
		c.Error(err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Booking cancelled successfully",
	})
}

// GetAllBookings handles GET /api/v1/admin/bookings
// Lists all bookings in the system with pagination and filters (admin only)
func (bc *BookingController) GetAllBookings(c *gin.Context) {
	// Parse pagination and filter parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")       // filter by status (optional)
	tripID := c.Query("trip_id")      // filter by trip (optional)
	passengerID, _ := strconv.ParseInt(c.Query("passenger_id"), 10, 64) // filter by passenger (optional)

	// Call service to get all bookings
	bookings, total, err := bc.bookingService.GetAllBookings(c.Request.Context(), page, limit, status, tripID, passengerID)
	if err != nil {
		c.Error(err)
		return
	}

	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)

	// Return success response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"bookings": bookings,
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": totalPages,
			},
		},
	})
}
