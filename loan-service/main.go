package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"loan-service/internal/database"
	"loan-service/internal/models"
	"loan-service/internal/routes"

	_ "loan-service/docs"
)

// @title           Loan Service API
// @version         1.0
// @description     Library system — Loan microservice (borrowing, returns, fines).
// @host            localhost:8082
// @BasePath        /
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	database.Connect()

	if err := database.DB.AutoMigrate(&models.Loan{}, &models.Fine{}); err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	log.Println("✅ Database migrated")

	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":8082")
}
