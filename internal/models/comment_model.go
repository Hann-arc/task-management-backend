package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID       uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TaskID   uuid.UUID  `json:"task_id" gorm:"not null;type:uuid"`
	UserID   uuid.UUID  `json:"user_id" gorm:"not null;type:uuid"`
	ParentID *uuid.UUID `json:"parent_id,omitempty" gorm:"type:uuid;index"`
	Content  string     `json:"content" gorm:"not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	Task    Task      `json:"task" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User    User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
