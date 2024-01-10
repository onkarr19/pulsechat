package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/onkarr19/pulsechat/messaging/internal/models"
)

var (
	redisPool *redis.Pool
	mu        sync.RWMutex
)

func init() {
	redisPool = &redis.Pool{
		MaxIdle:   10,
		MaxActive: 30,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

// Create a new room
func CreateRoom(c *gin.Context) {
	var requestData struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := uuid.New().String()

	newRoom := &models.Room{
		ID:    roomID,
		Name:  requestData.Name,
		Users: make(map[string]bool),
	}

	rooms[roomID] = newRoom

	createRoomInRedis(roomID)
	responseData := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		ID:   newRoom.ID,
		Name: newRoom.Name,
	}

	c.JSON(http.StatusOK, responseData)
}

// Get active rooms (already created)
func GetActiveRooms(c *gin.Context) {
	type responseData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	activeRooms := make([]responseData, 0, len(rooms))

	for _, room := range rooms {
		activeRooms = append(activeRooms, responseData{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	c.JSON(http.StatusOK, activeRooms)
}

// func createRoom(c *gin.Context) {
// 	roomID := c.PostForm("roomID")
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if existsRoom(roomID) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Room already exists"})
// 		return
// 	}
// 	createRoomInRedis(roomID)
// 	c.JSON(http.StatusOK, gin.H{"message": "Room created successfully"})
// }

func JoinRoom(c *gin.Context) {
	roomID := c.PostForm("roomID")
	userID := c.PostForm("userID")
	mu.Lock()
	defer mu.Unlock()
	if !ExistsRoom(roomID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room does not exist"})
		return
	}
	if !addUserToRoom(roomID, userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room is full"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User joined room successfully"})
}

func LeaveRoom(c *gin.Context) {
	roomID := c.PostForm("roomID")
	userID := c.PostForm("userID")
	mu.Lock()
	defer mu.Unlock()
	if !ExistsRoom(roomID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room does not exist"})
		return
	}
	removeUserFromRoom(roomID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "User left room successfully"})
}

func ExistsRoom(roomID string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", roomID))
	if err != nil {
		log.Println(err)
		return false
	}
	return exists
}

func createRoomInRedis(roomID string) {
	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("SET", roomID, "")
}

func addUserToRoom(roomID, userID string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	userCount, err := redis.Int(conn.Do("SADD", roomID, userID))
	if err != nil {
		log.Println(err)
		return false
	}
	return userCount <= 100
}

func removeUserFromRoom(roomID, userID string) {
	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("SREM", roomID, userID)
}
