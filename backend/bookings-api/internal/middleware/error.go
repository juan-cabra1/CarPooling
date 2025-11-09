package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"bookings-api/internal/domain"
)

// ErrorHandler is a middleware that handles errors and returns standardized JSON responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Log the error
		log.Error().
			Err(err).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("Request error")

		// Check if it's an AppError
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			statusCode := mapErrorCodeToHTTPStatus(appErr.Code)

			c.JSON(statusCode, gin.H{
				"success": false,
				"error": gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
					"details": appErr.Details,
				},
			})
			return
		}

		// Default error response for non-AppError errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An internal error occurred",
			},
		})
	}
}

// mapErrorCodeToHTTPStatus maps AppError codes to HTTP status codes
func mapErrorCodeToHTTPStatus(code string) int {
	switch code {
	case "BOOKING_NOT_FOUND", "TRIP_NOT_FOUND":
		return http.StatusNotFound
	case "UNAUTHORIZED":
		return http.StatusForbidden
	case "DUPLICATE_BOOKING":
		return http.StatusConflict
	case "INSUFFICIENT_SEATS", "CANNOT_BOOK_OWN_TRIP", "INVALID_INPUT", "TRIP_NOT_PUBLISHED", "CANNOT_CANCEL_COMPLETED":
		return http.StatusBadRequest
	case "TRIPS_API_UNAVAILABLE":
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
