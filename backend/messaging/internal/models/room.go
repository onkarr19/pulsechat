package models

// Room represents a chat room
type Room struct {
	ID    string
	Name  string
	Users map[string]bool // Map of user IDs
}
