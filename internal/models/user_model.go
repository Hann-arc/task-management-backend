package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	AvatarUrl    string    `json:"avatar_url"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	OwenedProjects []Project       `json:"owned_projects" gorm:"foreignKey:OwnerID;references:ID"`
	ProjectMembers []ProjectMember `json:"project_members" gorm:"foreignKey:UserID;references:ID"`
	TasksAssigned  []Task          `json:"task_assigned" gorm:"foreignKey:AssigneeID;references:ID"`
	TasksCreated   []Task          `json:"task_created" gorm:"foreignKey:CreatedBy;references:ID"`
	Comments       []Comment       `json:"comments" gorm:"foreignKey:UserID;references:ID"`
	Attachments    []Attachment    `json:"attachments" gorm:"foreignKey:UploadedBy;references:ID"`
	Notifications  []Notification  `json:"notifications" gorm:"foreignKey:UserID;references:ID"`
	ActivityLogs   []ActivityLog   `json:"activity_logs" gorm:"foreignKey:UserID;references:ID"`
	Invitation     []Invitation    `json:"invitations" gorm:"foreignKey:InviterID;references:ID"`
}
