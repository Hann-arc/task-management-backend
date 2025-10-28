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

type InvitationHandler struct {
	service *services.InvitationService
}

// NewInvitationHandler creates a new instance of InvitationHandler
func NewInvitationHandler(service *services.InvitationService) *InvitationHandler {
	return &InvitationHandler{service: service}
}

// CreateInvitation creates a new invitation for a user to join a project
func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	var req dto.CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	invitation, err := h.service.CreateInvitation(projectID, userID, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "Only project owner can send invitations", "")
		case errors.Is(err, apperrors.ErrCannotInviteSelf):
			return utils.Error(c, fiber.StatusBadRequest, "Cannot invite yourself", "")
		case errors.Is(err, apperrors.ErrAlreadyMember):
			return utils.Error(c, fiber.StatusBadRequest, "User is already a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to send invitation", err.Error())
		}
	}

	return utils.Created(c, "Invitation sent successfully", invitation)
}

// AcceptInvitation allows a user to accept an invitation to join a project
func (h *InvitationHandler) AcceptInvitation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.AcceptInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	if err := h.service.AcceptInvitation(req.Token, userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvitationNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Invitation not found or email mismatch", "")
		case errors.Is(err, apperrors.ErrInvitationUsed):
			return utils.Error(c, fiber.StatusBadRequest, "Invitation already used", "")
		case errors.Is(err, apperrors.ErrAlreadyMember):
			return utils.Error(c, fiber.StatusBadRequest, "You are already a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to accept invitation", err.Error())
		}
	}

	return utils.Success(c, "Invitation accepted successfully", nil)
}
