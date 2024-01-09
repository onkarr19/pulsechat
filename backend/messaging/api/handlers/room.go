package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/onkarr19/pulsechat/messaging/internal/models"
)

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
