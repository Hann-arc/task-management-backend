package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Hann-arc/task-management-backend/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnnDB() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to load .env")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	DB = db

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Board{},
		&models.Task{},
		&models.TaskLabel{},
		&models.Comment{},
		&models.Attachment{},
		&models.Notification{},
		&models.ActivityLog{},
		&models.Invitation{},
	)
}
