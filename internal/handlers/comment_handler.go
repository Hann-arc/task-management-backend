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

type CommentHandler struct {
	service *services.CommentService
}

// NewCommentHandler creates a new instance of CommentHandler
func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// CreateMainComment handles the creation of a main comment on a task
func (h *CommentHandler) CreateMainComment(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task id", "")
	}

	var req dto.CreateCommentRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid body request", "")
	}

	comment, err := h.service.CreateMainComment(taskID, userID, req.Content)

	if err != nil {
		switch {

		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to create comment", err.Error())
		}

	}

	return utils.Created(c, "Comment created successfully", comment)
}

// CreateReply handles the creation of a reply to an existing comment
func (h *CommentHandler) CreateReply(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	commentID, err := uuid.Parse(c.Params("commentId"))

	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid comment id", "")
	}

	var req dto.CreateCommentRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	comment, err := h.service.CreateReply(commentID, userID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrCommentNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Comment not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to create reply", err.Error())
		}
	}

	return utils.Created(c, "Reply created successfully", comment)
}

// GetComments retrieves all comments for a given task
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task ID", "")
	}

	comments, err := h.service.GetCommentsByTask(taskID, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch comments", err.Error())
		}
	}

	return utils.Success(c, "Comments fetched successfully", comments)
}

// DeleteComment handles the deletion of a comment
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid comment ID", "")
	}

	if err := h.service.DeleteComment(commentID, userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrCommentNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Comment not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "You can only delete your own comments or you are not project owner", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to delete comment", err.Error())
		}
	}

	return utils.Success(c, "Comment deleted successfully", nil)
}
