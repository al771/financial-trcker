package handlers

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)
type CreateTransactionRequest struct {
	Type        string    `json:"type" binding:"required,oneof=income expense"`
	Description string    `json:"description" binding:"min=1,max=150"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Date        time.Time `json:"date" binding:"required"`
	CategoryID  *uint     `json:"category_id"`
}

func CreateTransaction(c *gin.Context) {


	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)
	var request CreateTransactionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// ДУБЛИРОВАНИЕ!!
	var category models.Category
	var hasCategory bool

	if request.CategoryID != nil {
		result := database.DB.Where("id = ? AND user_id = ?", *request.CategoryID, userID).First(&category)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}
		hasCategory = true // запомнили что категория есть
	}

	transaction := models.Transaction{
		Amount:      request.Amount,
		Type:        request.Type,
		Description: request.Description,
		Date:        request.Date,
		UserID:      userID,
		CategoryID:  request.CategoryID,
	}
	if err := database.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create transaction"})
		return
	}

	response := gin.H{
		"id":          transaction.ID,
		"type":        transaction.Type,
		"description": transaction.Description,
		"date":        transaction.Date,
	}
	if hasCategory {
		response["category"] = gin.H{"name": category.Name}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction successfully created",
		"transaction": response,
	})
}



func GetTransactions(c *gin.Context){
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	var transactions []models.Transaction
	result := result := database.DB.Where("user_id = ? OR user_id IS NULL", userID).
		Order("name ASC").
		Find(&categories)

}