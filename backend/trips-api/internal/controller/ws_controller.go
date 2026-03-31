package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
	gorillaws "github.com/gorilla/websocket"

	ws "trips-api/internal/websocket"
)

var upgrader = gorillaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins — CORS is handled by the CORS middleware
		return true
	},
}

// WSController handles WebSocket upgrade requests for trip chat
type WSController struct {
	hub       *ws.Hub
	jwtSecret string
}

// NewWSController creates a new WebSocket controller
func NewWSController(hub *ws.Hub, jwtSecret string) *WSController {
	return &WSController{hub: hub, jwtSecret: jwtSecret}
}

// HandleWS upgrades an HTTP connection to WebSocket for trip chat.
// Auth is done via query param: GET /trips/:id/ws?token=<JWT>
// (WebSocket handshake headers can't carry Authorization in browsers)
func (wc *WSController) HandleWS(c *gin.Context) {
	tripID := c.Param("id")
	if tripID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Get JWT from query param
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "missing token",
		})
		return
	}

	// Validate JWT
	secret := wc.jwtSecret
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}

	token, err := jwtlib.Parse(tokenStr, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, jwtlib.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := int64(userIDFloat)

	// Upgrade HTTP → WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrade failure is already handled by gorilla (writes error response)
		return
	}

	client := &ws.Client{
		Hub:    wc.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		TripID: tripID,
		UserID: userID,
	}

	wc.hub.Register <- client

	// Start pumps in goroutines — these run until the connection is closed
	go client.WritePump()
	go client.ReadPump()
}
