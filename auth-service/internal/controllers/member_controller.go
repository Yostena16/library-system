package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-service/internal/database"
	"auth-service/internal/models"
)

// GetMyProfile godoc
// @Summary      Get my profile
// @Description  Returns the logged-in member's profile
// @Tags         members
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /members/me [get]
func GetMyProfile(c *gin.Context) {
	memberID, _ := c.Get("member_id")

	var member models.Member
	if err := database.DB.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member": member})
}
