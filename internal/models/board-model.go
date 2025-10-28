package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name       string         `json:"name"`
	OrderIndex int            `json:"order_index"`
	ProjectID  uuid.UUID      `json:"project_id" gorm:"type:uuid;not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships

	Project Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tasks   []Task  `json:"tasks" gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
