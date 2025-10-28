package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	DB *gorm.DB
}

// NewAttachmentRepository creates a new instance of AttachmentRepository
func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{DB: db}
}

// Create saves a new attachment record in the database
func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.DB.Create(attachment).Error
}

// FindByTaskID retrieves all attachments associated with a specific task
func (r *AttachmentRepository) FindByTaskID(taskID uuid.UUID) ([]models.Attachment, error) {
	var attachments []models.Attachment

	err := r.DB.Where("task_id = ?", taskID).Preload("Uploader", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Find(&attachments).Error

	return attachments, err
}

// IsTaskMember checks if a user is a member of the project that the task belongs to
func (r *AttachmentRepository) IsTaskMember(taskID, userID uuid.UUID) (bool, error) {
	var count int64

	err := r.DB.Table("tasks").Joins("JOIN boards ON tasks.board_id = boards.id").
		Joins("JOIN projects ON boards.project_id = projects.id").
		Joins("LEFT JOIN project_members ON projects.id = project_members.project_id AND project_members.user_id = ?", userID).
		Where("tasks.id = ? AND (projects.owner_id = ? OR project_members.user_id = ?)", taskID, userID, userID).
		Count(&count).Error

	return count > 0, err
}

// Delete removes an attachment record from the database
func (r *AttachmentRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Attachment{}, "id = ?", id).Error
}

// FindByID retrieves an attachment by its ID
func (r *AttachmentRepository) FindByID(id uuid.UUID) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.DB.First(&attachment, "id = ?", id).Error
	return &attachment, err
}

// IsOwnerOfAttachmentProject checks if the user is the owner of the project associated with the attachment
func (r *AttachmentRepository) IsOwnerOfAttachmentProject(attachmentID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Table("attachments").
		Joins("JOIN tasks ON attachments.task_id = tasks.id").
		Joins("JOIN boards ON tasks.board_id = boards.id").
		Joins("JOIN projects ON boards.project_id = projects.id").
		Where("attachments.id = ? AND projects.owner_id = ?", attachmentID, userID).
		Count(&count).Error
	return count > 0, err
}
