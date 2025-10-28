package services

import (
	"errors"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService struct {
	ProjectRepo        *repository.ProjectRepository
	UserRepo           *repository.UserRepository
	ActivityLogService *ActivityLogService
}

// NewProjectService creates a new instance of ProjectService
func NewProjectService(projectRepo *repository.ProjectRepository, userRepo *repository.UserRepository, activityLogService *ActivityLogService) *ProjectService {
	return &ProjectService{ProjectRepo: projectRepo, UserRepo: userRepo, ActivityLogService: activityLogService}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req *dto.CreateProjectRequest, ownerID uuid.UUID) (*dto.ProjectResponse, error) {
	project := &models.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
	}

	if err := s.ProjectRepo.Create(project); err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(project.ID, ownerID, "project.created", map[string]interface{}{
			"project_id": project.ID.String(),
			"name":       req.Name,
		})
	}

	return &dto.ProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID.String(),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		DeletedAt:   utils.ToTimePtr(project.DeletedAt),
	}, nil
}

// GetProjects retrieves all projects for a specific user
func (s *ProjectService) GetProjects(userID uuid.UUID) ([]dto.ProjectResponse, error) {
	projects, err := s.ProjectRepo.FindAllByUser(userID)
	if err != nil {
		return nil, err
	}

	var result []dto.ProjectResponse
	for _, p := range projects {
		result = append(result, dto.ProjectResponse{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			OwnerID:     p.OwnerID.String(),
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			DeletedAt:   utils.ToTimePtr(p.DeletedAt),
		})
	}

	return result, nil
}

// GetProjectDetail retrieves detailed information about a specific project
func (s *ProjectService) GetProjectDetail(projectID, userID uuid.UUID) (*dto.ProjectDetailResponse, error) {

	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	isMember := false
	if !isOwner {
		isMember, err = s.ProjectRepo.IsMember(projectID, userID)
		if err != nil {
			return nil, err
		}
	}
	if !isOwner && !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	project, err := s.ProjectRepo.FindByIDWithDetails(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrProjectNotFound
		}
		return nil, err
	}

	ownerResp := dto.UserResponse{
		ID:        project.Owner.ID.String(),
		Name:      project.Owner.Name,
		Email:     project.Owner.Email,
		AvatarUrl: project.Owner.AvatarUrl,
	}

	var boards []dto.BoardSummary
	for _, b := range project.Boards {
		boards = append(boards, dto.BoardSummary{
			ID:         b.ID.String(),
			Name:       b.Name,
			OrderIndex: b.OrderIndex,
		})
	}

	return &dto.ProjectDetailResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		Owner:       ownerResp,
		Boards:      boards,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		DeletedAt:   utils.ToTimePtr(project.DeletedAt),
	}, nil
}

// UpdateProject updates the details of a specific project
func (s *ProjectService) UpdateProject(projectID, userID uuid.UUID, req *dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrUnauthorizedOwnerOnly
	}

	if req.Name == nil && req.Description == nil {
		return nil, apperrors.ErrNoFieldsToUpdate
	}

	data := map[string]interface{}{}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}

	if err := s.ProjectRepo.Update(projectID, data); err != nil {
		return nil, err
	}

	updatedProject, err := s.ProjectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		details := map[string]interface{}{
			"project_id": projectID.String(),
		}
		if req.Name != nil {
			details["name"] = *req.Name
		}
		if req.Description != nil {
			details["description"] = *req.Description
		}
		s.ActivityLogService.LogActivity(projectID, userID, "project.updated", details)
	}

	return &dto.ProjectResponse{
		ID:          updatedProject.ID.String(),
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		OwnerID:     updatedProject.OwnerID.String(),
		CreatedAt:   updatedProject.CreatedAt,
		UpdatedAt:   updatedProject.UpdatedAt,
		DeletedAt:   utils.ToTimePtr(updatedProject.DeletedAt),
	}, nil
}

// DeleteProject deletes a specific project
func (s *ProjectService) DeleteProject(projectID, userID uuid.UUID) error {
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrUnauthorizedOwnerOnly
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, userID, "project.deleted", map[string]interface{}{
			"project_id": projectID.String(),
		})
	}

	return s.ProjectRepo.SoftDelete(projectID)
}
