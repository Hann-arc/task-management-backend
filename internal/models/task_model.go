package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BoardID     uuid.UUID      `json:"board_id" gorm:"type:uuid;not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Priority    string         `json:"priority" gorm:"not null"`
	DueDate     time.Time      `json:"due_date"`
	AssigneeID  *uuid.UUID     `json:"assignee_id,omitempty" gorm:"type:uuid"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	Board       Board        `json:"board" gorm:"foreignKey:BoardID;references:ID"`
	Assignee    User         `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID;references:ID"`
	Creator     User         `json:"creator" gorm:"foreignKey:CreatedBy;references:ID"`
	Labels      []TaskLabel  `json:"labels" gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments    []Comment    `json:"comments" gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachments []Attachment `json:"attachments" gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
