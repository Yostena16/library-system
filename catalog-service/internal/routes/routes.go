package routes

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"catalog-service/internal/controllers"
	"catalog-service/internal/middleware"
)

// SetupRoutes connects URLs to controller functions.
func SetupRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public — anyone can browse the catalog, no login required
	books := router.Group("/books")
	{
		books.GET("", controllers.GetBooks)
		books.GET("/:id", controllers.GetBook)
		books.GET("/:id/availability", controllers.GetAvailability)
	}

	// Librarian only — requires a valid JWT; the role=="librarian" check
	// happens inside each handler itself (same pattern as loan-service's
	// ReturnBook)
	librarianBooks := router.Group("/books")
	librarianBooks.Use(middleware.AuthMiddleware())
	{
		librarianBooks.POST("", controllers.CreateBook)
		librarianBooks.PUT("/:id", controllers.UpdateBook)
		librarianBooks.DELETE("/:id", controllers.DeleteBook)
	}
}
