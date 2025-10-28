package services

import (
	"errors"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"gorm.io/gorm"

	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if req.Name == nil && req.AvatarUrl == nil {
		return nil, apperrors.ErrInvalidUserData
	}

	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare update data
	data := map[string]interface{}{}

	if req.Name != nil && *req.Name != "" {
		data["name"] = *req.Name
	}

	if req.AvatarUrl != nil {
		if user.AvatarUrl != "" {
			utils.DeleteFromCloudinary(user.AvatarUrl)
		}
		data["avatar_url"] = *req.AvatarUrl
	}

	if err := s.userRepo.Update(id, data); err != nil {
		return nil, err
	}

	// Get updated user
	updatedUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        updatedUser.ID.String(),
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		AvatarUrl: updatedUser.AvatarUrl,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, dto.UserResponse{
			ID:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			AvatarUrl: u.AvatarUrl,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return result, nil
}
