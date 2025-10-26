package controller

import (
	"strconv"

	"github.com/carpooling-ucc/bookings-api/internal/domain"
	"github.com/carpooling-ucc/bookings-api/internal/service"
	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService service.BookingService
}

func NewBookingController(bookingService service.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// CreateBooking creates a new booking
func (ctrl *BookingController) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	booking, err := ctrl.bookingService.CreateBooking(req, userID.(int64))
	if err != nil {
		// Determine status code based on error
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		} else if contains(err.Error(), "not enough") || contains(err.Error(), "already has") || contains(err.Error(), "cannot be the driver") {
			statusCode = 409
		} else if contains(err.Error(), "validation failed") {
			statusCode = 400
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    booking,
	})
}

// GetBookingByID retrieves a booking by ID
func (ctrl *BookingController) GetBookingByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	booking, err := ctrl.bookingService.GetBookingByID(id, userID.(int64))
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		} else if contains(err.Error(), "forbidden") {
			statusCode = 403
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    booking,
	})
}

// CancelBooking cancels a booking
func (ctrl *BookingController) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	booking, err := ctrl.bookingService.CancelBooking(id, userID.(int64))
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		} else if contains(err.Error(), "forbidden") {
			statusCode = 403
		} else if contains(err.Error(), "cannot cancel") || contains(err.Error(), "already") {
			statusCode = 400
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    booking,
	})
}

// ConfirmArrival confirms safe arrival for a booking
func (ctrl *BookingController) ConfirmArrival(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	booking, err := ctrl.bookingService.ConfirmArrival(id, userID.(int64))
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		} else if contains(err.Error(), "forbidden") {
			statusCode = 403
		} else if contains(err.Error(), "must be") {
			statusCode = 400
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    booking,
	})
}

// GetUserBookings retrieves bookings for a user
func (ctrl *BookingController) GetUserBookings(c *gin.Context) {
	userIDParam := c.Param("userId")
	userIDFromToken, _ := c.Get("user_id")

	// Parse user ID from URL
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid user ID",
		})
		return
	}

	// Validate that user can only see their own bookings (unless admin)
	if userID != userIDFromToken.(int64) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(403, gin.H{
				"success": false,
				"error":   "forbidden: you can only view your own bookings",
			})
			return
		}
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	bookings, err := ctrl.bookingService.GetUserBookings(userID, page, limit, status)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    bookings,
	})
}

// GetTripBookings retrieves all bookings for a trip
func (ctrl *BookingController) GetTripBookings(c *gin.Context) {
	tripID := c.Param("tripId")
	userID, _ := c.Get("user_id")

	bookings, err := ctrl.bookingService.GetTripBookings(tripID, userID.(int64))
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		} else if contains(err.Error(), "forbidden") {
			statusCode = 403
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    bookings,
	})
}

// Internal endpoints (no auth required)

func (ctrl *BookingController) GetBookingByIDInternal(c *gin.Context) {
	id := c.Param("id")

	booking, err := ctrl.bookingService.GetBookingByIDInternal(id)
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    booking,
	})
}

func (ctrl *BookingController) GetTripBookingsInternal(c *gin.Context) {
	_ = c.Param("tripId")

	// For internal use, we bypass the driver check
	// This is used by trips-api to get booking information
	c.JSON(501, gin.H{
		"success": false,
		"error":   "not implemented yet",
	})
}

func (ctrl *BookingController) CompleteBookingInternal(c *gin.Context) {
	id := c.Param("id")

	booking, err := ctrl.bookingService.CompleteBookingInternal(id)
	if err != nil {
		statusCode := 500
		if contains(err.Error(), "not found") {
			statusCode = 404
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    booking,
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
