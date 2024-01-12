package api

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/messaging/api/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", handlers.Home)
	r.GET("/rooms", handlers.GetRooms)
	r.POST("/rooms", handlers.CreateRoom)
	r.PUT("/rooms/:id/join", handlers.JoinRoom)
	r.DELETE("/rooms/:id/leave", handlers.LeaveRoom)
	r.GET("/rooms/:id", handlers.GetParticipants)

	// r.GET("/rooms/:id/messages", handlers.GetMessages)
	// r.GET("/messages/:id", getPendingMessages)
	// r.POST("/messages", sendMessage)
}
