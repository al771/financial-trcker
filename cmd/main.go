package main

import (
	"financial-tracker/internal/database"
	"financial-tracker/internal/handlers"
	"financial-tracker/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found in .env")
	}
	fmt.Println("JWT Secret loaded:", jwtSecret)
	err = database.Connect()
	if err != nil {
		panic("Failed to connect database" + err.Error())
	}
	r := gin.Default()

	// Незащищённые маршруты
	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)

	// Защищённые маршруты
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Тестовый маршрут для проверки
		protected.GET("/profile", handlers.GetProfile)

		protected.POST("/categories", handlers.CreateCategory)
		protected.GET("/categories", handlers.GetCategories)
		protected.PUT("/categories/:id", handlers.UpdateCategory)
		protected.DELETE("categories/:id", handlers.DeleteCategory)

		protected.POST("/transactions", handlers.CreateTransaction)
	}

	r.Run(":8080")
}
