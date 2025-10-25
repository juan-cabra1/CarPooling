package controller

import (
	"strconv"
	"users-api/internal/domain"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController define la interfaz del controlador de usuarios
type UserController interface {
	GetUserByID(c *gin.Context)
	GetMe(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

// NewUserController crea una nueva instancia del controlador de usuarios
func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

// GetUserByID obtiene un usuario por su ID
// GET /users/:id
func (ctrl *userController) GetUserByID(c *gin.Context) {
	// Extraer ID del path
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	// Obtener usuario
	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(404, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
	})
}

// GetMe obtiene el perfil del usuario autenticado
// GET /users/me
func (ctrl *userController) GetMe(c *gin.Context) {
	// Extraer user_id del contexto (viene del middleware JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Obtener perfil
	user, err := ctrl.userService.GetUserProfile(userID.(int64))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUser actualiza el perfil de un usuario
// PUT /users/:id
func (ctrl *userController) UpdateUser(c *gin.Context) {
	// Extraer ID del path
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	// Extraer user_id del contexto (viene del middleware JWT)
	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Validar que el usuario solo puede actualizar su propio perfil
	if authUserID.(int64) != id {
		c.JSON(403, gin.H{
			"success": false,
			"error":   "no tienes permiso para actualizar este perfil",
		})
		return
	}

	// Bind request
	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	// Actualizar usuario
	user, err := ctrl.userService.UpdateUser(id, req)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(404, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    user,
	})
}

// DeleteUser elimina un usuario
// DELETE /users/:id
func (ctrl *userController) DeleteUser(c *gin.Context) {
	// Extraer ID del path
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	// Extraer user_id del contexto (viene del middleware JWT)
	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Validar que el usuario solo puede eliminar su propio perfil
	if authUserID.(int64) != id {
		c.JSON(403, gin.H{
			"success": false,
			"error":   "no tienes permiso para eliminar este perfil",
		})
		return
	}

	// Eliminar usuario
	if err := ctrl.userService.DeleteUser(id); err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(404, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "usuario eliminado exitosamente"},
	})
}
