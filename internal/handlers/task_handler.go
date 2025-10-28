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

type TaskHandler struct {
	service *services.TaskService
}

// NewTaskHandler creates a new instance of TaskHandler
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask handles the creation of a new task within a board
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	boardID, err := uuid.Parse(c.Params("boardId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid board ID", "")
	}

	var req dto.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	if req.Title == "" {
		return utils.Error(c, fiber.StatusBadRequest, "Title is required", "")
	}

	if req.Priority == "" {
		return utils.Error(c, fiber.StatusBadRequest, "Priority is required", "")
	}

	// Validate AssigneeID format if provided
	if req.AssigneeID != nil {
		if _, err := uuid.Parse(*req.AssigneeID); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "Invalid assignee ID format", "")
		}
	}

	task, err := h.service.CreateTask(&req, boardID, userID)
	if err != nil {

		switch {

		case errors.Is(err, apperrors.ErrUnauthorizedTask):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		case errors.Is(err, apperrors.ErrBoardNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Board not found", "")
		case errors.Is(err, apperrors.ErrAssigneeNotFound):
			return utils.Error(c, fiber.StatusBadRequest, "Assignee not found", "")
		case errors.Is(err, apperrors.ErrInvalidTaskData):
			return utils.Error(c, fiber.StatusBadRequest, "Invalid task data", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to create task", err.Error())
		}
	}

	return utils.Created(c, "Task created successfully", task)
}

// GetTasksByBoard retrieves all tasks for a given board
func (h *TaskHandler) GetTasksByBoard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	boardID, err := uuid.Parse(c.Params("boardId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid board ID", "")
	}

	tasks, err := h.service.GetTasksByBoard(boardID, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedTask):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch tasks", err.Error())
		}
	}

	return utils.Success(c, "Tasks fetched successfully", tasks)
}

// UpdateTask handles updating a task's details
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task ID", "")
	}

	var req dto.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	task, err := h.service.UpdateTask(taskID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTaskNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Task not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedTask):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		case errors.Is(err, apperrors.ErrAssigneeNotFound):
			return utils.Error(c, fiber.StatusBadRequest, "Assignee not found", "")
		case errors.Is(err, apperrors.ErrInvalidTaskData):
			return utils.Error(c, fiber.StatusBadRequest, "No valid fields to update", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to update task", err.Error())
		}
	}

	return utils.Success(c, "Task updated successfully", task)
}

// DeleteTask handles the deletion of a task
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	taskID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid task ID", "")
	}

	if err := h.service.DeleteTask(taskID, userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTaskNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Task not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedTask):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to delete task", err.Error())
		}
	}

	return utils.Success(c, "Task deleted successfully", nil)
}
