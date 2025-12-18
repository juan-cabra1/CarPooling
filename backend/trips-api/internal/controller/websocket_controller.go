package controller

import (
	"net/http"
	"trips-api/internal/service"
	ws "trips-api/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implementar verificación de origen adecuada en producción
		// Por ahora permitir todos los orígenes para desarrollo
		return true
	},
}

// WebSocketController define la interfaz del controlador de WebSocket
type WebSocketController interface {
	HandleWebSocket(c *gin.Context)
}

type websocketController struct {
	hub          *ws.Hub
	tripService  service.TripService
}

// NewWebSocketController crea una nueva instancia del controlador de WebSocket
func NewWebSocketController(hub *ws.Hub, tripService service.TripService) WebSocketController {
	return &websocketController{
		hub:          hub,
		tripService:  tripService,
	}
}

// HandleWebSocket maneja las conexiones WebSocket
// GET /ws
// Requiere autenticación (JWT) y query params: trip_id
func (ctrl *websocketController) HandleWebSocket(c *gin.Context) {
	// Extraer user_id del contexto (viene del middleware JWT)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		log.Error().Msg("User ID not found in context")
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		// Try converting from float64 (common with JSON parsing)
		if userIDFloat, ok := userIDVal.(float64); ok {
			userID = int64(userIDFloat)
		} else {
			log.Error().Msgf("Invalid user ID type: %T", userIDVal)
			c.JSON(400, gin.H{
				"success": false,
				"error":   "ID de usuario inválido",
			})
			return
		}
	}

	// Obtener trip_id de query params
	tripID := c.Query("trip_id")
	if tripID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "trip_id es requerido",
		})
		return
	}

	// Verificar que el viaje existe
	trip, err := ctrl.tripService.GetTrip(c.Request.Context(), tripID)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to get trip")
		c.JSON(404, gin.H{
			"success": false,
			"error":   "viaje no encontrado",
		})
		return
	}

	// Determinar si el usuario es el conductor
	isDriver := trip.DriverID == userID

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade to WebSocket")
		return
	}

	log.Info().
		Int64("user_id", userID).
		Str("trip_id", tripID).
		Bool("is_driver", isDriver).
		Msg("WebSocket connection established")

	// Servir la conexión WebSocket
	ws.ServeWs(ctrl.hub, conn, userID, isDriver)
}
