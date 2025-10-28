package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// Create adds a new comment to the database
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

// FindByID retrieves a comment by its ID
func (r *CommentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.DB.First(&comment, "id = ?", id).Error
	return &comment, err
}

// GetMainCommentsWithReplies retrieves main comments for a task along with their replies
func (r *CommentRepository) GetMainCommentsWithReplies(taskID uuid.UUID) ([]models.Comment, error) {
	var mainComments []models.Comment
	err := r.DB.Where("task_id = ? AND parent_id IS NULL", taskID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, avatar_url")
		}).
		Order("created_at ASC").
		Find(&mainComments).Error
	if err != nil {
		return nil, err
	}

	// Retrieve all replies for these main comments
	var allReplyIDs []uuid.UUID
	for _, main := range mainComments {
		allReplyIDs = append(allReplyIDs, main.ID)
	}

	if len(allReplyIDs) > 0 {
		var replies []models.Comment
		r.DB.Where("parent_id IN ?", allReplyIDs).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, avatar_url")
			}).
			Order("created_at ASC").
			Find(&replies)

		// Group replies by their parent comment
		replyMap := make(map[uuid.UUID][]models.Comment)
		for _, reply := range replies {
			replyMap[*reply.ParentID] = append(replyMap[*reply.ParentID], reply)
		}

		for i, main := range mainComments {
			if replies, exists := replyMap[main.ID]; exists {
				mainComments[i].Replies = replies
			}
		}
	}

	return mainComments, nil
}

// IsTaskMember checks if a user is a member of the project associated with a task
func (r *CommentRepository) IsTaskMember(taskID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Table("tasks").
		Joins("JOIN boards ON tasks.board_id = boards.id").
		Joins("JOIN projects ON boards.project_id = projects.id").
		Joins("LEFT JOIN project_members ON projects.id = project_members.project_id AND project_members.user_id = ?", userID).
		Where("tasks.id = ? AND (projects.owner_id = ? OR project_members.user_id = ?)", taskID, userID, userID).
		Count(&count).Error
	return count > 0, err
}

// SoftDelete marks a comment as deleted without removing it from the database
func (r *CommentRepository) SoftDelete(id uuid.UUID) error {
	return r.DB.Delete(&models.Comment{}, "id = ?", id).Error
}

// IsOwnerOfTaskProject checks if a user is the owner of the project associated with a task
func (r *CommentRepository) IsOwnerOfTaskProject(taskID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Table("tasks").
		Joins("JOIN boards ON tasks.board_id = boards.id").
		Joins("JOIN projects ON boards.project_id = projects.id").
		Where("tasks.id = ? AND projects.owner_id = ?", taskID, userID).
		Count(&count).Error
	return count > 0, err
}
