package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:uuid;not null"`
	Description string    `json:"description"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships

	Owner        User            `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Members      []ProjectMember `json:"members" gorm:"foreignKey:ProjectID"`
	Boards       []Board         `json:"boards" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActivityLogs []ActivityLog   `json:"activity_logs" gorm:"foreignKey:ProjectID"`
	Invitations  []Invitation    `json:"invitations" gorm:"foreignKey:ProjectID"`
}
