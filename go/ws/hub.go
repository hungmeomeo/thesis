package ws

import "log"

// Hub manages multiple rooms
type Hub struct {
	Rooms map[string]*Room
	// Channels to manage room creation/deletion at the hub level
	CreateRoom chan *Room
	DeleteRoom chan string
}

// NewHub initializes and returns a new hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		CreateRoom: make(chan *Room),
		DeleteRoom: make(chan string),
	}
}

// Run starts the hub to handle creating and deleting rooms
func (hub *Hub) Run() {
	for {
		select {
		case room := <-hub.CreateRoom:
			hub.Rooms[room.ID] = room
			go room.Run() // Run a goroutine for the room
			log.Printf("Room %s created", room.ID)
		case roomID := <-hub.DeleteRoom:
			if room, ok := hub.Rooms[roomID]; ok {
				room.Stop() // Gracefully stop the room's goroutine
				delete(hub.Rooms, roomID)
				log.Printf("Room %s deleted", roomID)
			}
		}
	}
}
