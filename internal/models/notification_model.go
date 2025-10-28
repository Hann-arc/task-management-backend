package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id            uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ActorID       uuid.UUID `json:"actor_id" gorm:"type:uuid;not null"`
	Action        string    `json:"action" gorm:"not null"`
	RelatedID     uuid.UUID `json:"related_id" gorm:"type:uuid;not null"`
	ReferenceType string    `json:"reference_type"`
	Message       string    `json:"message"`
	IsRead        bool      `json:"is_read" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
