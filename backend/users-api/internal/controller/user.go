package controller

import (
	"strconv"
	"users-api/internal/domain"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController define la interfaz del controlador de usuarios
type UserController interface {
	GetAllUsers(c *gin.Context)
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

// GetAllUsers obtiene todos los usuarios (solo admin)
// GET /admin/users
func (ctrl *userController) GetAllUsers(c *gin.Context) {
	// Parsear parámetros de query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")       // filtro opcional: "user" o "admin"
	search := c.Query("search")   // búsqueda por email o nombre

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Obtener usuarios con paginación
	users, total, err := ctrl.userService.GetAllUsers(page, limit, role, search)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
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

	// Extraer role del contexto
	authRole, roleExists := c.Get("role")
	if !roleExists {
		authRole = "user" // default
	}

	// Validar que el usuario puede actualizar: es admin O es su propio perfil
	if authRole != "admin" && authUserID.(int64) != id {
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

	// Extraer role del contexto
	authRole, roleExists := c.Get("role")
	if !roleExists {
		authRole = "user" // default
	}

	// Validar que el usuario puede eliminar: es admin O es su propio perfil
	if authRole != "admin" && authUserID.(int64) != id {
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
