package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tiempo máximo de escritura
	writeWait = 10 * time.Second

	// Tiempo máximo de lectura del siguiente pong
	pongWait = 60 * time.Second

	// Intervalo de ping (debe ser menor que pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo de mensaje
	maxMessageSize = 512
)

// Client representa un cliente WebSocket conectado
type Client struct {
	// Hub al que pertenece este cliente
	hub *Hub

	// Conexión WebSocket
	conn *websocket.Conn

	// Canal buffered para mensajes salientes
	send chan *Message

	// ID del usuario
	userID int64

	// ID del trip al que está conectado
	tripID string

	// Si el cliente es el conductor
	isDriver bool
}

// ClientMessage representa un mensaje recibido del cliente
type ClientMessage struct {
	Event   string          `json:"event"`
	TripID  string          `json:"trip_id"`
	Payload json.RawMessage `json:"payload"`
}

// LocationUpdatePayload representa la carga útil de una actualización de ubicación
type LocationUpdatePayload struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
	Timestamp int64   `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implementar verificación de origen adecuada en producción
		return true
	},
}

// readPump lee mensajes desde la conexión WebSocket
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parsear mensaje del cliente
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error parsing client message: %v", err)
			continue
		}

		// Procesar mensaje según el evento
		c.handleClientMessage(&clientMsg)
	}
}

// writePump escribe mensajes a la conexión WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Enviar mensaje como JSON
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleClientMessage procesa mensajes recibidos del cliente
func (c *Client) handleClientMessage(msg *ClientMessage) {
	switch msg.Event {
	case "join_trip":
		// Cliente se une a un trip room
		c.tripID = msg.TripID
		c.hub.register <- c
		log.Printf("User %d joined trip %s", c.userID, c.tripID)

	case "leave_trip":
		// Cliente sale del trip room
		c.hub.unregister <- c
		log.Printf("User %d left trip %s", c.userID, c.tripID)

	case "location_update":
		// Solo conductores pueden enviar actualizaciones de ubicación
		if !c.isDriver {
			log.Printf("Non-driver user %d attempted to send location update", c.userID)
			return
		}

		// Parsear payload de ubicación
		var payload LocationUpdatePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("Error parsing location update: %v", err)
			return
		}

		// Broadcast a todos los clientes del trip
		c.hub.BroadcastToTrip(msg.TripID, "location_update", payload)

	case "trip_start":
		// Solo conductores pueden iniciar viaje
		if !c.isDriver {
			log.Printf("Non-driver user %d attempted to start trip", c.userID)
			return
		}

		c.hub.BroadcastToTrip(msg.TripID, "trip_started", map[string]interface{}{
			"trip_id":   msg.TripID,
			"timestamp": time.Now().Unix(),
		})

	case "trip_complete":
		// Solo conductores pueden completar viaje
		if !c.isDriver {
			log.Printf("Non-driver user %d attempted to complete trip", c.userID)
			return
		}

		c.hub.BroadcastToTrip(msg.TripID, "trip_completed", map[string]interface{}{
			"trip_id":   msg.TripID,
			"timestamp": time.Now().Unix(),
		})

	default:
		log.Printf("Unknown event type: %s", msg.Event)
	}
}

// ServeWs maneja las solicitudes WebSocket del cliente
func ServeWs(hub *Hub, conn *websocket.Conn, userID int64, isDriver bool) {
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan *Message, 256),
		userID:   userID,
		isDriver: isDriver,
	}

	// Iniciar goroutines de lectura y escritura
	go client.writePump()
	go client.readPump()
}
