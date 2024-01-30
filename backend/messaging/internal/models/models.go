package models

type User struct {
	UserID   string
	Username string
}

// Room represents a chat room
type Room struct {
	ID    string
	Name  string
	Users map[string]bool // Map of user IDs
}

// Message represents a chat message
type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}
