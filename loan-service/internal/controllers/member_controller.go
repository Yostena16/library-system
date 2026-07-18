package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"loan-service/internal/database"
	"loan-service/internal/models"
)

// GetMyProfile returns the logged-in member's own profile.
func GetMyProfile(c *gin.Context) {
	// The middleware stored member_id in the context
	memberID, _ := c.Get("member_id")

	var member models.Member
	if err := database.DB.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member": member})
}
