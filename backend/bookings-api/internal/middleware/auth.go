package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT token and extracts claims to context
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "authentication token required",
			})
			c.Abort()
			return
		}

		// Verify "Bearer TOKEN" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "invalid token format, use: Bearer TOKEN",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Save claims to context
			c.Set("user_id", int64(claims["user_id"].(float64)))
			c.Set("email", claims["email"].(string))
			c.Set("role", claims["role"].(string))
		} else {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "invalid token claims",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
