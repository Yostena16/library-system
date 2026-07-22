package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/internal/routes"

	_ "auth-service/docs"
)

// @title           Auth Service API
// @version         1.0
// @description     Library system — Auth microservice (registration, login, member profile).
// @host            localhost:8083
// @BasePath        /
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	database.Connect()

	if err := database.DB.AutoMigrate(&models.Member{}); err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	log.Println("✅ Database migrated")

	database.SeedLibrarian()

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8083")
}
