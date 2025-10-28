package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardRepository struct {
	DB *gorm.DB
}

// NewBoardRepository creates a new instance of BoardRepository
func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{DB: db}
}

// Create saves a new board record in the database
func (r *BoardRepository) Create(tx *gorm.DB, board *models.Board) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.Create(board).Error
}

// FindByProjectID retrieves all boards associated with a specific project
func (r *BoardRepository) FindByProjectID(projectID uuid.UUID) ([]models.Board, error) {
	var boards []models.Board
	err := r.DB.Where("project_id = ?", projectID).
		Order("order_index ASC").
		Find(&boards).Error
	return boards, err
}

// FindByID retrieves a board by its ID
func (r *BoardRepository) FindByID(id uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := r.DB.First(&board, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// ShiftBoardOrder adjusts the order_index of boards when a board is moved
func (r *BoardRepository) ShiftBoardOrder(tx *gorm.DB, projectID uuid.UUID, oldIndex, newIndex int) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	if newIndex > oldIndex {
		return db.Model(&models.Board{}).
			Where("project_id = ? AND order_index > ? AND order_index <= ?", projectID, oldIndex, newIndex).
			Update("order_index", gorm.Expr("order_index - 1")).Error
	}
	if newIndex < oldIndex {
		return db.Model(&models.Board{}).
			Where("project_id = ? AND order_index >= ? AND order_index < ?", projectID, newIndex, oldIndex).
			Update("order_index", gorm.Expr("order_index + 1")).Error
	}
	return nil
}

// GetBoardCount returns the number of boards in a project
func (r *BoardRepository) GetBoardCount(projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Board{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// GetMaxOrderIndex retrieves the maximum order_index for boards in a project
func (r *BoardRepository) GetMaxOrderIndex(projectID uuid.UUID) (int, error) {
	var max int
	err := r.DB.Table("boards").
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(order_index), 0)").
		Scan(&max).Error
	return max, err
}

// Update modifies the details of a board
func (r *BoardRepository) Update(tx *gorm.DB, id uuid.UUID, data map[string]interface{}) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.Model(&models.Board{}).Where("id = ?", id).Updates(data).Error
}

// SoftDelete marks a board as deleted without removing it from the database
func (r *BoardRepository) SoftDelete(tx *gorm.DB, id uuid.UUID) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.Board{}, "id = ?", id).Error
}
