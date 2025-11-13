package middleware

import (
	"errors"
	"net/http"

	"search-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

		// Log the error with context
		log.Error().
			Err(err).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
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
	// 404 Not Found
	case "TRIP_NOT_FOUND", "USER_NOT_FOUND", "SEARCH_TRIP_NOT_FOUND":
		return http.StatusNotFound

	// 400 Bad Request - validation errors
	case "INVALID_QUERY", "INVALID_GEO_COORDS", "INVALID_INPUT":
		return http.StatusBadRequest

	// 401 Unauthorized
	case "UNAUTHORIZED":
		return http.StatusUnauthorized

	// 503 Service Unavailable - infrastructure errors
	case "SOLR_UNAVAILABLE", "SERVICE_UNAVAILABLE":
		return http.StatusServiceUnavailable

	// Note: CACHE_UNAVAILABLE is not mapped to 503 because cache failures
	// should be handled gracefully without returning error to client

	// 500 Internal Server Error - default for everything else
	default:
		return http.StatusInternalServerError
	}
}
