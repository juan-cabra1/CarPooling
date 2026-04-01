package controller

import (
	"strconv"
	"users-api/internal/domain"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// VerificationController maneja las operaciones de verificación de identidad
type VerificationController interface {
	SubmitVerification(c *gin.Context)
	GetVerificationStatus(c *gin.Context)
	GetPendingVerifications(c *gin.Context)
	ReviewVerification(c *gin.Context)
}

type verificationController struct {
	userService service.UserService
}

// NewVerificationController crea una nueva instancia del controlador
func NewVerificationController(userService service.UserService) VerificationController {
	return &verificationController{userService: userService}
}

// SubmitVerification - el usuario envía datos de DNI o licencia
// POST /users/:id/verification
func (ctrl *verificationController) SubmitVerification(c *gin.Context) {
	userIDParam := c.Param("id")
	paramID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "ID inválido"})
		return
	}

	// Solo el propio usuario puede subir su verificación
	authedID, _ := c.Get("user_id")
	if authedID.(int64) != paramID {
		c.JSON(403, gin.H{"success": false, "error": "no autorizado"})
		return
	}

	var req domain.SubmitVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "datos inválidos: " + err.Error()})
		return
	}

	user, err := ctrl.userService.SubmitVerification(paramID, req)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Información enviada correctamente. Será revisada por un administrador.",
		"data":    user,
	})
}

// GetVerificationStatus - el usuario consulta su estado de verificación
// GET /users/:id/verification
func (ctrl *verificationController) GetVerificationStatus(c *gin.Context) {
	userIDParam := c.Param("id")
	paramID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "ID inválido"})
		return
	}

	authedID, _ := c.Get("user_id")
	authedRole, _ := c.Get("role")
	if authedID.(int64) != paramID && authedRole != "admin" {
		c.JSON(403, gin.H{"success": false, "error": "no autorizado"})
		return
	}

	user, err := ctrl.userService.GetUserByID(paramID)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "error": "usuario no encontrado"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"verification_status":  user.VerificationStatus,
			"dni_verified":         user.DNIVerified,
			"dni_verified_at":      user.DNIVerifiedAt,
			"license_verified":     user.LicenseVerified,
			"license_verified_at":  user.LicenseVerifiedAt,
			"rejection_reason":     user.RejectionReason,
		},
	})
}

// GetPendingVerifications - admin obtiene la lista de verificaciones pendientes
// GET /admin/verifications
func (ctrl *verificationController) GetPendingVerifications(c *gin.Context) {
	users, total, err := ctrl.userService.GetPendingVerifications()
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"total": total,
		},
	})
}

// ReviewVerification - admin aprueba o rechaza una verificación
// POST /admin/verifications/:id/review
func (ctrl *verificationController) ReviewVerification(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "error": "ID inválido"})
		return
	}

	var req domain.ReviewVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "datos inválidos: " + err.Error()})
		return
	}

	user, err := ctrl.userService.ReviewVerification(userID, req)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Verificación procesada correctamente",
		"data":    user,
	})
}
