package handlers

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(500, gin.H{"error": "User ID not found"})
		return
	}

	userID := userIDInterface.(uint)
	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Profile successfully found",
		"name":    user.Username,
		"email":   user.Email,
	})
}
