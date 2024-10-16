package ws

import (
	"demo/db" // Assuming this handles your database interactions
	"log"
	"time"
)

type MessageManager struct {
	roomID        string
	latestMessage []byte
	dbInterval    time.Duration
	stop          chan bool
}

// NewMessageManager creates a new manager that tracks broadcast messages
func NewMessageManager(roomID string, dbInterval time.Duration) *MessageManager {
	return &MessageManager{
		roomID:     roomID,
		dbInterval: dbInterval,
		stop:       make(chan bool),
	}
}

// UpdateMessage updates the latest broadcast message
func (m *MessageManager) UpdateMessage(message []byte) {
	m.latestMessage = message
}

// Run starts a background process that periodically saves the message to the database
func (m *MessageManager) Run() {
	ticker := time.NewTicker(m.dbInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if m.latestMessage != nil {
				// Save the latest message to the database

				log.Printf("Saving message to DB")
				err := db.SetRedis(m.roomID, string(m.latestMessage))
				if err != nil {
					log.Printf("Error saving message to DB for room %s: %v", m.roomID, err)
				} else {
					log.Printf("Saved message to DB for room %s", m.roomID)
				}
				m.latestMessage = nil // Reset after saving
			}
		case <-m.stop:
			// Stop saving and clean up
			return
		}
	}
}

// Stop stops the message manager
func (m *MessageManager) Stop() {
	m.stop <- true
}
