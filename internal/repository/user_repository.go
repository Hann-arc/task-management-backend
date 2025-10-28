package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// function to create a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// function to find user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// function to find user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// function to update user by ID
func (r *UserRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.DB.Model(&models.User{}).Where("id = ?", id).Updates(data).Error
}

// function to get all users
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}
