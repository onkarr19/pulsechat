package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/profile/api"
)

func main() {
	r := gin.Default()

	api.RegisterRoutes(r)

	// Start the service
	err := r.Run("localhost:8085")
	if err != nil {
		panic(err)
	}
}
