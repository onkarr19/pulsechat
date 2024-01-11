package models

import "time"

// Room represents a chat room
type Room struct {
	ID        string
	Name      string
	CreatedAt time.Time
	Users     map[string]bool // Map of user IDs
}
