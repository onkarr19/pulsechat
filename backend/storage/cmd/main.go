package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/storage/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	// Start the service
	err := r.Run("localhost:8084")
	if err != nil {
		panic(err)
	}
}
