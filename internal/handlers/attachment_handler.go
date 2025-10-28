package handlers

import (
	"errors"

	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentHandler struct {
	service *services.AttachmentService
}

// NewAttachmentHandler creates a new instance of AttachmentHandler
func NewAttachmentHandler(service *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

// UploadAttachment handles the uploading of an attachment to a task
func (h *AttachmentHandler) UploadAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task ID", "")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "File is required", "")
	}

	if file.Size > 10*1024*1024 {
		return utils.Error(c, fiber.StatusBadRequest, "File too large (max 10MB)", "")
	}

	attachment, err := h.service.UploadAttachment(taskID, userID, file)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to upload attachment", err.Error())
		}
	}

	return utils.Created(c, "Attachment uploaded successfully", attachment)
}

// GetAttachments retrieves all attachments for a specific task
func (h *AttachmentHandler) GetAttachments(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task ID", "")
	}

	attachments, err := h.service.GetAttachments(taskID, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch attachments", err.Error())
		}
	}

	return utils.Success(c, "Attachments fetched successfully", attachments)
}

// DeleteAttachment handles the deletion of an attachment
func (h *AttachmentHandler) DeleteAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	attachmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid attachment ID", "")
	}

	if err := h.service.DeleteAttachment(attachmentID, userID); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Attachment not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "You can only delete your own attachments or you are not project owner", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to delete attachment", err.Error())
		}
	}

	return utils.Success(c, "Attachment deleted successfully", nil)
}
