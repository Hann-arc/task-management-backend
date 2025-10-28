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

type ProjectHandler struct {
	service *services.ProjectService
}

// NewProjectHandler creates a new instance of ProjectHandler
func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// CreateProject handles the creation of a new project
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	project, err := h.service.CreateProject(&req, userID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to create project", err.Error())
	}

	return utils.Created(c, "Project created successfully", project)
}

// ListProjects retrieves all projects for the authenticated user
func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	projects, err := h.service.GetProjects(userID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch projects", err.Error())
	}

	return utils.Success(c, "Projects fetched successfully", projects)
}

// GetProject retrieves details of a specific project
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	userID := c.Locals("user_id").(uuid.UUID)

	project, err := h.service.GetProjectDetail(projectID, userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrProjectNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Project not found", "")
		case errors.Is(err, apperrors.ErrUnauthorizedProject):
			return utils.Error(c, fiber.StatusForbidden, "You are not a member of this project", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to fetch project", err.Error())
		}
	}

	return utils.Success(c, "Project fetched successfully", project)
}

// UpdateProject updates the details of a specific project
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", "")
	}

	project, err := h.service.UpdateProject(projectID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "Only project owner can update this project", "")
		case errors.Is(err, apperrors.ErrNoFieldsToUpdate):
			return utils.Error(c, fiber.StatusBadRequest, "No fields to update", "")
		case errors.Is(err, apperrors.ErrProjectNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Project not found", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to update project", err.Error())
		}
	}

	return utils.Success(c, "Project updated successfully", project)
}

// DeleteProject deletes a specific project
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid project ID", "")
	}

	userID := c.Locals("user_id").(uuid.UUID)

	if err := h.service.DeleteProject(projectID, userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUnauthorizedOwnerOnly):
			return utils.Error(c, fiber.StatusForbidden, "Only project owner can delete this project", "")
		case errors.Is(err, apperrors.ErrProjectNotFound):
			return utils.Error(c, fiber.StatusNotFound, "Project not found", "")
		default:
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to delete project", err.Error())
		}
	}

	return utils.Success(c, "Project deleted successfully", nil)
}
