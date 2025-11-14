package middleware

import (
	"strings"
	"users-api/internal/repository"
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

// RequireVerifiedEmail valida que el usuario tenga su email verificado
// Este middleware debe usarse DESPUÉS de AuthMiddleware
func RequireVerifiedEmail(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el user_id del contexto (viene de AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "no autenticado",
			})
			c.Abort()
			return
		}

		// Buscar el usuario en la base de datos
		user, err := userRepo.FindByID(userID.(int64))
		if err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "usuario no encontrado",
			})
			c.Abort()
			return
		}

		// Verificar si el email está verificado
		if !user.EmailVerified {
			c.JSON(403, gin.H{
				"success": false,
				"error":   "debes verificar tu correo electrónico para acceder a esta funcionalidad",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
