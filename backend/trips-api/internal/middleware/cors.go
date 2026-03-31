package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura CORS para la aplicación.
// ALLOWED_ORIGINS debe ser una lista de orígenes separados por comas.
// Cuando se usan credenciales (cookies/Authorization), el browser requiere
// que el origen se refleje exactamente — no se permite wildcard con credentials.
func CORSMiddleware() gin.HandlerFunc {
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	allowed := map[string]bool{}
	for _, o := range strings.Split(rawOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" {
			if allowed[origin] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Vary", "Origin")
			} else if len(allowed) == 0 {
				// Sin lista configurada → modo desarrollo, permitir todo sin credentials
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
