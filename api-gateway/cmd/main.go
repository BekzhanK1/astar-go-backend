package main

import (
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/services"

	// "api-gateway/internal/utils"
	"log"
	// "net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	// config := utils.LoadConfig()
	authService, err := services.NewAuthServiceClient("localhost:50052")
	if err != nil {
		log.Fatalf("Failed to create auth service client: %v", err)
	}

	userService, err := services.NewUserServiceClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	authHandlers := handlers.NewAuthHandler(authService)
	userHandlers := handlers.NewUserHandler(userService)

	// Initialize Gin router
	r := gin.Default()

	// Public routes
	r.POST("/auth/login", authHandlers.Login)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthMiddleware())
	// {
	protected.GET("/user/profile", userHandlers.UserProfile)
	// 	protected.GET("/school/data", handlers.SchoolData)
	// }

	// Start the server
	log.Printf("API Gateway running on port %s", "8000")
	if err := r.Run(":" + "8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
