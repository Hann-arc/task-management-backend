package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
}

// Hub manages all active WebSocket connections
type Hub struct {
	clients map[string][]*Client
	mutex   sync.RWMutex
}

// Instance global Hub
var GlobalHub = &Hub{
	clients: make(map[string][]*Client),
}

// register adds a new client connection to the hub
func (h *Hub) Register(userID string, conn *websocket.Conn) {
	client := &Client{Conn: conn, UserID: userID}
	h.mutex.Lock()
	h.clients[userID] = append(h.clients[userID], client)
	h.mutex.Unlock()
	log.Printf("User %s connected. Total: %d", userID, len(h.clients[userID]))
}

// Unregister removes a dead connection
func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	clients := h.clients[userID]
	for i, client := range clients {
		if client.Conn == conn {
			client.Conn.Close()
			h.clients[userID] = append(clients[:i], clients[i+1:]...)
			log.Printf("User %s disconnected. Remaining: %d", userID, len(h.clients[userID]))
			break
		}
	}
}

// SendToUser sends a message to all user connections
func (h *Hub) SendToUser(userID string, message interface{}) {
	h.mutex.RLock()
	clients, exists := h.clients[userID]
	h.mutex.RUnlock()

	if !exists || len(clients) == 0 {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("WebSocket marshal error: %v", err)
		return
	}

	for _, client := range clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("WebSocket send error: %v", err)
			go h.Unregister(userID, client.Conn)
		}
	}
}
