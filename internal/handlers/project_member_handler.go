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

type ProjectMemberHandler struct {
	service *services.ProjectMemberService
}

// NewProjectMemberHandler creates a new instance of ProjectMemberHandler
func NewProjectMemberHandler(service *services.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{service: service}
}

// AddMember adds a new member to the project
func (h *ProjectMemberHandler) AddMember(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	var req dto.AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	member, err := h.service.AddMember(projectID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "Only project owner can add members", "")
		case errors.Is(err, apperrors.ErrUserNotFound):
			return utils.Error(c, fiber.StatusBadRequest, "User with this email not found", "")
		case errors.Is(err, apperrors.ErrAlreadyMember):
			return utils.Error(c, fiber.StatusBadRequest, "User is already a member of this project", "")
		case errors.Is(err, apperrors.ErrInviteeIsOwner):
			return utils.Error(c, fiber.StatusBadRequest, "Cannot invite project owner", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to add member", err.Error())
		}
	}

	return utils.Created(c, "Member added successfully", member)
}

// GetMembers retrieves all members of the project
func (h *ProjectMemberHandler) GetMembers(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	members, err := h.service.GetMembers(projectID, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch members", err.Error())
		}
	}

	return utils.Success(c, "Members fetched successfully", members)
}

// RemoveMember removes a member from the project
func (h *ProjectMemberHandler) RemoveMember(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	targetUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid user ID", "")
	}

	if err := h.service.RemoveMember(projectID, userID, targetUserID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "Only project owner can remove members", "")
		case errors.Is(err, apperrors.ErrCannotRemoveSelf):
			return utils.Error(c, fiber.StatusBadRequest, "You cannot remove yourself", "")
		case errors.Is(err, apperrors.ErrProjectMemberNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Member not found", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to remove member", err.Error())
		}
	}

	return utils.Success(c, "Member removed successfully", nil)
}
