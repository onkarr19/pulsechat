// messaging/cmd/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/ws", handleMessage)

	// Start the service
	err := r.Run(":8082")
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}

func handleMessage(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
	})
}
