package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader configures the WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for simplicity
	},
}

// HandleWebSocket manages the WebSocket connection for each client
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	clientID := r.URL.Query().Get("clientID")

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
		// Create a new room with an allowed clients list
		allowedClients := []string{"client1", "client2", "client3"} // Modify as needed
		room = NewRoom(allowedClients)
		hub.Rooms[roomID] = room
		go room.Run() // Start the room's goroutine, not the hub's
	}

	// Check if the client ID is allowed to join the room
	if !room.AllowedClients[clientID] {
		http.Error(w, "Unauthorized: Client ID not allowed in this room", http.StatusForbidden)
		conn.Close()
		return
	}

	// Create and register a new client
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
