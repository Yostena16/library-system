package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"loan-service/internal/controllers"
	"loan-service/internal/middleware"
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

	members := router.Group("/members")
	members.Use(middleware.AuthMiddleware())
	{
		members.GET("/me", controllers.GetMyProfile)
	}

	//  Loans — require login
	loans := router.Group("/loans")
	loans.Use(middleware.AuthMiddleware())
	{
		loans.POST("", controllers.BorrowBook) // POST /loans
		loans.POST("/:id/return", controllers.ReturnBook)
		loans.GET("", controllers.GetMyLoans)
	}

	fines := router.Group("/fines")
	fines.Use(middleware.AuthMiddleware())
	{
		fines.GET("", controllers.GetMyFines)
	}
}
