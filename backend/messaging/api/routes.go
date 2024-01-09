package api

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/messaging/api/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", handlers.Home)
	r.GET("/rooms", handlers.GetActiveRooms)
	r.POST("/rooms", handlers.CreateRoom)
	r.GET("/broadcast/:roomID/:userID", handlers.WebSocket)
}
