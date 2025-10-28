package handlers

import (
	"strconv"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service *services.NotificationService
}

// NewNotificationHandler creates a new instance of NotificationHandler
func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetNotifications returns a list of notifications for the logged-in user
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	notifications, err := h.service.GetNotifications(userID, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch notifications", err.Error())
	}

	return utils.Success(c, "Notifications fetched successfully", notifications)
}

// MarkAsRead marks a specific notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	if err := h.service.MarkAsRead(req.IDs, userID); err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to mark notifications as read", err.Error())
	}

	return utils.Success(c, "Notifications marked as read", nil)
}

// MarkAllAsRead marks all notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	if err := h.service.MarkAllAsRead(userID); err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to mark all notifications as read", err.Error())
	}

	return utils.Success(c, "All notifications marked as read", nil)
}
