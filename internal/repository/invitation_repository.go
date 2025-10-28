package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvitationRepository struct {
	DB *gorm.DB
}

// NewInvitationRepository creates a new instance of InvitationRepository
func NewInvitationRepository(db *gorm.DB) *InvitationRepository {
	return &InvitationRepository{DB: db}
}

// Create adds a new invitation to the database
func (r *InvitationRepository) Create(invitation *models.Invitation) error {
	return r.DB.Create(invitation).Error
}

// FindByToken retrieves an invitation by its token
func (r *InvitationRepository) FindByToken(token string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.DB.Where("token = ?", token).
		Preload("Project").
		First(&invitation).Error
	return &invitation, err
}

// FindPendingByProjectAndEmail retrieves a pending invitation by project ID and email
func (r *InvitationRepository) FindPendingByProjectAndEmail(projectID uuid.UUID, email string) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.DB.Where("project_id = ? AND email = ? AND status = ?", projectID, email, "pending").
		First(&invitation).Error
	return &invitation, err
}

// UpdateStatus updates the status of an invitation
func (r *InvitationRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.DB.Model(&models.Invitation{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// IsMember checks if a user is already a member of a specific project
func (r *InvitationRepository) IsMember(projectID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}
