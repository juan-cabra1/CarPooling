package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAdminRole valida que el usuario autenticado tenga rol de administrador
func RequireAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el rol del contexto (establecido por AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "rol no encontrado en el token",
			})
			c.Abort()
			return
		}

		// Verificar que el rol sea admin
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "acceso denegado - se requiere rol de administrador",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
