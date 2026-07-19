package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"catalog-service/internal/database"
	"catalog-service/internal/models"
	"catalog-service/internal/routes"
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

	routes.SetupRoutes(router)

	router.Run(":8081")
}
