package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents the database schema for a user
type User struct {
	gorm.Model
	Name         string    `json:"name" binding:"required"`
	Age          int       `json:"age" binding:"required"`
	Email        string    `json:"email" binding:"required,email" gorm:"unique;index"`
	Password     string    `json:"password" binding:"required,min=6"`
	DOB          time.Time `json:"dob" binding:"required"`
	RefreshToken string    `json:"-" gorm:"index"` // Store refresh token for rotation
}
