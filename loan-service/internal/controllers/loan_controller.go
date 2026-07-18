package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"loan-service/internal/database"
	"loan-service/internal/models"
)

// BorrowInput is the JSON we expect when borrowing.
type BorrowInput struct {
	BookID uint `json:"book_id" binding:"required"`
}

// BorrowBook creates a loan for the logged-in member.
func BorrowBook(c *gin.Context) {
	var input BorrowInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Who is borrowing? The middleware stored member_id in the context.
	memberID, _ := c.Get("member_id")

	// TODO (next step): ask the Catalog service if this book is available.

	// Create the loan: borrowed now, due in 14 days.
	loan := models.Loan{
		MemberID:   memberID.(uint),
		BookID:     input.BookID,
		BorrowedAt: time.Now(),
		DueDate:    time.Now().Add(14 * 24 * time.Hour),
		Status:     "borrowed",
	}

	if err := database.DB.Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create loan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "book borrowed successfully",
		"loan":    loan,
	})
}
