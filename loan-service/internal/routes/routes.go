package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"loan-service/internal/controllers"
)

// SetupRoutes connects URLs to controller functions.
func SetupRoutes(router *gin.Engine) {
	router.GET("/book", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "loan-service"})
	})

	// Group all auth routes under /auth
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}
