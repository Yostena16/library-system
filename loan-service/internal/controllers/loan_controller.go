package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"loan-service/internal/clients"
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

	// Ask the Catalog service whether the book is available
	availability, err := clients.CheckAvailability(input.BookID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "catalog service unavailable"})
		return
	}
	if !availability.Available {
		c.JSON(http.StatusConflict, gin.H{"error": "book is not available"})
		return
	}

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
