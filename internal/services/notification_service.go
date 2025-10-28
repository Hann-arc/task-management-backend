package services

import (
	"github.com/Hann-arc/task-management-backend/internal/dto"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/websocket"
	"github.com/google/uuid"
)

type NotificationService struct {
	Repo *repository.NotificationRepository
}

// NewNotificationService creates a new instance of NotificationService
func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

// CreateNotification creates a new notification and sends it via WebSocket
func (s *NotificationService) CreateNotification(
	userID, actorID uuid.UUID,
	action, referenceType string,
	relatedID uuid.UUID,
	message string,
) error {
	notification := &models.Notification{
		UserID:        userID,
		ActorID:       actorID,
		Action:        action,
		RelatedID:     relatedID,
		ReferenceType: referenceType,
		Message:       message,
		IsRead:        false,
	}

	if err := s.Repo.Create(notification); err != nil {
		return err
	}

	// Send notification via WebSocket
	go func() {
		websocket.GlobalHub.SendToUser(userID.String(), dto.NotificationResponse{
			ID:            notification.Id.String(),
			UserID:        notification.UserID.String(),
			ActorID:       notification.ActorID.String(),
			Action:        notification.Action,
			RelatedID:     notification.RelatedID.String(),
			ReferenceType: notification.ReferenceType,
			Message:       notification.Message,
			IsRead:        notification.IsRead,
			CreatedAt:     notification.CreatedAt,
		})
	}()

	return nil
}

// GetNotifications retrieves notifications for a specific user with pagination
func (s *NotificationService) GetNotifications(userID uuid.UUID, limit, offset int) ([]dto.NotificationResponse, error) {
	notifications, err := s.Repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []dto.NotificationResponse
	for _, n := range notifications {
		result = append(result, dto.NotificationResponse{
			ID:            n.Id.String(),
			UserID:        n.UserID.String(),
			ActorID:       n.ActorID.String(),
			Action:        n.Action,
			RelatedID:     n.RelatedID.String(),
			ReferenceType: n.ReferenceType,
			Message:       n.Message,
			IsRead:        n.IsRead,
			CreatedAt:     n.CreatedAt,
		})
	}

	return result, nil
}

// MarkAsRead marks a specific notification as read
func (s *NotificationService) MarkAsRead(notificationIDs []string, userID uuid.UUID) error {
	var ids []uuid.UUID
	for _, idStr := range notificationIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return s.Repo.MarkAsRead(ids, userID)
}

// MarkAllAsRead marks all user notifications as read
func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.Repo.MarkAllAsRead(userID)
}
