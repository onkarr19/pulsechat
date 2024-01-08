package main

import (
	"github.com/gin-gonic/gin"
	"github.com/onkarr19/pulsechat/api-gateway/internal"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize service discovery
	sd := internal.NewServiceDiscovery()

	// Define routes for microservices
	r.Any("/auth/*any", internal.ForwardRequest(sd, "auth"))
	r.Any("/profile/*any", internal.ForwardRequest(sd, "profile"))
	r.Any("/messaging/*any", internal.ForwardRequest(sd, "messaging"))
	r.Any("/notification/*any", internal.ForwardRequest(sd, "notification"))
	r.Any("/storage/*any", internal.ForwardRequest(sd, "storage"))

	// Run the API gateway
	r.Run(":8080")
}
