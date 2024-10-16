package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
type Client struct {
	Conn   *websocket.Conn
	Hub    *Hub
	RoomID string
	Send   chan []byte
}

type Message struct {
	ClientID string `json:"clientID"`
	Content  string `json:"content"`
}

// Room represents a chat room where clients can communicate
type Room struct {
	AllowedClients map[string]bool // List of allowed client IDs
	Clients        map[*Client]bool
	Broadcast      chan []byte
	Register       chan *Client
	Unregister     chan *Client
}

// Hub manages multiple rooms
type Hub struct {
	Rooms map[string]*Room
}

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for simplicity
	},
}

// NewRoom initializes and returns a new room with allowed client IDs
func NewRoom(allowedClients []string) *Room {
	clientMap := make(map[string]bool)
	for _, id := range allowedClients {
		clientMap[id] = true
	}

	return &Room{
		AllowedClients: clientMap,
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan []byte),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
	}
}

// NewHub initializes and returns a new hub
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// Run starts the hub to handle registering, unregistering, and broadcasting messages
func (hub *Hub) Run() {
	for {
		for roomID, room := range hub.Rooms {
			select {
			case client := <-room.Register:
				room.Clients[client] = true
				log.Printf("Client joined room %s", roomID)
			case client := <-room.Unregister:
				if _, ok := room.Clients[client]; ok {
					delete(room.Clients, client)
					close(client.Send)
					log.Printf("Client left room %s", roomID)
				}
			case message := <-room.Broadcast:
				for client := range room.Clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(room.Clients, client)
					}
				}
			}
		}
	}
}

// HandleWebSocket manages the WebSocket connection for each client
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	clientID := r.URL.Query().Get("clientID") // Get the client ID from the query

	if roomID == "" || clientID == "" {
		http.Error(w, "Room ID and Client ID are required", http.StatusBadRequest)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Ensure the room exists in the hub
	room, exists := hub.Rooms[roomID]
	if !exists {
		// Define allowed clients for the room (example data)
		allowedClients := []string{"client1", "client2", "client3"} // Modify as needed
		room = NewRoom(allowedClients)
		hub.Rooms[roomID] = room
		go hub.Run()
	}

	// Check if the client ID is allowed to join the room
	if !room.AllowedClients[clientID] {
		http.Error(w, "Unauthorized: Client ID not allowed in this room", http.StatusForbidden)
		conn.Close()
		return
	}

	client := &Client{
		Conn:   conn,
		Hub:    hub,
		RoomID: roomID,
		Send:   make(chan []byte, 256),
	}
	room.Register <- client

	// Start goroutines for reading from and writing to the WebSocket
	go client.Read()
	go client.Write()
}

// Read listens for incoming WebSocket messages from the client
func (c *Client) Read() {
	defer func() {
		c.Hub.Rooms[c.RoomID].Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Unmarshal the JSON message
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("JSON Unmarshal error:", err)
			continue
		}

		// Log the received message
		log.Printf("Received message from client %s: %s", msg.ClientID, msg.Content)

		// Broadcast the message to clients in the same room
		c.Hub.Rooms[c.RoomID].Broadcast <- message
	}
}

// Write sends outgoing WebSocket messages to the client
func (c *Client) Write() {
	defer c.Conn.Close()

	for message := range c.Send {
		// Assume the message is already a JSON formatted byte slice, so no need to re-marshal
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(hub, w, r)
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
