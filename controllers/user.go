package controllers

import (
	"net/http"
	"time"

	"golang-auth-api/config"
	"golang-auth-api/models"
	"golang-auth-api/utils"
	"github.com/gin-gonic/gin"
)

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles user profile updates
func UpdateUser(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	var input struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Password string `json:"password"`
		DOB      string `json:"dob"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Age != 0 {
		user.Age = input.Age
	}
	if input.Password != "" && len(input.Password) >= 6 {
		user.Password, _ = utils.HashPassword(input.Password)
	}
	if input.DOB != "" {
		dob, err := time.Parse("2006-01-02", input.DOB)
		if err == nil {
			user.DOB = dob
		}
	}

	// Save changes
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

// DeleteUser deletes the current user's account
func DeleteUser(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	// Soft delete (gorm.Model default behavior)
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Clear Cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsers is an optional helper to list all users in the DB
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
