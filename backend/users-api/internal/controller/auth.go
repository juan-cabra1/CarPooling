package controller

import (
	"users-api/internal/domain"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthController define la interfaz del controlador de autenticación
type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	VerifyEmail(c *gin.Context)
	ResendVerificationEmail(c *gin.Context)
	RequestPasswordReset(c *gin.Context)
	ResetPassword(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type authController struct {
	authService service.AuthService
}

// NewAuthController crea una nueva instancia del controlador de autenticación
func NewAuthController(authService service.AuthService) AuthController {
	return &authController{authService: authService}
}

// Register maneja el registro de nuevos usuarios
// POST /users
func (ctrl *authController) Register(c *gin.Context) {
	var req domain.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	user, err := ctrl.authService.Register(req)
	if err != nil {
		// Si el email ya existe, retornar 409 Conflict
		if err.Error() == "el email ya está registrado" {
			c.JSON(409, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    user,
	})
}

// Login maneja la autenticación de usuarios
// POST /login
func (ctrl *authController) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	response, err := ctrl.authService.Login(req)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    response,
	})
}

// VerifyEmail verifica el email del usuario
// GET /verify-email?token=xxx
func (ctrl *authController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "token requerido",
		})
		return
	}

	if err := ctrl.authService.VerifyEmail(token); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "email verificado exitosamente"},
	})
}

// ResendVerificationEmail reenvía el email de verificación
// POST /resend-verification
func (ctrl *authController) ResendVerificationEmail(c *gin.Context) {
	var req domain.ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	if err := ctrl.authService.ResendVerificationEmail(req.Email); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "email de verificación enviado"},
	})
}

// RequestPasswordReset solicita el restablecimiento de contraseña
// POST /forgot-password
func (ctrl *authController) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	if err := ctrl.authService.RequestPasswordReset(req.Email); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "si el email existe, recibirás instrucciones para restablecer tu contraseña"},
	})
}

// ResetPassword restablece la contraseña usando el token
// POST /reset-password
func (ctrl *authController) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	if err := ctrl.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "contraseña restablecida exitosamente"},
	})
}

// ChangePassword cambia la contraseña del usuario autenticado
// POST /change-password (requiere autenticación)
func (ctrl *authController) ChangePassword(c *gin.Context) {
	var req domain.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	// Obtener user_id del contexto (viene del middleware JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	if err := ctrl.authService.ChangePassword(userID.(int64), req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    gin.H{"message": "contraseña cambiada exitosamente"},
	})
}
