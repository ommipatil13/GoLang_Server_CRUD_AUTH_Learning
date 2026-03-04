package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang-auth-api/config"
	"golang-auth-api/controllers"
	"golang-auth-api/middlewares"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// Only load .env if not running in Lambda
	if os.Getenv("LAMBDA_TASK_ROOT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	// Connect to Database
	config.ConnectToDB()
}

func setupRouter() *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	// Default Route (Welcome)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go Auth API 🚀",
			"status":  "Healthy",
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
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/user/update", controllers.UpdateUser)
		protected.DELETE("/user/delete", controllers.DeleteUser)
		protected.GET("/users", controllers.GetAllUsers)
		protected.POST("/logout", controllers.Logout)
	}

	return r
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// Create the adapter only once
		ginLambda = ginadapter.New(setupRouter())
	}
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		// Running on AWS Lambda
		lambda.Start(Handler)
	} else {
		// Running locally
		r := setupRouter()
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		fmt.Printf("✅ Server starting and running on port %s...\n", port)
		r.Run(":" + port)
	}
}
