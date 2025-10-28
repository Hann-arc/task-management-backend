package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMemberRepository struct {
	DB *gorm.DB
}

// NewProjectMemberRepository creates a new instance of ProjectMemberRepository
func NewProjectMemberRepository(db *gorm.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{DB: db}
}

// Create adds a new project member to the database
func (r *ProjectMemberRepository) Create(member *models.ProjectMember) error {
	return r.DB.Create(member).Error
}

// FindByProjectID retrieves all members of a specific project
func (r *ProjectMemberRepository) FindByProjectID(projectID uuid.UUID) ([]models.ProjectMember, error) {
	var members []models.ProjectMember
	err := r.DB.Where("project_id = ?", projectID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, email")
		}).
		Preload("Role").
		Find(&members).Error
	return members, err
}

// FindByProjectAndUser retrieves a project member by project ID and user ID
func (r *ProjectMemberRepository) FindByProjectAndUser(projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	var member models.ProjectMember
	err := r.DB.Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&member).Error
	return &member, err
}

// Delete removes a project member from the database
func (r *ProjectMemberRepository) Delete(projectID, userID uuid.UUID) error {
	return r.DB.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

// UserExistsByEmail checks if a user exists by their email
func (r *ProjectMemberRepository) UserExistsByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// IsOwner checks if a user is the owner of a specific project
func (r *ProjectMemberRepository) IsOwner(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Project{}).
		Where("id = ? AND owner_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}
