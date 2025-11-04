package services

import (
	"errors"

	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthservice creates a new instance of AuthService
func NewAuthservice(repo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

// Register a new user
func (s *AuthService) Register(name, email, password string) error {
	hash, _ := utils.HashPassword(password)
	user := &models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
	}
	return s.userRepo.Create(user)

}

// Login an existing user
func (s *AuthService) Login(email, password string) (string, error) {

	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("Email not found")
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("Invalid password")
	}

	// Generate JWT token
	return utils.GenerateToken(user.ID)
}
