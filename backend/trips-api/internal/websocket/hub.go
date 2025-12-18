package websocket

import (
	"log"
	"sync"
)

// Hub mantiene el conjunto de clientes activos y transmite mensajes a los clientes
type Hub struct {
	// Clientes registrados agrupados por tripID
	// map[tripID]map[*Client]bool
	trips map[string]map[*Client]bool

	// Canal para registrar clientes
	register chan *Client

	// Canal para desregistrar clientes
	unregister chan *Client

	// Canal para broadcast de mensajes a un trip específico
	broadcast chan *Message

	// Mutex para acceso concurrente seguro
	mu sync.RWMutex
}

// Message representa un mensaje a transmitir
type Message struct {
	TripID  string      `json:"trip_id"`
	Event   string      `json:"event"` // location_update, trip_started, trip_completed, passenger_joined
	Payload interface{} `json:"payload"`
}

// NewHub crea una nueva instancia de Hub
func NewHub() *Hub {
	return &Hub{
		trips:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256), // Buffer para evitar bloqueos
	}
}

// Run inicia el loop principal del hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient registra un nuevo cliente en un trip
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	tripID := client.tripID

	// Si el trip no existe, crear el mapa
	if h.trips[tripID] == nil {
		h.trips[tripID] = make(map[*Client]bool)
		log.Printf("Created new trip room: %s", tripID)
	}

	// Agregar cliente al trip
	h.trips[tripID][client] = true
	log.Printf("Client %d joined trip %s (total clients: %d)", client.userID, tripID, len(h.trips[tripID]))
}

// unregisterClient desregistra un cliente y limpia si es necesario
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	tripID := client.tripID

	if clients, ok := h.trips[tripID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			log.Printf("Client %d left trip %s (remaining clients: %d)", client.userID, tripID, len(clients))

			// Si no quedan clientes, eliminar el trip room
			if len(clients) == 0 {
				delete(h.trips, tripID)
				log.Printf("Removed empty trip room: %s", tripID)
			}
		}
	}
}

// broadcastMessage envía un mensaje a todos los clientes de un trip
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	tripID := message.TripID

	if clients, ok := h.trips[tripID]; ok {
		log.Printf("Broadcasting '%s' to %d clients in trip %s", message.Event, len(clients), tripID)

		for client := range clients {
			select {
			case client.send <- message:
				// Mensaje enviado exitosamente
			default:
				// Cliente bloqueado, desregistrar
				log.Printf("Client %d blocked, unregistering", client.userID)
				go func(c *Client) {
					h.unregister <- c
				}(client)
			}
		}
	} else {
		log.Printf("No clients in trip %s to broadcast to", tripID)
	}
}

// BroadcastToTrip envía un mensaje a todos los clientes de un trip (método público)
func (h *Hub) BroadcastToTrip(tripID string, event string, payload interface{}) {
	message := &Message{
		TripID:  tripID,
		Event:   event,
		Payload: payload,
	}

	select {
	case h.broadcast <- message:
		// Mensaje encolado
	default:
		log.Printf("Broadcast channel full, message dropped for trip %s", tripID)
	}
}

// GetTripClientCount devuelve el número de clientes conectados a un trip
func (h *Hub) GetTripClientCount(tripID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.trips[tripID]; ok {
		return len(clients)
	}
	return 0
}

// GetActiveTripIDs devuelve la lista de IDs de trips con clientes activos
func (h *Hub) GetActiveTripIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	tripIDs := make([]string, 0, len(h.trips))
	for tripID := range h.trips {
		tripIDs = append(tripIDs, tripID)
	}

	return tripIDs
}
