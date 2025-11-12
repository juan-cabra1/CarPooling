package middleware

import "github.com/gin-gonic/gin"

// ErrorHandler middleware handles errors that occur during request processing
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if any errors were attached during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
}
