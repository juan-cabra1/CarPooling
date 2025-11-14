package middleware

import (
	"strings"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware valida el token JWT y extrae los claims al contexto
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
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
			c.JSON(401, gin.H{
				"success": false,
				"error":   "token inválido o expirado",
			})
			c.Abort()
			return
		}

		// Extraer claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Guardar claims en el contexto
			c.Set("user_id", int64(claims["user_id"].(float64)))
			c.Set("email", claims["email"].(string))
			c.Set("role", claims["role"].(string))
		} else {
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
