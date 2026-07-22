package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"loan-service/internal/controllers"
	"loan-service/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes connects URLs to controller functions.
func SetupRoutes(router *gin.Engine) {
	router.GET("/book", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "loan-service"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Loans — require login
	loans := router.Group("/loans")
	loans.Use(middleware.AuthMiddleware())
	{
		loans.POST("", controllers.BorrowBook)
		loans.POST("/:id/return", controllers.ReturnBook)
		loans.GET("", controllers.GetMyLoans)
	}

	fines := router.Group("/fines")
	fines.Use(middleware.AuthMiddleware())
	{
		fines.GET("", controllers.GetMyFines)
	}
}
