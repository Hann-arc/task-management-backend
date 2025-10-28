package handlers

import (
	"errors"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetProfile retrieves the profile of the authenticated user
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "User not found", "")
		}
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch user", err.Error())
	}

	return utils.Success(c, "User profile fetched successfully", user)
}

// UpdateProfile updates the profile of the authenticated user
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	name := c.FormValue("name")

	// Get avatar file if exists
	avatarFile, err := c.FormFile("avatar")
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid file upload", err.Error())
	}

	req := dto.UpdateUserRequest{
		Name: &name,
	}

	// If avatar file exists, upload to cloudinary
	if avatarFile != nil {
		if avatarFile.Size > 5*1024*1024 {
			return utils.Error(c, fiber.StatusBadRequest, "Avatar file too large (max 5MB)", "")
		}

		// Validate file type
		contentType := avatarFile.Header.Get("Content-Type")
		if !isImageFile(contentType) {
			return utils.Error(c, fiber.StatusBadRequest, "Invalid file type. Only images are allowed", "")
		}

		avatarURL, err := utils.UploadToCloudinary(avatarFile)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to upload avatar", err.Error())
		}
		req.AvatarUrl = &avatarURL
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidUserData):
			return utils.Error(c, fiber.StatusBadRequest, "No valid fields to update", "")
		case errors.Is(err, apperrors.ErrUserNotFound):
			return utils.Error(c, fiber.StatusNotFound, "User not found", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
		}
	}

	return utils.Success(c, "User profile updated successfully", user)
}

func isImageFile(contentType string) bool {
	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/gif" ||
		contentType == "image/webp"
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch users", err.Error())
	}
	return utils.Success(c, "Users fetched successfully", users)
}
