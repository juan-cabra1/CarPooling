package domain

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extrae el user_id del contexto de Gin de forma segura
// Retorna el user_id si existe, o un error si no se encuentra o tiene tipo incorrecto
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	// Obtener el valor del contexto
	value, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user_id not found in context")
	}

	// Validar el tipo
	userID, ok := value.(int64)
	if !ok {
		return 0, fmt.Errorf("user_id has invalid type: %T", value)
	}

	return userID, nil
}

// GetEmailFromContext extrae el email del contexto de Gin de forma segura
// Retorna el email si existe, o un error si no se encuentra o tiene tipo incorrecto
func GetEmailFromContext(c *gin.Context) (string, error) {
	// Obtener el valor del contexto
	value, exists := c.Get("email")
	if !exists {
		return "", fmt.Errorf("email not found in context")
	}

	// Validar el tipo
	email, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("email has invalid type: %T", value)
	}

	return email, nil
}

// GetRoleFromContext extrae el role del contexto de Gin de forma segura
// Retorna el role si existe, o un error si no se encuentra o tiene tipo incorrecto
func GetRoleFromContext(c *gin.Context) (string, error) {
	// Obtener el valor del contexto
	value, exists := c.Get("role")
	if !exists {
		return "", fmt.Errorf("role not found in context")
	}

	// Validar el tipo
	role, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("role has invalid type: %T", value)
	}

	return role, nil
}
