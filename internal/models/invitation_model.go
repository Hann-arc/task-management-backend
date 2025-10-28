package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	Email     string    `json:"email" gorm:"not null"`
	InviterID uuid.UUID `json:"inviter_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"not null;unique"`
	Status    string    `json:"status" gorm:"not null;default:'pending'"`

	CreatedAt time.Time `json:"created_at"`

	// Relationships

	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
	Inviter User    `json:"inviter" gorm:"foreignKey:InviterID"`
}
