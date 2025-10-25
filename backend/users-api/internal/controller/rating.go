package controller

import (
	"strconv"
	"users-api/internal/domain"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// RatingController define la interfaz del controlador de calificaciones
type RatingController interface {
	CreateRating(c *gin.Context)
	GetUserRatings(c *gin.Context)
}

type ratingController struct {
	ratingService service.RatingService
}

// NewRatingController crea una nueva instancia del controlador de calificaciones
func NewRatingController(ratingService service.RatingService) RatingController {
	return &ratingController{ratingService: ratingService}
}

// CreateRating crea una nueva calificación
// POST /internal/ratings
func (ctrl *ratingController) CreateRating(c *gin.Context) {
	var req domain.CreateRatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	// Crear calificación
	if err := ctrl.ratingService.CreateRating(req); err != nil {
		if err.Error() == "ya existe una calificación para este viaje" {
			c.JSON(400, gin.H{
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

	c.JSON(201, gin.H{
		"success": true,
		"data":    gin.H{"message": "calificación creada exitosamente"},
	})
}

// GetUserRatings obtiene las calificaciones de un usuario
// GET /users/:id/ratings?page=1&limit=10
func (ctrl *ratingController) GetUserRatings(c *gin.Context) {
	// Extraer ID del path
	idParam := c.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Obtener ratings
	ratings, total, err := ctrl.ratingService.GetUserRatings(userID, page, limit)
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
			"ratings": ratings,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}
