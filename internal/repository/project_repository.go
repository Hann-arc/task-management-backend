package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

// NewProjectRepository creates a new instance of ProjectRepository
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

// Create adds a new project to the database
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.DB.Create(project).Error
}

// FindByID retrieves a project by its ID
func (r *ProjectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.DB.Where("id = ?", id).First(&project).Error
	return &project, err
}

// FindByIDWithDetails retrieves a project by its ID along with its owner and boards
func (r *ProjectRepository) FindByIDWithDetails(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.DB.Preload("Owner").
		Preload("Boards", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, order_index, project_id")
		}).
		Where("id = ?", id).First(&project).Error
	return &project, err
}

// FindAllByUser retrieves all projects where the user is either the owner or a member
func (r *ProjectRepository) FindAllByUser(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	subQuery := r.DB.Model(&models.ProjectMember{}).Select("project_id").Where("user_id = ?", userID)

	err := r.DB.Where("owner_id = ? OR id IN (?)", userID, subQuery).
		Find(&projects).Error
	return projects, err
}

// Update modifies an existing project's details
func (r *ProjectRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.DB.Model(&models.Project{}).Where("id = ?", id).Updates(data).Error
}

// SoftDelete marks a project as deleted without removing it from the database
func (r *ProjectRepository) SoftDelete(id uuid.UUID) error {
	return r.DB.Delete(&models.Project{}, "id = ?", id).Error
}

// IsOwner checks if a user is the owner of a specific project
func (r *ProjectRepository) IsOwner(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Project{}).
		Where("id = ? AND owner_id = ?", projectID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsMember checks if a user is a member of a specific project
func (r *ProjectRepository) IsMember(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return err == nil && count > 0, nil
}
