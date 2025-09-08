package handlers

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func GetTransactions(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID := userIDInterface.(uint)

	query := database.DB.Preload("Category").Where("user_id = ?", userID)

	if categoryID := c.Query("category_id"); categoryID != "" {
		categoryIDInt, err := strconv.Atoi(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id parameter"})
			return
		}
		query = query.Where("category_id = ?", categoryIDInt)
	}

	if transactionType := c.Query("type"); transactionType != "" {
		if transactionType != "income" && transactionType != "expense" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
			return
		}
		query = query.Where("type = ?", transactionType)
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		_, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_from format It must be YYYY-MM-DD"})
			return
		}
		query = query.Where("date >= ?", dateFrom)
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		_, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_to format It must be YYYY-MM-DD"})
			return
		}
		query = query.Where("date <= ?", dateTo)
	}

	query = query.Order("date DESC")

	limit := c.DefaultQuery("limit", "20")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		c.JSON(400, gin.H{"error": "Invalid limit parameter. It must be more then 0"})
		return
	}
	query = query.Limit(limitInt)

	var transactions []models.Transaction
	query.Find(&transactions)

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
