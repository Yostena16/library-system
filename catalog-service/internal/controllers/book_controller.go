package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"catalog-service/internal/database"
	"catalog-service/internal/models"
)

// GetBooks godoc
// @Summary      List all books
// @Description  Public — anyone can browse the catalog
// @Tags         books
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /books [get]
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"books": books})
}

// GetBook godoc
// @Summary      Get one book
// @Description  Public — fetch a single book by ID
// @Tags         books
// @Produce      json
// @Param        id   path  int  true  "Book ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [get]
func GetBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book})
}

// GetAvailability godoc
// @Summary      Check book availability
// @Description  Public — used by loan-service before allowing a borrow
// @Tags         books
// @Produce      json
// @Param        id   path  int  true  "Book ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id}/availability [get]
func GetAvailability(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// ⚠️ Contract with loan-service: field names "available" and "copies"
	// must stay exactly this — loan-service/internal/clients/catalog.go
	// decodes the response into a struct with these exact json tags.
	c.JSON(http.StatusOK, gin.H{
		"book_id":   book.ID,
		"available": book.AvailableCopies > 0,
		"copies":    book.AvailableCopies,
	})
}
