package websocket

import (
	"encoding/json"
	"sync"
)

// Hub manages WebSocket connections grouped by trip room
type Hub struct {
	// rooms maps trip_id → set of connected clients
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex

	// Channels for register/unregister/broadcast operations
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan BroadcastMessage
}

// BroadcastMessage holds a payload to send to all clients of a trip
type BroadcastMessage struct {
	TripID  string
	Payload []byte
}

// NewHub creates and returns a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		broadcast:  make(chan BroadcastMessage, 256),
	}
}

// Run processes register/unregister/broadcast events in a dedicated goroutine.
// Must be called with `go hub.Run()`.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.TripID] == nil {
				h.rooms[client.TripID] = make(map[*Client]bool)
			}
			h.rooms[client.TripID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.TripID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.TripID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.TripID]
			h.mu.RUnlock()
			for client := range clients {
				select {
				case client.Send <- msg.Payload:
				default:
					// Client send buffer full — disconnect slow client
					h.Unregister <- client
				}
			}
		}
	}
}

// Broadcast sends a JSON payload to all WebSocket clients connected to the trip room
func (h *Hub) Broadcast(tripID string, payload []byte) {
	h.broadcast <- BroadcastMessage{TripID: tripID, Payload: payload}
}

// BroadcastToTrip serializes an event+payload struct and broadcasts it to all trip clients
func (h *Hub) BroadcastToTrip(tripID string, event string, payload interface{}) {
	data, err := json.Marshal(map[string]interface{}{
		"event":   event,
		"payload": payload,
	})
	if err != nil {
		return
	}
	h.Broadcast(tripID, data)
}

// GetTripClientCount returns the number of clients connected to a trip room
func (h *Hub) GetTripClientCount(tripID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[tripID])
}
