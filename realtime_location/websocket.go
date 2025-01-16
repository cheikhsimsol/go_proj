package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketHub is the implementation of the Hub interface
type WebSocketHub struct {
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

// NewWebSocketHub initializes a new WebSocketHub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
	}
}

// Push sends a payload to all connected clients
func (hub *WebSocketHub) Push(payload Location) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for clientID, conn := range hub.clients {
		if err := conn.WriteJSON(payload); err != nil {
			// Handle errors gracefully, such as removing disconnected clients
			conn.Close()
			delete(hub.clients, clientID)
			return err
		}
	}

	return nil
}

// Register adds a new client connection to the hub
func (hub *WebSocketHub) Register(clientID string, conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.clients[clientID] = conn
}

// DeRegister removes a client connection from the hub
func (hub *WebSocketHub) DeRegister(clientID string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	// check to prevent any panics
	if conn, exists := hub.clients[clientID]; exists {
		// gracefully close the connection
		conn.Close()
		delete(hub.clients, clientID)
	}
}
