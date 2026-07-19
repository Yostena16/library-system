package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"catalog-service/internal/database"
	"catalog-service/internal/models"
)

func main() {
	// Load the .env file into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	// Connect to PostgreSQL
	database.Connect()

	// Auto-migrate: make the tables match our models
	if err := database.DB.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	log.Println("✅ Database migrated")

	router := gin.Default()

	// A "health check" route: GET /health tells us the service is alive
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "catalog-service",
		})
	})

	// Routes (public book browsing + librarian management) get wired up
	// via routes.SetupRoutes() once that package exists — for now this
	// keeps the service runnable at each step.

	router.Run(":8081")
}
