package handlers

import (
	"errors"

	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BoardHandler struct {
	service *services.BoardService
}

// NewBoardHandler creates a new instance of BoardHandler
func NewBoardHandler(service *services.BoardService) *BoardHandler {
	return &BoardHandler{service: service}
}

// CreateBoard handles the creation of a new board within a project
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return utils.Error(c, fiber.StatusBadRequest, "Board name is required", "")
	}

	board, err := h.service.CreateBoard(projectID, userID, req.Name)
	if err != nil {
		if errors.Is(err, apperrors.ErrUnauthorizedBoardAction) {
			return utils.Error(c, fiber.StatusForbidden, "You are not authorized to create boards in this project", "")
		}
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to create board", err.Error())
	}

	return utils.Created(c, "Board created successfully", board)
}

// GetBoards retrieves all boards for a given project
func (h *BoardHandler) GetBoards(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	boards, err := h.service.GetBoards(projectID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch boards", err.Error())
	}

	return utils.Success(c, "Boards fetched successfully", boards)
}

// UpdateBoard handles updating a board's details such as name and order index
func (h *BoardHandler) UpdateBoard(c *fiber.Ctx) error {
	boardID, err := uuid.Parse(c.Params("boardId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid board ID", "")
	}

	userID := c.Locals("user_id").(uuid.UUID)

	var payload struct {
		Name       *string `json:"name"`
		OrderIndex *int    `json:"order_index"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid payload", "")
	}

	board, err := h.service.UpdateBoard(boardID, userID, payload.Name, payload.OrderIndex)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedBoardAction):
			return utils.Error(c, fiber.StatusForbidden, "You are not authorized to update this board", "")
		case errors.Is(err, apperrors.ErrInvalidOrderIndex), errors.Is(err, apperrors.ErrNoFieldsToUpdate):
			return utils.Error(c, fiber.StatusBadRequest, "Invalid request", err.Error())
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to update board", err.Error())
		}
	}

	return utils.Success(c, "Board updated successfully", board)
}

// DeleteBoard handles the deletion of a board
func (h *BoardHandler) DeleteBoard(c *fiber.Ctx) error {
	boardID, err := uuid.Parse(c.Params("boardId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid board ID", "")
	}

	userID := c.Locals("user_id").(uuid.UUID)

	if err := h.service.DeleteBoard(boardID, userID); err != nil {
		if errors.Is(err, apperrors.ErrUnauthorizedBoardAction) {
			return utils.Error(c, fiber.StatusForbidden, "You are not authorized to delete this board", "")
		}
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to delete board", err.Error())
	}

	return utils.Success(c, "Board deleted successfully", nil)
}
