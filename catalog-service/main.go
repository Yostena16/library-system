package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.Default() creates a router with built-in logging + crash recovery
	router := gin.Default()

	// A "health check" route: GET /health tells us the service is alive
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "catalog-service",
		})
	})

	// Start the server on port 8081
	router.Run(":8081")
}
