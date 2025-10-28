package services

import (
	"errors"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	CommentRepo         *repository.CommentRepository
	TaskRepo            *repository.TaskRepository
	ActivityLogService  *ActivityLogService
	NotificationService *NotificationService
}

// NewCommentService creates a new instance of CommentService
func NewCommentService(commentRepo *repository.CommentRepository, taskRepo *repository.TaskRepository, activityLogService *ActivityLogService, notificationService *NotificationService) *CommentService {
	return &CommentService{CommentRepo: commentRepo, TaskRepo: taskRepo, ActivityLogService: activityLogService, NotificationService: notificationService}
}

// CreateMainComment handles the creation of a main comment on a task
func (s *CommentService) CreateMainComment(taskId, userID uuid.UUID, content string) (*dto.CommentResponse, error) {
	isMember, err := s.CommentRepo.IsTaskMember(taskId, userID)

	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	comment := &models.Comment{
		TaskID:  taskId,
		UserID:  userID,
		Content: content,
	}

	if err := s.CommentRepo.Create(comment); err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		var projectID uuid.UUID
		s.TaskRepo.DB.Model(&models.Task{}).
			Select("projects.id").
			Joins("JOIN boards ON tasks.board_id = boards.id").
			Joins("JOIN projects ON boards.project_id = projects.id").
			Where("tasks.id = ?", taskId).
			Scan(&projectID)

		s.ActivityLogService.LogActivity(projectID, userID, "comment.added", map[string]interface{}{
			"comment_id": comment.ID.String(),
			"task_id":    taskId.String(),
			"content":    content,
		})
	}

	// Send notifications
	if s.NotificationService != nil {
		var task models.Task
		s.TaskRepo.DB.Preload("Assignee").Preload("Creator").First(&task, taskId)

		// Notification to assignee
		if task.AssigneeID != nil && *task.AssigneeID != userID {
			go s.NotificationService.CreateNotification(
				*task.AssigneeID,
				userID,
				"comment.added",
				"task",
				taskId,
				"You have a new comment on your task",
			)
		}

		// Notification to task creator
		if task.CreatedBy != userID {
			go s.NotificationService.CreateNotification(
				task.CreatedBy,
				userID,
				"comment.added",
				"task",
				taskId,
				"New comment on your task",
			)
		}
	}

	return s.buildCommentResponse(comment), nil
}

// CreateReply handles the creation of a reply to an existing comment
func (s *CommentService) CreateReply(targetCommentID, userID uuid.UUID, content string) (*dto.CommentResponse, error) {
	targetComment, err := s.CommentRepo.FindByID(targetCommentID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCommentNotFound
		}
		return nil, err
	}

	isMember, err := s.CommentRepo.IsTaskMember(targetComment.TaskID, userID)

	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	var actualParentID uuid.UUID

	if targetComment.ParentID == nil {
		actualParentID = targetComment.ID
	} else {
		actualParentID = *targetComment.ParentID
	}

	comment := models.Comment{
		TaskID:   targetComment.TaskID,
		UserID:   userID,
		ParentID: &actualParentID,
		Content:  content,
	}

	if err := s.CommentRepo.Create(&comment); err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		var projectID uuid.UUID
		s.TaskRepo.DB.Model(&models.Task{}).
			Select("projects.id").
			Joins("JOIN boards ON tasks.board_id = boards.id").
			Joins("JOIN projects ON boards.project_id = projects.id").
			Where("tasks.id = ?", comment.TaskID).
			Scan(&projectID)

		s.ActivityLogService.LogActivity(projectID, userID, "comment.added", map[string]interface{}{
			"comment_id": comment.ID.String(),
			"task_id":    comment.TaskID.String(),
			"content":    content,
		})
	}

	if s.NotificationService != nil {
		var task models.Task
		s.TaskRepo.DB.Preload("Assignee").Preload("Creator").First(&task, task.ID)

		// Notifikasi ke assignee
		if task.AssigneeID != nil && *task.AssigneeID != userID {
			go s.NotificationService.CreateNotification(
				*task.AssigneeID,
				userID,
				"comment.added",
				"task",
				task.ID,
				"You have a new comment on your task",
			)
		}

		// Notifikasi ke creator task
		if task.CreatedBy != userID {
			go s.NotificationService.CreateNotification(
				task.CreatedBy,
				userID,
				"comment.added",
				"task",
				task.ID,
				"New comment on your task",
			)
		}
	}

	return s.buildCommentResponse(&comment), nil

}

// GetCommentsByTask retrieves all comments and their replies for a given task
func (s *CommentService) GetCommentsByTask(taskID, userID uuid.UUID) ([]dto.CommentResponse, error) {
	isMember, err := s.CommentRepo.IsTaskMember(taskID, userID)

	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	mainComments, err := s.CommentRepo.GetMainCommentsWithReplies(taskID)

	if err != nil {
		return nil, err
	}

	var result []dto.CommentResponse
	for _, main := range mainComments {
		result = append(result, *s.buildCommentResponse(&main))

		for _, reply := range main.Replies {
			result = append(result, *s.buildCommentResponse(&reply))
		}
	}

	return result, nil
}

// DeleteComment handles the deletion of a comment
func (s *CommentService) DeleteComment(commentID, userID uuid.UUID) error {
	comment, err := s.CommentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	isOwner, err := s.CommentRepo.IsOwnerOfTaskProject(comment.TaskID, userID)
	if err != nil {
		return err
	}
	if isOwner || comment.UserID == userID {

		// Log activity
		if s.ActivityLogService != nil {
			var projectID uuid.UUID
			s.CommentRepo.DB.Model(&models.Comment{}).
				Select("projects.id").
				Joins("JOIN tasks ON comments.task_id = tasks.id").
				Joins("JOIN boards ON tasks.board_id = boards.id").
				Joins("JOIN projects ON boards.project_id = projects.id").
				Where("comments.id = ?", commentID).
				Scan(&projectID)

			s.ActivityLogService.LogActivity(projectID, userID, "comment.deleted", map[string]interface{}{
				"comment_id": commentID.String(),
			})
		}

		return s.CommentRepo.SoftDelete(commentID)
	}

	return apperrors.ErrUnauthorizedOwnerOnly
}

// buildCommentResponse transforms a Comment model into a CommentResponse DTO
func (s *CommentService) buildCommentResponse(comment *models.Comment) *dto.CommentResponse {
	res := &dto.CommentResponse{
		ID:      comment.ID.String(),
		TaskID:  comment.TaskID.String(),
		UserID:  comment.UserID.String(),
		Content: comment.Content,
		User: dto.CommentUser{
			ID:        comment.User.ID.String(),
			Name:      comment.User.Name,
			AvatarUrl: comment.User.AvatarUrl,
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	if comment.ParentID != nil {
		idStr := comment.ParentID.String()
		res.ParentID = &idStr
	}

	return res
}
