package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", homeHandler)
	r.POST("/sendMessage", sendMessageHandler)
}

func homeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to PulseChat!",
	})
}

func sendMessageHandler(c *gin.Context) {

}
