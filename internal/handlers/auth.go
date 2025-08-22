package handlers

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(c *gin.Context) {
	var request RegisterRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	var existingUser models.User
	result := database.DB.Where("email = ?", request.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password",
		})
		return
	}
	newUser := models.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}
	result = database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to save new user",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{ // 201 = Created
		"message":  "User created successfully!",
		"user_id":  newUser.ID,
		"username": request.Username,
		"email":    request.Email,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var request LoginRequest
	err := c.ShouldBindJSON(&request)
	//проверка корректности json
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	var user models.User
	result := database.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
	}
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // истекает через 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // выдан сейчас
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"expires": time.Now().Add(24 * time.Hour),
	})

}
