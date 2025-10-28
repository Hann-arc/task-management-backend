package services

import (
	"errors"
	"time"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskService struct {
	TaskRepo            *repository.TaskRepository
	ProjectRepo         *repository.ProjectRepository
	UserRepo            *repository.UserRepository
	ActivityLogService  *ActivityLogService
	NotificationService *NotificationService
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(
	taskRepo *repository.TaskRepository,
	projectRepo *repository.ProjectRepository,
	userRepo *repository.UserRepository,
	activityLogService *ActivityLogService,
	notificationService *NotificationService,
) *TaskService {
	return &TaskService{
		TaskRepo:            taskRepo,
		ProjectRepo:         projectRepo,
		UserRepo:            userRepo,
		ActivityLogService:  activityLogService,
		NotificationService: notificationService,
	}
}

// CreateTask handles the creation of a new task within a board
func (s *TaskService) CreateTask(req *dto.CreateTaskRequest, boardID, userID uuid.UUID) (*dto.TaskResponse, error) {

	// validate priority
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
		"urgent": true,
	}

	if !validPriorities[req.Priority] {
		return nil, apperrors.ErrInvalidTaskData
	}

	//  Validate if board exists and user is a member of the project
	boardExists, err := s.TaskRepo.BoardExists(boardID)
	if err != nil {
		return nil, err
	}
	if !boardExists {
		return nil, apperrors.ErrBoardNotFound
	}

	var board models.Board
	if err := s.TaskRepo.DB.Select("project_id").Where("id = ?", boardID).First(&board).Error; err != nil {
		return nil, err
	}

	projectID := board.ProjectID

	isMember, err := s.ProjectRepo.IsMember(projectID, userID)
	if err != nil {
		return nil, err
	}
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember && !isOwner {
		return nil, apperrors.ErrUnauthorizedTask
	}

	//  Validate assignee if provided
	var assigneeID *uuid.UUID
	if req.AssigneeID != nil {
		id, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			return nil, apperrors.ErrInvalidTaskData
		}
		exists, err := s.TaskRepo.UserExists(id)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apperrors.ErrAssigneeNotFound
		}
		assigneeID = &id
	}

	// Parse due date
	var dueDate time.Time
	if req.DueDate != nil {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return nil, apperrors.ErrInvalidTaskData
		}
		dueDate = t
	}

	task := &models.Task{
		ID:          uuid.New(),
		BoardID:     boardID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     dueDate,
		AssigneeID:  assigneeID,
		CreatedBy:   userID,
	}

	if err := s.TaskRepo.Create(task); err != nil {
		return nil, err
	}

	// Handle labels
	if len(req.Labels) > 0 {
		var labels []models.TaskLabel
		for _, l := range req.Labels {
			labels = append(labels, models.TaskLabel{
				ID:     uuid.New(),
				TaskID: task.ID,
				Name:   l.Name,
				Color:  l.Color,
			})
		}
		if err := s.TaskRepo.ReplaceLabels(task.ID, labels); err != nil {
			return nil, err
		}
	}

	// Log activity
	if s.ActivityLogService != nil {
		details := map[string]interface{}{
			"task_id":    task.ID.String(),
			"title":      req.Title,
			"board_id":   boardID.String(),
			"project_id": projectID.String(),
		}
		if req.AssigneeID != nil {
			details["assignee_id"] = *req.AssigneeID
		}
		s.ActivityLogService.LogActivity(projectID, userID, "task.created", details)
	}

	// Send notification to assignee
	if s.NotificationService != nil && req.AssigneeID != nil {
		assigneeID, _ := uuid.Parse(*req.AssigneeID)
		if assigneeID != userID {
			go s.NotificationService.CreateNotification(
				assigneeID,
				userID,
				"task.assigned",
				"task",
				task.ID,
				"You have been assigned to a new task",
			)
		}
	}

	return s.buildTaskResponse(task), nil
}

// GetTasksByBoard retrieves all tasks for a specific board, ensuring the user has access
func (s *TaskService) GetTasksByBoard(boardID, userID uuid.UUID) ([]dto.TaskResponse, error) {

	var board models.Board
	if err := s.TaskRepo.DB.Select("project_id").Where("id = ?", boardID).First(&board).Error; err != nil {
		return nil, err
	}

	projectID := board.ProjectID

	isMember, err := s.ProjectRepo.IsMember(projectID, userID)
	if err != nil {
		return nil, err
	}
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember && !isOwner {
		return nil, apperrors.ErrUnauthorizedTask
	}

	tasks, err := s.TaskRepo.FindByBoardID(boardID)
	if err != nil {
		return nil, err
	}

	var result []dto.TaskResponse
	for _, t := range tasks {
		result = append(result, *s.buildTaskResponse(&t))
	}
	return result, nil
}

// UpdateTask modifies the details of an existing task
func (s *TaskService) UpdateTask(taskID, userID uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, err
	}

	// Validate access
	var board models.Board
	if err := s.TaskRepo.DB.Model(&models.Board{}).
		Select("project_id").
		Where("id = ?", task.BoardID).
		First(&board).Error; err != nil {
		return nil, err
	}

	projectID := board.ProjectID

	isMember, err := s.ProjectRepo.IsMember(projectID, userID)
	if err != nil {
		return nil, err
	}
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember && !isOwner {
		return nil, apperrors.ErrUnauthorizedTask
	}

	data := map[string]interface{}{}

	if req.Title != nil {
		data["title"] = *req.Title
	}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if req.Priority != nil {
		data["priority"] = *req.Priority
	}

	if req.DueDate != nil {
		if *req.DueDate == "" {
			data["due_date"] = nil
		} else {
			t, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return nil, apperrors.ErrInvalidTaskData
			}
			data["due_date"] = t
		}
	}

	if req.AssigneeID != nil {
		if *req.AssigneeID == "" {
			data["assignee_id"] = nil
		} else {
			id, err := uuid.Parse(*req.AssigneeID)
			if err != nil {
				return nil, apperrors.ErrInvalidTaskData
			}
			exists, err := s.TaskRepo.UserExists(id)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, apperrors.ErrAssigneeNotFound
			}
			data["assignee_id"] = id
		}
	}

	if len(data) == 0 && req.Labels == nil {
		return nil, apperrors.ErrInvalidTaskData
	}

	if len(data) > 0 {
		if err := s.TaskRepo.Update(taskID, data); err != nil {
			return nil, err
		}
	}

	if req.Labels != nil {
		var labels []models.TaskLabel
		for _, l := range *req.Labels {
			labels = append(labels, models.TaskLabel{
				ID:     uuid.New(),
				TaskID: taskID,
				Name:   l.Name,
				Color:  l.Color,
			})
		}
		if err := s.TaskRepo.ReplaceLabels(taskID, labels); err != nil {
			return nil, err
		}
	}

	updatedTask, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// Log activity
	if s.ActivityLogService != nil {
		details := map[string]interface{}{
			"task_id": taskID.String(),
		}

		if req.Title != nil {
			details["title"] = *req.Title
		}

		if req.Description != nil {
			details["descripton"] = *req.Description
		}

		if req.DueDate != nil {
			details["due_date"] = *req.DueDate
		}

		if req.AssigneeID != nil {
			details["assignee_id"] = *req.AssigneeID
		}

		if req.Priority != nil {
			details["priority"] = *req.Priority
		}

		s.ActivityLogService.LogActivity(projectID, userID, "task.updated", details)
	}

	// Send notification to assignee
	if s.NotificationService != nil && req.AssigneeID != nil {
		oldAssignee := task.AssigneeID
		newAssignee, _ := uuid.Parse(*req.AssigneeID)

		// Notify to new assignee if changed
		if newAssignee != userID && (oldAssignee == nil || *oldAssignee != newAssignee) {
			go s.NotificationService.CreateNotification(
				newAssignee,
				userID,
				"task.assigned",
				"task",
				taskID,
				"You have been assigned to a task",
			)
		}
	}

	return s.buildTaskResponse(updatedTask), nil
}

// DeleteTask performs a soft delete of a task after validating user access
func (s *TaskService) DeleteTask(taskID, userID uuid.UUID) error {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrTaskNotFound
		}
		return err
	}

	// Validate access
	var board models.Board
	if err := s.TaskRepo.DB.Model(&models.Board{}).
		Select("project_id").
		Where("id = ?", task.BoardID).
		First(&board).Error; err != nil {
		return err
	}

	projectID := board.ProjectID

	isMember, err := s.ProjectRepo.IsMember(projectID, userID)
	if err != nil {
		return err
	}
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return err
	}
	if !isMember && !isOwner {
		return apperrors.ErrUnauthorizedTask
	}

	if err := s.TaskRepo.SoftDelete(taskID); err != nil {
		return err
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, userID, "task.deleted", map[string]interface{}{
			"task_id": taskID.String(),
		})
	}

	return nil
}

// Helper to converts a task model to a task response DTO
func (s *TaskService) buildTaskResponse(task *models.Task) *dto.TaskResponse {
	resp := &dto.TaskResponse{
		ID:          task.ID.String(),
		BoardID:     task.BoardID.String(),
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		CreatedBy:   task.CreatedBy.String(),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		DeletedAt:   utils.ToTimePtr(task.DeletedAt),
	}

	if task.DueDate != (time.Time{}) {
		resp.DueDate = &task.DueDate
	}

	if task.AssigneeID != nil {
		idStr := task.AssigneeID.String()
		resp.AssigneeID = &idStr
		if task.Assignee.ID != uuid.Nil {
			resp.Assignee = &dto.UserBasic{
				ID:   task.Assignee.ID.String(),
				Name: task.Assignee.Name,
			}
		}
	}

	for _, l := range task.Labels {
		resp.Labels = append(resp.Labels, dto.LabelDTO{
			ID:    l.ID.String(),
			Name:  l.Name,
			Color: l.Color,
		})
	}

	return resp
}
