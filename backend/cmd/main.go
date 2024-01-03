package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/internal/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	// Start the server
	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
