package models

import "time"

// Room represents a chat room
type Room struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreatedAt    time.Time
	CreatorID    string `json:"creator_id"`
	Participants []User `json:"participants"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	LastSeen  time.Time `json:"last_seen"`
	Connected bool      `json:"connected"`
}
