package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/messaging/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	// Start the service
	err := r.Run("localhost:8082")
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
