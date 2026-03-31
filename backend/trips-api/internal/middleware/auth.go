package middleware

import (
	"strings"
	"trips-api/internal/service"

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
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(401, gin.H{"success": false, "error": "claims del token inválidos"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(401, gin.H{"success": false, "error": "token inválido: user_id faltante"})
			c.Abort()
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(401, gin.H{"success": false, "error": "token inválido: email faltante"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(401, gin.H{"success": false, "error": "token inválido: role faltante"})
			c.Abort()
			return
		}

		c.Set("user_id", int64(userIDFloat))
		c.Set("email", email)
		c.Set("role", role)
		if name, ok := claims["name"].(string); ok {
			c.Set("user_name", name)
		}

		c.Next()
	}
}
