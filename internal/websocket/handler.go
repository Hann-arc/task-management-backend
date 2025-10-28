package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Handler to handle WebSocket connections for notifications
func Handler(c *websocket.Conn) {

	userID := c.Query("user_id")
	if userID == "" {
		c.Close()
		return
	}

	// Register the new connection
	GlobalHub.Register(userID, c)

	// Keep the connection alive
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			GlobalHub.Unregister(userID, c)
			break
		}
	}
}

// Route returns the Fiber handler for WebSocket
func Route() func(*fiber.Ctx) error {
	return websocket.New(Handler)
}
