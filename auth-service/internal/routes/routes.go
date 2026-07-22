package routes

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"auth-service/internal/controllers"
	"auth-service/internal/middleware"
)

// SetupRoutes connects URLs to controller functions.
func SetupRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public — anyone can register or log in
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Requires a valid token
	members := router.Group("/members")
	members.Use(middleware.AuthMiddleware())
	{
		members.GET("/me", controllers.GetMyProfile)
	}
}
