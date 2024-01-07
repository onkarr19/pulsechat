package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", homeHandler)
}

func homeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to PulseChat Auth Service!",
	})
}
