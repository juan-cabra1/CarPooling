package controller

import (
	"trips-api/internal/domain"
	"trips-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// TrackingController define la interfaz del controlador de tracking
type TrackingController interface {
	StartTrip(c *gin.Context)
	UpdateLocation(c *gin.Context)
	CompleteTrip(c *gin.Context)
	GetTripTracking(c *gin.Context)
}

type trackingController struct {
	trackingService service.TrackingService
}

// NewTrackingController crea una nueva instancia del controlador de tracking
func NewTrackingController(trackingService service.TrackingService) TrackingController {
	return &trackingController{
		trackingService: trackingService,
	}
}

// StartTrip inicia un viaje
// POST /trips/:id/start
// Requiere autenticación (JWT) y ser el conductor
func (ctrl *trackingController) StartTrip(c *gin.Context) {
	// Extraer user_id del contexto
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		if userIDFloat, ok := userIDVal.(float64); ok {
			userID = int64(userIDFloat)
		} else {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}
	}

	// Extraer trip_id del path
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "trip_id es requerido",
		})
		return
	}

	// Iniciar viaje
	err := ctrl.trackingService.StartTrip(c.Request.Context(), tripID, userID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to start trip")
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Viaje iniciado exitosamente",
	})
}

// UpdateLocation actualiza la ubicación del conductor
// POST /trips/:id/location
// Requiere autenticación (JWT) y ser el conductor
func (ctrl *trackingController) UpdateLocation(c *gin.Context) {
	// Extraer user_id del contexto
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		if userIDFloat, ok := userIDVal.(float64); ok {
			userID = int64(userIDFloat)
		} else {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}
	}

	// Extraer trip_id del path
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "trip_id es requerido",
		})
		return
	}

	// Parsear body
	var req struct {
		Lat       float64 `json:"lat" binding:"required"`
		Lng       float64 `json:"lng" binding:"required"`
		Speed     float64 `json:"speed"`
		Heading   float64 `json:"heading"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos de ubicación inválidos: " + err.Error(),
		})
		return
	}

	// Crear LocationPoint
	location := domain.LocationPoint{
		Lat:       req.Lat,
		Lng:       req.Lng,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Timestamp: domain.UnixToTime(req.Timestamp),
	}

	// Actualizar ubicación
	err := ctrl.trackingService.UpdateLocation(c.Request.Context(), tripID, userID, location)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to update location")
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Ubicación actualizada",
	})
}

// CompleteTrip finaliza un viaje
// POST /trips/:id/complete
// Requiere autenticación (JWT) y ser el conductor
func (ctrl *trackingController) CompleteTrip(c *gin.Context) {
	// Extraer user_id del contexto
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		if userIDFloat, ok := userIDVal.(float64); ok {
			userID = int64(userIDFloat)
		} else {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}
	}

	// Extraer trip_id del path
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "trip_id es requerido",
		})
		return
	}

	// Completar viaje
	err := ctrl.trackingService.CompleteTrip(c.Request.Context(), tripID, userID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to complete trip")
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Viaje completado exitosamente",
	})
}

// GetTripTracking obtiene el estado de tracking de un viaje
// GET /trips/:id/tracking
// Requiere autenticación (JWT)
func (ctrl *trackingController) GetTripTracking(c *gin.Context) {
	// Extraer trip_id del path
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "trip_id es requerido",
		})
		return
	}

	// Obtener tracking
	trip, err := ctrl.trackingService.GetTripTracking(c.Request.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to get trip tracking")
		c.JSON(404, gin.H{
			"success": false,
			"error":   "viaje no encontrado",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    trip,
	})
}
