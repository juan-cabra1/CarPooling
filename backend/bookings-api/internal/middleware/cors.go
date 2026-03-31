package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures Cross-Origin Resource Sharing (CORS) headers.
// Allowed origins are read from ALLOWED_ORIGINS env var (comma-separated).
// When credentials are used, browsers require the response to echo the exact
// request origin rather than a wildcard — this middleware handles that correctly.
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
				// No allowlist configured — development mode, allow all without credentials
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
