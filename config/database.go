package config

import (
	"fmt"
	"log"
	"os"

	"golang-auth-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectToDB initializes the connection to PostgreSQL using GORM
func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Database connection established")

	// AutoMigrate will create/update the table based on the User model
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
