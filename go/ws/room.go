package ws

import (
	"demo/db"
	"log"
	"time"
)

type Room struct {
	ID             string
	Clients        map[*Client]bool
	AllowedClients map[string]bool // Tracks allowed clients
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan []byte
	stop           chan bool
	latestMessage  []byte    // To store the latest broadcast message
	lastUpdateTime time.Time // To track the last update time
}

func NewRoom(id string, allowedClients []string) *Room {
	allowedClientsMap := make(map[string]bool)
	for _, clientID := range allowedClients {
		allowedClientsMap[clientID] = true
	}

	return &Room{
		ID:             id,
		Clients:        make(map[*Client]bool),
		AllowedClients: allowedClientsMap,
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan []byte),
		stop:           make(chan bool),
		lastUpdateTime: time.Now(), // Initialize with the current time
	}
}

// Run starts the room to handle registering, unregistering, and broadcasting messages
func (room *Room) Run() {
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
			log.Printf("Client joined room %s", room.ID)

		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
				log.Printf("Client left room %s", room.ID)
			}

			// Check if no clients are left, and stop the room if so
			if len(room.Clients) == 0 {
				room.Stop()
				log.Printf("No clients left in room %s, stopping the room.", room.ID)
			}

		case message := <-room.Broadcast:
			room.latestMessage = message     // Store the latest broadcast message
			room.lastUpdateTime = time.Now() // Update the last update time
			for client := range room.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.Clients, client)
				}
			}

		case <-ticker.C:
			// Check if 10 seconds have passed since the last update
			if time.Since(room.lastUpdateTime) >= 10*time.Second {
				if room.latestMessage != nil {
					// Store the latest broadcast message in Redis
					db.ConnectRedis()
					err := db.SetRedis(room.ID, string(room.latestMessage))
					if err != nil {
						log.Printf("Error storing message in Redis for room %s: %v", room.ID, err)
					} else {
						log.Printf("Stored message in Redis for room %s", room.ID)
					}
					db.CloseRedis()
					// Reset the latestMessage after storing
					room.latestMessage = nil
				}
			}

		case <-room.stop:
			// Clean up when the room is stopped
			for client := range room.Clients {
				close(client.Send)
				delete(room.Clients, client)
			}
			log.Printf("Room %s stopped", room.ID)
			return
		}
	}
}

// Stop gracefully stops the room's loop
func (room *Room) Stop() {
	room.stop <- true
}
