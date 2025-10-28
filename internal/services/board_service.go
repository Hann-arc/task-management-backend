package services

import (
	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardService struct {
	BoardRepo          *repository.BoardRepository
	ProjectRepo        *repository.ProjectRepository
	DB                 *gorm.DB
	ActivityLogService *ActivityLogService
}

// NewBoardService creates a new instance of BoardService
func NewBoardService(db *gorm.DB, boardRepo *repository.BoardRepository, projectRepo *repository.ProjectRepository, activityLogService *ActivityLogService) *BoardService {
	return &BoardService{DB: db, BoardRepo: boardRepo, ProjectRepo: projectRepo, ActivityLogService: activityLogService}
}

// CreateBoard creates a new board within a project
func (s *BoardService) CreateBoard(projectID, userID uuid.UUID, name string) (*dto.BoardResponse, error) {
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrUnauthorizedBoardAction
	}

	count, err := s.BoardRepo.GetBoardCount(projectID)
	if err != nil {
		return nil, err
	}

	board := &models.Board{
		ID:         uuid.New(),
		Name:       name,
		ProjectID:  projectID,
		OrderIndex: int(count) + 1,
	}

	if err := s.BoardRepo.Create(nil, board); err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, userID, "board.created", map[string]interface{}{
			"board_id":   board.ID.String(),
			"name":       name,
			"project_id": projectID.String(),
		})
	}

	return &dto.BoardResponse{
		ID:         board.ID.String(),
		Name:       board.Name,
		OrderIndex: board.OrderIndex,
		ProjectID:  board.ProjectID.String(),
		CreatedAt:  board.CreatedAt,
		UpdatedAt:  board.UpdatedAt,
	}, nil
}

// GetBoards retrieves all boards for a given project
func (s *BoardService) GetBoards(projectID uuid.UUID) ([]dto.BoardResponse, error) {
	boards, err := s.BoardRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	var result []dto.BoardResponse
	for _, b := range boards {
		result = append(result, dto.BoardResponse{
			ID:         b.ID.String(),
			Name:       b.Name,
			OrderIndex: b.OrderIndex,
			ProjectID:  b.ProjectID.String(),
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,
		})
	}
	return result, nil
}

// UpdateBoard updates the details of a board, including its name and order index
func (s *BoardService) UpdateBoard(boardID, userID uuid.UUID, name *string, orderIndex *int) (*dto.BoardResponse, error) {
	board, err := s.BoardRepo.FindByID(boardID)
	if err != nil {
		return nil, err
	}

	isOwner, err := s.ProjectRepo.IsOwner(board.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrUnauthorizedBoardAction
	}

	if name == nil && orderIndex == nil {
		return nil, apperrors.ErrNoFieldsToUpdate
	}

	var newOrderIndex *int
	if orderIndex != nil {
		maxOrder, err := s.BoardRepo.GetMaxOrderIndex(board.ProjectID)
		if err != nil {
			return nil, err
		}
		if *orderIndex < 1 || *orderIndex > maxOrder+1 {
			return nil, apperrors.ErrInvalidOrderIndex
		}
		if *orderIndex != board.OrderIndex {
			newOrderIndex = orderIndex
		}
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	data := map[string]interface{}{}
	if name != nil {
		data["name"] = *name
	}

	if newOrderIndex != nil {
		if err := s.BoardRepo.ShiftBoardOrder(tx, board.ProjectID, board.OrderIndex, *newOrderIndex); err != nil {
			tx.Rollback()
			return nil, err
		}
		data["order_index"] = *newOrderIndex
	}

	if len(data) > 0 {
		if err := s.BoardRepo.Update(tx, boardID, data); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	updatedBoard, err := s.BoardRepo.FindByID(boardID)
	if err != nil {
		return nil, err
	}

	// Log activity

	if s.ActivityLogService != nil {
		details := map[string]interface{}{
			"board_id": boardID.String(),
		}
		if name != nil {
			details["name"] = *name
		}
		if orderIndex != nil {
			details["order_index"] = *orderIndex
		}
		s.ActivityLogService.LogActivity(board.ProjectID, userID, "board.updated", details)
	}

	return &dto.BoardResponse{
		ID:         updatedBoard.ID.String(),
		Name:       updatedBoard.Name,
		OrderIndex: updatedBoard.OrderIndex,
		ProjectID:  updatedBoard.ProjectID.String(),
		CreatedAt:  updatedBoard.CreatedAt,
		UpdatedAt:  updatedBoard.UpdatedAt,
	}, nil
}

// DeleteBoard handles the deletion of a board
func (s *BoardService) DeleteBoard(boardID, userID uuid.UUID) error {
	board, err := s.BoardRepo.FindByID(boardID)
	if err != nil {
		return err
	}

	isOwner, err := s.ProjectRepo.IsOwner(board.ProjectID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrUnauthorizedBoardAction
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := s.BoardRepo.SoftDelete(tx, boardID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.BoardRepo.ShiftBoardOrder(tx, board.ProjectID, board.OrderIndex, 0); err != nil {
		tx.Rollback()
		return err
	}

	// Log activity

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(board.ProjectID, userID, "board.deleted", map[string]interface{}{
			"board_id": boardID.String(),
		})
	}

	return tx.Commit().Error
}
