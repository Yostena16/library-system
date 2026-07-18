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
	// ⚠️ TEMPORARY STUB — your friend replaces this with real DB logic later.
	// Answers whether a given book can be borrowed.
	router.GET("/books/:id/availability", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"book_id":   id,
			"available": true,
			"copies":    3,
		})
	})

	// Start the server on port 8081
	router.Run(":8081")
}
