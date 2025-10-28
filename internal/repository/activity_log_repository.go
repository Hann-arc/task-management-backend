package repository

import (
	"encoding/json"
	"time"

	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	DB *gorm.DB
}

// NewActivityLogRepository creates a new instance of ActivityLogRepository
func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

// Create adds a new activity log entry to the database
func (r *ActivityLogRepository) Create(projectID, userID uuid.UUID, action string, details map[string]interface{}) error {
	var detailsBytes []byte
	if details != nil {
		var err error
		detailsBytes, err = json.Marshal(details)
		if err != nil {
			return err
		}
	}

	log := &models.ActivityLog{
		ProjectID: projectID,
		UserID:    userID,
		Action:    action,
		Details:   detailsBytes,
		CreatedAt: time.Now(),
	}

	return r.DB.Create(log).Error
}

// FindByProjectID retrieves activity logs for a specific project with pagination
func (r *ActivityLogRepository) FindByProjectID(projectID uuid.UUID, limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.DB.Where("project_id = ?", projectID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// IsProjectMember checks if a user is a member of a specific project
func (r *ActivityLogRepository) IsProjectMember(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Table("projects").
		Where("id = ? AND owner_id = ?", projectID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	err = r.DB.Table("project_members").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return err == nil && count > 0, nil
}
