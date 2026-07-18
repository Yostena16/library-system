package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"loan-service/internal/database"
)

func main() {
	// Load the .env file into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	// Connect to PostgreSQL
	database.Connect()

	router := gin.Default()

	router.GET("/book", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "loan-service",
		})
	})

	router.Run(":8082")
}
