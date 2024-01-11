package api

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/messaging/api/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", handlers.Home)
	r.GET("/active-rooms", handlers.GetActiveRooms)
	r.POST("/create-room", handlers.CreateRoom)
	r.POST("/join-room", handlers.CreateRoom)
	r.POST("/leave-room", handlers.CreateRoom)
	r.GET("/broadcast/:roomID/:userID", handlers.WebSocket)
}
