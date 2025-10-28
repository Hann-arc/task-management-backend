package services

import (
	"encoding/json"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/google/uuid"
)

type ActivityLogService struct {
	Repo *repository.ActivityLogRepository
}

// NewActivityLogService creates a new instance of ActivityLogService
func NewActivityLogService(repo *repository.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{Repo: repo}
}

// LogActivity records a new activity in the activity log
func (s *ActivityLogService) LogActivity(projectID, userID uuid.UUID, action string, details map[string]interface{}) error {
	return s.Repo.Create(projectID, userID, action, details)
}

// GetActivityLogs retrieves activity logs for a specific project
func (s *ActivityLogService) GetActivityLogs(projectID, userID uuid.UUID, limit, offset int) ([]dto.ActivityLogResponse, error) {
	isMemeber, err := s.Repo.IsProjectMember(projectID, userID)

	if err != nil {
		return nil, err
	}

	if !isMemeber {
		return nil, apperrors.ErrUnauthorizedProject
	}

	logs, err := s.Repo.FindByProjectID(projectID, limit, offset)

	if err != nil {
		return nil, err
	}

	var result []dto.ActivityLogResponse

	for _, log := range logs {
		details := make(map[string]interface{})
		if len(log.Details) > 0 {
			_ = json.Unmarshal(log.Details, &details)
		}

		result = append(result, dto.ActivityLogResponse{
			ID:        log.ID.String(),
			ProjectID: log.ProjectID.String(),
			UserID:    log.UserID.String(),
			User: dto.ActivityLogUser{
				ID:   log.User.ID.String(),
				Name: log.User.Name,
			},
			Action:    log.Action,
			Details:   details,
			CreatedAt: log.CreatedAt,
		})

	}

	return result, nil
}
