package middleware

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET")

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Fail to return token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userIDFloat := claims["user_id"].(float64) // JWT возвращает float64
		userID := uint(userIDFloat)
		var user models.User
		result := database.DB.Where("id = ?", userID).First(&user)
		if result.Error != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}

}
