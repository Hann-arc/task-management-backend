package repository

import (
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

// NewNotificationRepository creates a new instance of NotificationRepository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// Create saves a new notification to the database
func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

// FindByUserID retrieves notifications by user ID with pagination
func (r *NotificationRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a specific notification as read
func (r *NotificationRepository) MarkAsRead(notificationIDs []uuid.UUID, userID uuid.UUID) error {
	return r.DB.Model(&models.Notification{}).
		Where("id IN ? AND user_id = ?", notificationIDs, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead marks all user notifications as read
func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
