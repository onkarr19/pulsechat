package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/onkarr19/pulsechat/messaging/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var rooms = make(map[string]*models.Room)

var connections = make(map[string]*websocket.Conn)

// WebSocketHandler handles WebSocket connections
func WebSocket(c *gin.Context) {
	roomID := c.Param("roomID")
	userID := c.Param("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Add the user's WebSocket connection to the connections map
	connections[userID] = conn

	// Check if the room exists
	room, exists := rooms[roomID]
	if !exists {
		log.Printf("Room %s does not exist\n", roomID)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Register user to the room
	room.Participants = append(room.Participants, models.User{ID: userID})

	// Handle WebSocket messages
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// Remove the user from the room and connections map when they disconnect
			// delete(room.Participants, userID)
			delete(connections, userID)
			return
		}

		// Handle chat message
		if messageType == websocket.TextMessage {
			message := string(p)

			// Include the user's ID in the message
			messageWithSender := fmt.Sprintf("[%s]: %s", userID, message)

			// Broadcast the message to all users in the room
			for i := range room.Participants {
				recipientConn, found := connections[room.Participants[i].ID]
				if !found {
					// Handle case where recipient is not connected
					log.Printf("User %s is not connected\n", room.Participants[i].ID)
					continue
				}
				if err := recipientConn.WriteMessage(websocket.TextMessage, []byte(messageWithSender)); err != nil {
					// Handle write error (e.g., user disconnected)
					log.Printf("Error sending message to user %s: %v\n", room.Participants[i].ID, err)
				}
			}
		}
	}
}
