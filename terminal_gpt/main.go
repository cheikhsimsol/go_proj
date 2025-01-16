package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"github.com/gorilla/websocket"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GPT interface
type GPT interface {
	// Register a client session by ID
	Register(clientID string, c *websocket.Conn)
	DeRegister(clientID string)
	// SendMessage sends a chat message to LLM
	SendMessage(clientID, message string) error
}

type Session struct {
	ChatSession *genai.ChatSession
	Conn        *websocket.Conn
}

// GeminiAI struct to manage API key, endpoint, and client sessions
type GeminiAI struct {
	client  *genai.Client
	clients map[string]Session
	mu      sync.Mutex // To ensure thread-safe access to the clients map
}

// NewGenAI creates a new GeminiAI instance
func NewGenAI(client *genai.Client) *GeminiAI {

	return &GeminiAI{
		client:  client,
		clients: make(map[string]Session),
	}
}

// Register establishes a WebSocket connection for a given clientID
func (c *GeminiAI) Register(clientID string, ws *websocket.Conn) {
	model := c.client.GenerativeModel("gemini-1.5-flash")
	chat := model.StartChat()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.clients[clientID] = Session{
		chat,
		ws,
	}
}

// SendMessage sends a message through the WebSocket connection
func (c *GeminiAI) SendMessage(clientID, message string) error {
	c.mu.Lock()
	s, exists := c.clients[clientID]
	c.mu.Unlock()

	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}

	iter := s.ChatSession.SendMessageStream(context.Background(), genai.Text(message))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {

					t, isText := part.(genai.Text)

					if !isText {
						log.Println("Unsupported context received")
						continue
					}

					err = s.Conn.WriteMessage(
						websocket.TextMessage,
						[]byte(t),
					)
				}
			}
		}
	}

	return nil
}

// DeRegister removes a client connection
func (c *GeminiAI) DeRegister(clientID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check to prevent any panics
	if conn, exists := c.clients[clientID]; exists {
		// chat history can also be backed up too.
		// gracefully close the connection
		conn.Conn.Close()
		delete(c.clients, clientID)
	}
}

// ServeWs handles WebSocket connections,
// registers clients and sessions with gemini
func ServeWs(h GPT) http.HandlerFunc {
	// WebSocket upgrader configuration
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

		// Register the client
		// openai is websocket conn
		h.Register(clientID, conn)

		defer h.DeRegister(clientID)

		// Listen for incoming messages (optional, can be used to receive client data)
		// this loop will hold function until broken by rec. error.
		for {

			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			if len(message) > 0 {
				err = h.SendMessage(clientID, string(message))
				if err != nil {
					log.Println("error sending message", err)
					break
				}
			}

		}
	}
}

func main() {

	client, err := genai.NewClient(
		context.Background(),
		option.WithAPIKey(os.Getenv("GEMINI_KEY")),
	)

	if err != nil {
		log.Fatalf("error constructing gemini client: %s", err)
	}

	defer client.Close()

	gpt := NewGenAI(client)

	http.Handle("/chat", ServeWs(gpt))

	log.Println("Listenning at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
