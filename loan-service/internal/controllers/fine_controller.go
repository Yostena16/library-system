package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"loan-service/internal/database"
	"loan-service/internal/models"
)

// GetMyFines lists all fines belonging to the logged-in member.
func GetMyFines(c *gin.Context) {
	memberID, _ := c.Get("member_id")

	var fines []models.Fine
	database.DB.Where("member_id = ?", memberID).Order("id desc").Find(&fines)

	c.JSON(http.StatusOK, gin.H{"fines": fines})
}
