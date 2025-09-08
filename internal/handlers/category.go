package handlers

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=200"`
}

func categoryExists(name string, userID uint) bool {
	var existingCategory models.Category
	result := database.DB.Where("name = ? AND (user_id = ? OR user_id IS NULL)", name, userID).First(&existingCategory)
	return result.Error == nil // true если найдена
}

func CreateCategory(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	var request CreateCategoryRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if categoryExists(request.Name, userID) {
		c.JSON(http.StatusConflict, gin.H{"error": "This category already exists"})
		return
	}
	category := models.Category{
		Name:   request.Name,
		UserID: &userID,
	}
	result := database.DB.Create(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	// 5. Возврат результата
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category is created",
		"category": category,
	})
}

func GetCategories(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	var categories []models.Category
	result := database.DB.Where("user_id = ? OR user_id IS NULL", userID).
		Order("name ASC").
		Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "Categories are found",
		"categories": categories,
		"count":      len(categories),
	})
}

func UpdateCategory(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	categoryID := c.Param("id")

	var request UpdateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Category
	result := database.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or access denied"})
		return
	}

	if category.Name != request.Name && categoryExists(request.Name, userID) {
		c.JSON(http.StatusConflict, gin.H{"error": "Category with this name already exists"})
		return
	}

	category.Name = request.Name
	result = database.DB.Save(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category updated successfully",
		"category": category,
	})

}

func DeleteCategory(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	categoryID := c.Param("id")

	var category models.Category
	result := database.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found or access denied"})
		return
	}
	result = database.DB.Delete(&category)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})

}
