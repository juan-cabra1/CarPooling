package middleware

import "github.com/gin-gonic/gin"

// CORSMiddleware configures Cross-Origin Resource Sharing (CORS) headers
// This middleware allows frontend applications running on different origins
// to make requests to the bookings-api
//
// CORS headers configured:
//   - Access-Control-Allow-Origin: Allows requests from any origin (*)
//   - Access-Control-Allow-Methods: Permitted HTTP methods
//   - Access-Control-Allow-Headers: Headers that can be sent in requests
//   - Access-Control-Expose-Headers: Headers that frontend can read
//   - Access-Control-Allow-Credentials: Allows cookies/auth headers
//   - Access-Control-Max-Age: Cache preflight response for 12 hours
//
// Preflight Requests:
//   - Browser sends OPTIONS request before actual request
//   - This middleware responds with 204 No Content for OPTIONS
//   - Tells browser which origins/methods/headers are allowed
//
// Security Note:
//   - Currently allows all origins (Access-Control-Allow-Origin: *)
//   - For production, should restrict to specific frontend domain(s)
//   - Can be made configurable via environment variable
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin
		// TODO: In production, restrict to specific frontend domains
		// e.g., "https://carpooling.example.com"
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow these HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")

		// Allow these request headers
		// Content-Type: For JSON requests
		// Authorization: For JWT tokens
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Expose these response headers to JavaScript
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		// Allow credentials (cookies, authorization headers)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Cache preflight response for 12 hours (43200 seconds)
		// Reduces preflight requests from browser
		c.Writer.Header().Set("Access-Control-Max-Age", "43200")

		// Handle preflight OPTIONS request
		// Browser sends OPTIONS before actual request to check permissions
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}

		// Continue to next middleware/handler
		c.Next()
	}
}
