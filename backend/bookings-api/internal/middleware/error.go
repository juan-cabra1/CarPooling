package middleware

import (
	"github.com/gin-gonic/gin"
)

// ErrorHandler handles global application errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
}
