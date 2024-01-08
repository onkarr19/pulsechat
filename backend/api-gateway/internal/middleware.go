package internal

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ForwardRequest is a middleware to forward requests to microservices
func ForwardRequest(sd ServiceDiscovery, serviceName string) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Forward the request to the respective microservice
		serviceURL := sd.GetServiceURL(serviceName)

		// Extract the path after the service prefix, e.g., /messaging/ws -> /ws
		originalPath := c.Param("any")
		forwardedPath := strings.TrimPrefix(originalPath, "/"+serviceName)

		// Create a new request with the same method, URL, and body as the original request
		forwardedRequest, err := http.NewRequest(c.Request.Method, serviceURL+forwardedPath, c.Request.Body)
		if err != nil {
			// Handle error, e.g., log it and return an error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Copy headers from the original request to the forwarded request
		forwardedRequest.Header = c.Request.Header

		// Use the HTTP client to send the forwarded request
		client := &http.Client{}
		response, err := client.Do(forwardedRequest)
		if err != nil {
			// Handle error, e.g., log it and return an error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer response.Body.Close()

		// Copy the response status code and headers back to the original response
		c.Status(response.StatusCode)
		for key, values := range response.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy the response body back to the original response
		_, err = io.Copy(c.Writer, response.Body)
		if err != nil {
			// Handle error, e.g., log it and return an error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
}
