package middleware

import (
	"bookings-api/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware valida el token JWT y extrae los claims al contexto
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
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

		// Verificar formato "Bearer TOKEN"
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

		// Validar el token
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

		// Extraer claims con validaciones de tipo
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extraer user_id con validación de tipo
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

			// Extraer email con validación de tipo
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

			// Extraer role con validación de tipo
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

			// Guardar claims en el contexto
			c.Set("user_id", int64(userIDFloat))
			c.Set("email", email)
			c.Set("role", role)

			log.Debug().
				Int64("user_id", int64(userIDFloat)).
				Str("email", email).
				Str("role", role).
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
