package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := uuid.New().String()

	newRoom := &models.Room{
		ID:        roomID,
		Name:      requestBody.Name,
		CreatedAt: time.Now(),
		Users:     make(map[string]bool),
	}

	rooms[roomID] = newRoom

	if err := createRoomInRedis(roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a room"})
		return
	}
	responseData := struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created-at"`
	}{
		ID:        newRoom.ID,
		Name:      newRoom.Name,
		CreatedAt: newRoom.CreatedAt,
	}

	c.JSON(http.StatusOK, responseData)
}

// Get active rooms
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

func JoinRoom(c *gin.Context) {
	roomID := c.PostForm("roomID")
	userID := c.PostForm("userID")
	mu.Lock()
	defer mu.Unlock()
	if !existsRoom(roomID) {
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
	if !existsRoom(roomID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room does not exist"})
		return
	}
	removeUserFromRoom(roomID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "User left room successfully"})
}

func existsRoom(roomID string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", roomID))
	if err != nil {
		log.Println(err)
		return false
	}
	return exists
}

func createRoomInRedis(roomID string) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", roomID, "")
	if err != nil {
		log.Printf("Error setting room in Redis: %v", err)
		return err
	}

	fmt.Println("reaching in redis")
	return nil
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
