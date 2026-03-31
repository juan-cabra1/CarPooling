package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// InternalAuthMiddleware protege los endpoints /internal/* con un token compartido
// entre servicios. Cada servicio que llame a estos endpoints debe enviar
// el header X-Internal-Token con el valor de INTERNAL_API_TOKEN.
func InternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := os.Getenv("INTERNAL_API_TOKEN")
		if expected == "" {
			// Si no hay token configurado, rechazar siempre (fail-secure)
			c.JSON(401, gin.H{"success": false, "error": "servicio no configurado para llamadas internas"})
			c.Abort()
			return
		}

		token := c.GetHeader("X-Internal-Token")
		if token != expected {
			c.JSON(401, gin.H{"success": false, "error": "token interno inválido"})
			c.Abort()
			return
		}

		c.Next()
	}
}
