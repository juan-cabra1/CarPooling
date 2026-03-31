package controller

import (
	"net/http"
	"trips-api/internal/service"
	ws "trips-api/internal/websocket"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var trackingUpgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketController handles WebSocket connections for real-time location tracking
type WebSocketController interface {
	HandleWebSocket(c *gin.Context)
}

type websocketController struct {
	hub         *ws.Hub
	tripService service.TripService
}

func NewWebSocketController(hub *ws.Hub, tripService service.TripService) WebSocketController {
	return &websocketController{hub: hub, tripService: tripService}
}

// HandleWebSocket — GET /ws?trip_id=<id>
// Requires JWT middleware. Drivers can send location_update events; passengers receive them.
func (ctrl *websocketController) HandleWebSocket(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		log.Error().Msg("User ID not found in context")
		c.JSON(401, gin.H{"success": false, "error": "usuario no autenticado"})
		return
	}

	var userID int64
	switch v := userIDVal.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(400, gin.H{"success": false, "error": "ID de usuario inválido"})
		return
	}

	tripID := c.Query("trip_id")
	if tripID == "" {
		c.JSON(400, gin.H{"success": false, "error": "trip_id es requerido"})
		return
	}

	// Verify trip exists and determine if user is the driver
	trip, err := ctrl.tripService.GetTrip(c.Request.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to get trip")
		c.JSON(404, gin.H{"success": false, "error": "viaje no encontrado"})
		return
	}

	isDriver := trip.DriverID == userID

	conn, err := trackingUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade to WebSocket")
		return
	}

	log.Info().
		Int64("user_id", userID).
		Str("trip_id", tripID).
		Bool("is_driver", isDriver).
		Msg("WebSocket tracking connection established")

	ws.ServeWs(ctrl.hub, conn, userID, tripID, isDriver)
}
