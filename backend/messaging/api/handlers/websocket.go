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
)

type Connection struct {
	WebSocket *websocket.Conn
	UserID    string
}

var rooms = make(map[string]*models.Room)

// WebSocketHandler handles WebSocket connections
func WebSocket(c *gin.Context) {
	room := c.Param("room")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Create or retrieve the RoomProcess for the given room
	roomMutex.Lock()
	rp, ok := RoomProcesses[room]
	if !ok {
		rp = NewRoomProcess(room)
		RoomProcesses[room] = rp
		rp.Start()
	}
	rp.AddConnection(conn)
	roomMutex.Unlock()

	// Close the connection when the function returns
	defer func() {
		roomMutex.Lock()
		rp.RemoveConnection(conn)
		roomMutex.Unlock()
	}()

	// Read and broadcast messages
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Publish the message to the RabbitMQ fanout exchange
		PublishToRabbitMQ(room, p)
	}
}

func ReadMessages(c *Connection) {
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

		// broadcast <- message
	}
}
