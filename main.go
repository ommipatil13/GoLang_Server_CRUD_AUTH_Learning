package main

import (
	"fmt"
	"log"
	"os"

	"golang-auth-api/config"
	"golang-auth-api/controllers"
	"golang-auth-api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// init() runs before main()
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Database
	config.ConnectToDB()
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Default Route (Welcome)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go Auth API 🚀",
			"status": "Healthy",
		})
	})

	// PUBLIC ROUTES (Auth)
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.Refresh)
	}

	// PROTECTED ROUTES (Requires JWT)
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", controllers.GetProfile)
		
		// CRUD Operations for Users
		protected.PUT("/user/update", controllers.UpdateUser)
		protected.DELETE("/user/delete", controllers.DeleteUser)
		protected.GET("/users", controllers.GetAllUsers) // List all users

		// Logout
		protected.POST("/logout", controllers.Logout)
	}

	// Determine Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("✅ Server starting and running on port %s...\n", port)
	r.Run(":" + port)
}
