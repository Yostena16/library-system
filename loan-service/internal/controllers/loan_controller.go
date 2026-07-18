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

// ReturnBook marks a loan as returned and creates a fine if it's overdue.
func ReturnBook(c *gin.Context) {
	loanID := c.Param("id")           // the loan id from the URL
	memberID, _ := c.Get("member_id") // who's returning (from the token)

	// Find the loan
	var loan models.Loan
	if err := database.DB.First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	// Security: you can only return YOUR OWN loan
	if loan.MemberID != memberID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "this is not your loan"})
		return
	}

	// Can't return something already returned
	if loan.Status == "returned" {
		c.JSON(http.StatusConflict, gin.H{"error": "book already returned"})
		return
	}

	// Mark it returned
	now := time.Now()
	loan.ReturnedAt = &now
	loan.Status = "returned"
	database.DB.Save(&loan)

	// If overdue → create a fine ($1 per day late)
	var fine *models.Fine
	if now.After(loan.DueDate) {
		daysLate := int(now.Sub(loan.DueDate).Hours() / 24)
		newFine := models.Fine{
			LoanID:   loan.ID,
			MemberID: loan.MemberID,
			Amount:   float64(daysLate) * 1.0,
		}
		database.DB.Create(&newFine)
		fine = &newFine
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "book returned successfully",
		"loan":    loan,
		"fine":    fine, // null if returned on time
	})
}
