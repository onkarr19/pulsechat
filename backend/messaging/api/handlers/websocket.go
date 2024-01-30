package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/onkarr19/pulsechat/messaging/internal/models"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	broadcast = make(chan models.Message)
)

type Connection struct {
	WebSocket *websocket.Conn
	UserID    string
}

var rooms = make(map[string]*models.Room)
var connections = make(map[*Connection]bool)

// WebSocketHandler handles WebSocket connections
func WebSocket(c *gin.Context) {
	userID := c.Param("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	connection := &Connection{
		WebSocket: conn,
		UserID:    userID,
	}

	connections[connection] = true

	// Listen for messages from the client
	go readMessages(connection)

	// Listen for close events
	<-c.Request.Context().Done()
	delete(connections, connection)
}

func readMessages(c *Connection) {
	for {
		_, msg, err := c.WebSocket.ReadMessage()
		if err != nil {
			break
		}

		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Error decoding message:", err)
			continue
		}

		// TODO: Handle the received message (forward it to RabbitMQ)

		broadcast <- message
	}
}
