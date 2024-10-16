package ws

import "log"

// Hub manages multiple rooms
type Hub struct {
	Rooms map[string]*Room
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
