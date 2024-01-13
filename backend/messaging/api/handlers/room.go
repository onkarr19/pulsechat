package handlers

import (
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
	var room models.Room
	if err := c.BindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.ID = uuid.New().String()
	room.CreatedAt = time.Now()
	room.Participants = []models.User{}

	rooms[room.ID] = &room

	if err := createRoomInRedis(room.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a room"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// Get all rooms
func GetRooms(c *gin.Context) {
	// Get all active rooms from the "rooms" set in Redis
	roomIDs, err := redis.Strings(redisPool.Get().Do("SMEMBERS", "rooms"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active rooms"})
		return
	}

	// Fetch details for each active room
	var activeRooms []models.Room
	for _, roomID := range roomIDs {
		// Fetch room details from Redis or your storage
		// Adjust this based on your requirements
		participants, err := redis.Strings(redisPool.Get().Do("SMEMBERS", "room:"+roomID+":participants"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room participants"})
			return
		}

		room := models.Room{
			ID:           roomID,
			Participants: make([]models.User, len(participants)),
		}

		// Fetch user details for each participant
		for i, participantID := range participants {
			// Fetch user details from Redis or your storage
			// Adjust this based on your requirements
			user := models.User{
				ID: participantID,
				// Fetch other user details as needed
			}
			room.Participants[i] = user
		}

		activeRooms = append(activeRooms, room)
	}

	c.JSON(http.StatusOK, activeRooms)
}

func JoinRoom(c *gin.Context) {

	roomID := c.Param("id")

	var requestBody struct {
		UserID string `json:"userID"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := requestBody.UserID

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

	_, err := redisPool.Get().Do("SADD", "rooms", roomID)
	if err != nil {
		log.Printf("Error setting room in Redis: %v", err)
		return err
	}

	return nil
}

func addUserToRoom(roomID, userID string) bool {
	log.Println("Adding user to room")
	conn := redisPool.Get()
	defer conn.Close()
	userCount, err := redis.Int(conn.Do("SADD", roomID, userID))
	if err != nil {
		return false
	}
	return userCount <= 100
}

func removeUserFromRoom(roomID, userID string) {
	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("SREM", roomID, userID)
}

func GetParticipants(c *gin.Context) {
	roomID := c.Param("roomID")
	mu.RLock()
	defer mu.RUnlock()
	if !existsRoom(roomID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room does not exist"})
		return
	}
	participants, err := getRoomParticipants(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve participants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

func getRoomParticipants(roomID string) ([]string, error) {
	conn := redisPool.Get()
	defer conn.Close()
	participants, err := redis.Strings(conn.Do("SMEMBERS", roomID))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return participants, nil
}
