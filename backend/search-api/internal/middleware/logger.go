package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger is a middleware that logs HTTP requests with structured logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		durationMs := int(duration.Milliseconds())

		// Get status code
		statusCode := c.Writer.Status()

		// Build log event
		logEvent := log.Info()

		// Add error level if status code indicates error
		if statusCode >= 500 {
			logEvent = log.Error()
		} else if statusCode >= 400 {
			logEvent = log.Warn()
		}

		// Build log with context
		logEvent.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Int("duration_ms", durationMs).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		// Add query parameters if present
		if raw != "" {
			logEvent.Str("query", raw)
		}

		// Add result count if available in context (for search queries)
		if resultsCount, exists := c.Get("results_count"); exists {
			if count, ok := resultsCount.(int); ok {
				logEvent.Int("results_count", count)
			}
		}

		// Add source if available in context (cache/solr/mongodb)
		if source, exists := c.Get("source"); exists {
			if sourceStr, ok := source.(string); ok {
				logEvent.Str("source", sourceStr)
			}
		}

		logEvent.Msg("HTTP request")
	}
}
