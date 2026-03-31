package websocket

import (
	"encoding/json"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
)

// Client represents a single WebSocket connection
type Client struct {
	Hub      *Hub
	Conn     *gorillaws.Conn
	Send     chan []byte
	TripID   string
	UserID   int64
	IsDriver bool // true for drivers sending location updates
}

// locationUpdateMsg is the expected shape of a location_update message from the driver
type locationUpdateMsg struct {
	Event  string `json:"event"`
	TripID string `json:"trip_id"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Speed  float64 `json:"speed"`
	Heading float64 `json:"heading"`
	Timestamp int64 `json:"timestamp"`
}

// WritePump pumps messages from the Send channel to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(gorillaws.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(gorillaws.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorillaws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump reads from the connection to handle pong responses and detect disconnection.
// For driver clients (IsDriver=true), incoming messages are processed as location events.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		if c.IsDriver && len(msg) > 0 {
			c.handleDriverMessage(msg)
		}
	}
}

// handleDriverMessage processes location events sent by the driver over WebSocket
func (c *Client) handleDriverMessage(raw []byte) {
	var msg locationUpdateMsg
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}
	switch msg.Event {
	case "location_update":
		tripID := msg.TripID
		if tripID == "" {
			tripID = c.TripID
		}
		c.Hub.BroadcastToTrip(tripID, "location_update", map[string]interface{}{
			"lat":       msg.Lat,
			"lng":       msg.Lng,
			"speed":     msg.Speed,
			"heading":   msg.Heading,
			"timestamp": msg.Timestamp,
		})
	case "trip_start":
		tripID := msg.TripID
		if tripID == "" {
			tripID = c.TripID
		}
		c.Hub.BroadcastToTrip(tripID, "trip_started", map[string]interface{}{
			"trip_id":   tripID,
			"timestamp": time.Now().Unix(),
		})
	case "trip_complete":
		tripID := msg.TripID
		if tripID == "" {
			tripID = c.TripID
		}
		c.Hub.BroadcastToTrip(tripID, "trip_completed", map[string]interface{}{
			"trip_id":   tripID,
			"timestamp": time.Now().Unix(),
		})
	}
}

// ServeWs registers a new tracking WebSocket client and starts its read/write pumps.
func ServeWs(hub *Hub, conn *gorillaws.Conn, userID int64, tripID string, isDriver bool) {
	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		TripID:   tripID,
		UserID:   userID,
		IsDriver: isDriver,
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
