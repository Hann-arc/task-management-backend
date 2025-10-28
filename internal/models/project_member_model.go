package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RoleID    uuid.UUID `json:"role_id" gorm:"type:uuid;not null"`
	InvitedAt time.Time `json:"invited_at"`
	JoinedAt  time.Time `json:"joined_at"`

	// relationships

	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
	Role    Role    `gorm:"foreignKey:RoleID"`
}
