package middleware

import (
	"github.com/gin-gonic/gin"
)

// ErrorHandler maneja los errores globales de la aplicación
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verificar si hay errores
		if len(c.Errors) > 0 {
			// Tomar el último error
			err := c.Errors.Last()

			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
}
