package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub interface {
	Push(Location) error
	Register(string, *websocket.Conn)
	DeRegister(string)
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DeviceId  string  `json:"device_id"`
}

// ServeWs handles WebSocket connections,
// registers clients, and manages disconnections
func ServeWs(h Hub) http.HandlerFunc {
	// WebSocket upgrader configuration
	// Constructing here and not in returned function,
	// because doing so would result in a new
	// upgrader on each request.
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins; customize as needed
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Generate a unique client ID (for simplicity, using remote address)
		clientID := r.RemoteAddr

		// Register the client in the Hub
		h.Register(clientID, conn)

		// Ensure the client is deregistered on connection close. ending lifecycle
		defer h.DeRegister(clientID)

		// Listen for incoming messages (optional, can be used to receive client data)
		// this loop will hold function until broken by rec. error.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Exit loop on connection error or close
				break
			}
		}
	}
}

func PushLocation(h Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Ensure the request is a POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the JSON body into a Location struct
		var location Location
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid client id", http.StatusBadRequest)
			return
		}

		// Naturally, the id should come from
		// a JWT or session value. Don't do in production.
		location.DeviceId = ip

		// Push the parsed location to the hub
		if err := h.Push(location); err != nil {
			http.Error(w, "Failed to push location", http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Location pushed successfully"))
	}
}

func main() {

	h := NewWebSocketHub()

	http.Handle("/push", PushLocation(h))
	http.Handle("/track", ServeWs(h))

	log.Println("Listenning at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
