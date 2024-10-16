package ws

// Room represents a chat room where clients can communicate
type Room struct {
	AllowedClients map[string]bool // List of allowed client IDs
	Clients        map[*Client]bool
	Broadcast      chan []byte
	Register       chan *Client
	Unregister     chan *Client
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
