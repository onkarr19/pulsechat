package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", homeHandler)
	r.POST("/ws", handleMessage)

}

func homeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to PulseChat Messaging Service!",
	})
}
func handleMessage(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
	})
}
