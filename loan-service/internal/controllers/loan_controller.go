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

// BorrowBook godoc
// @Summary      Borrow a book
// @Description  A logged-in member borrows a book (availability checked with Catalog)
// @Tags         loans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body  BorrowInput  true  "Book to borrow"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /loans [post]
func BorrowBook(c *gin.Context) {
	var input BorrowInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// ReturnBook godoc
// @Summary      Return a book (librarian only)
// @Description  A librarian processes a return; creates a fine if overdue
// @Tags         loans
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Loan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /loans/{id}/return [post]
func ReturnBook(c *gin.Context) {
	loanID := c.Param("id")

	// Only librarians can process returns
	role, _ := c.Get("role")
	if role != "librarian" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only librarians can process returns"})
		return
	}

	var loan models.Loan
	if err := database.DB.First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.Status == "returned" {
		c.JSON(http.StatusConflict, gin.H{"error": "book already returned"})
		return
	}

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
		"fine":    fine,
	})
}

// GetMyLoans godoc
// @Summary      List my loans
// @Description  List all loans belonging to the logged-in member
// @Tags         loans
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /loans [get]
func GetMyLoans(c *gin.Context) {
	memberID, _ := c.Get("member_id")

	var loans []models.Loan
	database.DB.Where("member_id = ?", memberID).Order("id desc").Find(&loans)

	c.JSON(http.StatusOK, gin.H{"loans": loans})
}
