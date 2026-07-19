package database

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"loan-service/internal/models"
)

// SeedLibrarian creates the librarian account on startup if it doesn't exist.
func SeedLibrarian() {
	email := os.Getenv("LIBRARIAN_EMAIL")
	password := os.Getenv("LIBRARIAN_PASSWORD")
	if email == "" || password == "" {
		return // no librarian configured — skip
	}

	// Already exists? Don't create a duplicate.
	var existing models.Member
	if err := DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	librarian := models.Member{
		Name:     "Librarian",
		Email:    email,
		Password: string(hashed),
		Role:     "librarian", // ← role assigned at creation
	}
	DB.Create(&librarian)
	log.Println("✅ Librarian account ready:", email)
}
