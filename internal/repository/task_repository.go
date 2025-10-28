package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// Create adds a new task to the database
func (r *TaskRepository) Create(task *models.Task) error {
	return r.DB.Create(task).Error
}

// FindByBoardID retrieves all tasks associated with a specific board
func (r *TaskRepository) FindByBoardID(boardID uuid.UUID) ([]models.Task, error) {

	var tasks []models.Task

	err := r.DB.Where("board_id = ?", boardID).Preload("Assignee", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Labels").Find(&tasks).Error

	return tasks, err
}

// FindByID retrieves a task by its ID
func (r *TaskRepository) FindByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task

	err := r.DB.Preload("Assignee", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Labels").
		First(&task, "id = ?", id).Error
	return &task, err
}

// Update modifies an existing task's details
func (r *TaskRepository) Update(id uuid.UUID, data map[string]interface{}) error {
	return r.DB.Model(&models.Task{}).Where("id = ?", id).Updates(data).Error
}

// SoftDelete marks a task as deleted without removing it from the database
func (r *TaskRepository) SoftDelete(id uuid.UUID) error {
	return r.DB.Delete(&models.Task{}, "id = ?", id).Error
}

// BoardExists checks if a board with the given ID exists
func (r *TaskRepository) BoardExists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Board{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// UserExists checks if a user with the given ID exists
func (r *TaskRepository) UserExists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// ReplaceLabels replaces all labels associated with a task
func (r *TaskRepository) ReplaceLabels(taskID uuid.UUID, labels []models.TaskLabel) error {
	if err := r.DB.Where("task_id = ?", taskID).Delete(&models.TaskLabel{}).Error; err != nil {
		return err
	}
	if len(labels) > 0 {
		if err := r.DB.CreateInBatches(labels, 100).Error; err != nil {
			return err
		}
	}
	return nil
}
