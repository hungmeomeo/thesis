package ws

import (
	"log"
	"time"
)

type Room struct {
	ID             string
	Clients        map[*Client]bool
	AllowedClients map[string]bool
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan []byte
	stop           chan bool
	messageManager *MessageManager
}

// NewRoom creates and returns a new room with allowed clients
func NewRoom(ID string, allowedClients []string) *Room {
	allowedClientsMap := make(map[string]bool)
	for _, clientID := range allowedClients {
		allowedClientsMap[clientID] = true
	}

	// Create a new MessageManager with a 10-second database interval
	manager := NewMessageManager(ID, 30*time.Second)

	room := &Room{
		ID:             ID,
		Clients:        make(map[*Client]bool),
		AllowedClients: allowedClientsMap,
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan []byte),
		stop:           make(chan bool),
		messageManager: manager,
	}

	go manager.Run() // Start the MessageManager
	return room
}

// Run starts the room to handle registering, unregistering, and broadcasting messages
func (room *Room) Run() {
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
		case message := <-room.Broadcast:
			// Broadcast the message to clients in the room
			for client := range room.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.Clients, client)
				}
			}
			// Update the message manager with the latest message
			room.messageManager.UpdateMessage(message)
		case <-room.stop:
			// Clean up when the room is stopped
			room.messageManager.Stop() // Stop the message manager
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
