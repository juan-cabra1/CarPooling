package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ErrorMiddleware handles panics and unexpected errors
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered in ErrorMiddleware")

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "error interno del servidor",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
