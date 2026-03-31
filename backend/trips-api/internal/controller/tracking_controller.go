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

func NewTrackingController(trackingService service.TrackingService) TrackingController {
	return &trackingController{trackingService: trackingService}
}

func extractUserID(c *gin.Context) (int64, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"success": false, "error": "usuario no autenticado"})
		return 0, false
	}
	if id, ok := userIDVal.(int64); ok {
		return id, true
	}
	if f, ok := userIDVal.(float64); ok {
		return int64(f), true
	}
	c.JSON(400, gin.H{"success": false, "error": "ID de usuario inválido"})
	return 0, false
}

// StartTrip — POST /trips/:id/start
func (ctrl *trackingController) StartTrip(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		return
	}
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{"success": false, "error": "trip_id es requerido"})
		return
	}
	if err := ctrl.trackingService.StartTrip(c.Request.Context(), tripID, userID); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to start trip")
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Viaje iniciado exitosamente"})
}

// UpdateLocation — POST /trips/:id/location
func (ctrl *trackingController) UpdateLocation(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		return
	}
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{"success": false, "error": "trip_id es requerido"})
		return
	}

	var req struct {
		Lat       float64 `json:"lat" binding:"required"`
		Lng       float64 `json:"lng" binding:"required"`
		Speed     float64 `json:"speed"`
		Heading   float64 `json:"heading"`
		Timestamp int64   `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "datos de ubicación inválidos: " + err.Error()})
		return
	}

	location := domain.LocationPoint{
		Lat:       req.Lat,
		Lng:       req.Lng,
		Speed:     req.Speed,
		Heading:   req.Heading,
		Timestamp: domain.UnixToTime(req.Timestamp),
	}

	if err := ctrl.trackingService.UpdateLocation(c.Request.Context(), tripID, userID, location); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to update location")
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Ubicación actualizada"})
}

// CompleteTrip — POST /trips/:id/complete
func (ctrl *trackingController) CompleteTrip(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		return
	}
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{"success": false, "error": "trip_id es requerido"})
		return
	}
	if err := ctrl.trackingService.CompleteTrip(c.Request.Context(), tripID, userID); err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Int64("user_id", userID).Msg("Failed to complete trip")
		c.JSON(400, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Viaje completado exitosamente"})
}

// GetTripTracking — GET /trips/:id/tracking
func (ctrl *trackingController) GetTripTracking(c *gin.Context) {
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(400, gin.H{"success": false, "error": "trip_id es requerido"})
		return
	}
	trip, err := ctrl.trackingService.GetTripTracking(c.Request.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to get trip tracking")
		c.JSON(404, gin.H{"success": false, "error": "viaje no encontrado"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": trip})
}
