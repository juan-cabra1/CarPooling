package middleware

import (
	"payments-api/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware validates JWT token and extracts claims to context
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("Request missing Authorization header")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "token de autenticación requerido",
			})
			c.Abort()
			return
		}

		// Verify "Bearer TOKEN" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("auth_header", authHeader).
				Msg("Invalid Authorization header format")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "formato de token inválido, usar: Bearer TOKEN",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Warn().
				Err(err).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("Token validation failed")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "token inválido o expirado",
			})
			c.Abort()
			return
		}

		// Extract claims with type validations
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user_id with type validation
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				log.Error().
					Str("path", c.Request.URL.Path).
					Interface("user_id_claim", claims["user_id"]).
					Msg("Invalid user_id claim type")
				c.JSON(401, gin.H{
					"success": false,
					"error":   "claims del token inválidos",
				})
				c.Abort()
				return
			}

			// Extract email with type validation
			email, ok := claims["email"].(string)
			if !ok {
				log.Error().
					Str("path", c.Request.URL.Path).
					Interface("email_claim", claims["email"]).
					Msg("Invalid email claim type")
				c.JSON(401, gin.H{
					"success": false,
					"error":   "claims del token inválidos",
				})
				c.Abort()
				return
			}

			// Extract role with type validation
			role, ok := claims["role"].(string)
			if !ok {
				log.Error().
					Str("path", c.Request.URL.Path).
					Interface("role_claim", claims["role"]).
					Msg("Invalid role claim type")
				c.JSON(401, gin.H{
					"success": false,
					"error":   "claims del token inválidos",
				})
				c.Abort()
				return
			}

			// Extract user_uuid if present (for external references)
			userUUID := ""
			if uuid, ok := claims["user_uuid"].(string); ok {
				userUUID = uuid
			}

			// Save claims in context
			c.Set("user_id", int64(userIDFloat))
			c.Set("email", email)
			c.Set("role", role)
			c.Set("user_uuid", userUUID)

			log.Debug().
				Int64("user_id", int64(userIDFloat)).
				Str("email", email).
				Str("role", role).
				Str("user_uuid", userUUID).
				Str("path", c.Request.URL.Path).
				Msg("JWT authentication successful")
		} else {
			log.Error().
				Str("path", c.Request.URL.Path).
				Msg("Invalid token claims")
			c.JSON(401, gin.H{
				"success": false,
				"error":   "claims del token inválidos",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID extracts user_id from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUserUUID extracts user_uuid from context
func GetUserUUID(c *gin.Context) (string, bool) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		return "", false
	}
	uuid, ok := userUUID.(string)
	return uuid, ok
}

// GetUserEmail extracts email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetUserRole extracts role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	r, ok := role.(string)
	return r, ok
}

// GetUserIdentifier extracts user identifier (UUID if available, otherwise user_id as string)
// This provides compatibility with JWTs that only have user_id
func GetUserIdentifier(c *gin.Context) (string, bool) {
	// Try UUID first
	userUUID, exists := c.Get("user_uuid")
	if exists {
		if uuid, ok := userUUID.(string); ok && uuid != "" {
			return uuid, true
		}
	}

	// Fallback to user_id converted to string
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	if id, ok := userID.(int64); ok {
		return strconv.FormatInt(id, 10), true
	}
	return "", false
}
