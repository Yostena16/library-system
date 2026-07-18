package controllers // ← matches the folder name

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"loan-service/internal/database"
	"loan-service/internal/models"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register creates a new member account.
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. Read + validate the JSON body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the password (so we never store it in plain text)
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// 3. Build the member
	member := models.Member{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	// 4. Save to the database
	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "member registered successfully",
		"member":  member,
	})
}
