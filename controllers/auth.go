package controllers

import (
	"net/http"
	"os"
	"time"

	"golang-auth-api/config"
	"golang-auth-api/models"
	"golang-auth-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Register handles user signup
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Age      int    `json:"age" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		DOB      string `json:"dob" binding:"required"` // Format: YYYY-MM-DD
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse DOB
	dob, err := time.Parse("2006-01-02", input.DOB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DOB format, use YYYY-MM-DD"})
		return
	}

	// Hash Password
	hashedPassword, _ := utils.HashPassword(input.Password)

	user := models.User{
		Name:     input.Name,
		Age:      input.Age,
		Email:    input.Email,
		Password: hashedPassword,
		DOB:      dob,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user (Email might already exist)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate Tokens
	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	// Store Refresh Token in DB for rotation
	config.DB.Model(&user).Update("refresh_token", refreshToken)

	// Set Refresh Token in HttpOnly Cookie
	// SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true) // Secure=true on Prod

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"message":      "Logged in successfully",
	})
}

// Refresh handles token rotation
func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	// Verify Refresh Token
	token, err := utils.ValidateToken(refreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check if token matches stored token (Rotation check)
	if user.RefreshToken != refreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token reuse detected or invalid"})
		return
	}

	// Issue NEW tokens
	newAccessToken, _ := utils.GenerateAccessToken(user.ID)
	newRefreshToken, _ := utils.GenerateRefreshToken(user.ID)

	// Update DB with NEW refresh token
	config.DB.Model(&user).Update("refresh_token", newRefreshToken)

	// Update Cookie
	c.SetCookie("refresh_token", newRefreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// Logout clears tokens
func Logout(c *gin.Context) {
	// Find user from context (set by middleware)
	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	// Nullify refresh token in DB
	config.DB.Model(&user).Update("refresh_token", "")

	// Clear Cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
