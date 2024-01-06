package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Start the service
	err := r.Run("localhost:8081")
	if err != nil {
		panic(err)
	}
}
