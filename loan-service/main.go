package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Health check — confirms the loan service is alive
	router.GET("/book", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "loan-service",
		})
	})

	// Loan service runs on port 8082
	router.Run(":8082")
}
