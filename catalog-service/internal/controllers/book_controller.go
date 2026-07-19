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

// BookInput is the request body for creating or updating a book.
type BookInput struct {
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author"`
	Category        string `json:"category"`
	ISBN            string `json:"isbn"`
	TotalCopies     int    `json:"total_copies"`
	AvailableCopies int    `json:"available_copies"`
}

// CreateBook godoc
// @Summary      Add a book (librarian only)
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body  BookInput  true  "Book to add"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /books [post]
func CreateBook(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only librarians can add books"})
		return
	}

	var input BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		Title:           input.Title,
		Author:          input.Author,
		Category:        input.Category,
		ISBN:            input.ISBN,
		TotalCopies:     input.TotalCopies,
		AvailableCopies: input.AvailableCopies,
	}
	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book added successfully", "book": book})
}

// UpdateBook godoc
// @Summary      Edit a book (librarian only)
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path  int        true  "Book ID"
// @Param        input  body  BookInput  true  "Updated book fields"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [put]
func UpdateBook(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only librarians can edit books"})
		return
	}

	id := c.Param("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var input BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.Title = input.Title
	book.Author = input.Author
	book.Category = input.Category
	book.ISBN = input.ISBN
	book.TotalCopies = input.TotalCopies
	book.AvailableCopies = input.AvailableCopies
	database.DB.Save(&book)

	c.JSON(http.StatusOK, gin.H{"message": "book updated successfully", "book": book})
}

// DeleteBook godoc
// @Summary      Remove a book (librarian only)
// @Tags         books
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Book ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only librarians can remove books"})
		return
	}

	id := c.Param("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	database.DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "book removed successfully"})
}
