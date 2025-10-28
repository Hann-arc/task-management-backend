package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Action    string    `json:"action" gorm:"not null"`
	Details   []byte    `json:"details" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at"`

	// Relationships

	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}
