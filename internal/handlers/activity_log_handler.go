package handlers

import (
	"errors"
	"strconv"

	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ActivityLogHandler struct {
	service *services.ActivityLogService
}

// NewActivityLogHandler creates a new instance of ActivityLogHandler
func NewActivityLogHandler(service *services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

// GetActivityLogs retrieves activity logs for a specific project
func (h *ActivityLogHandler) GetActivityLogs(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project id", "")
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	logs, err := h.service.GetActivityLogs(projectID, userID, limit, offset)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to featch activity log", "")
		}
	}

	return utils.Success(c, "Activoty log fetched successfully", logs)
}
