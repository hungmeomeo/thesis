package ws

import (
	//"demo/db"

	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection
type Client struct {
	Conn   *websocket.Conn
	Hub    *Hub
	RoomID string
	Send   chan []byte
}

// Message struct for broadcasting
type Message struct {
	ClientID string `json:"clientID"`
	Content  string `json:"content"`
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
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}

		// Unmarshal the message directly here
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue // Skip logging if unmarshalling fails
		}

		// Log the content correctly
		log.Printf("Sent message to client %s, code is %s", c.RoomID, msg.Content)
		//codeMap.AddCode(msg.Content, c.RoomID)
	}
}
