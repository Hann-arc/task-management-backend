package models

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string    `json:"name" gorm:"not null"`

	// Relationships
	ProjectMembers []ProjectMember `json:"project_members" gorm:"foreignKey:RoleID"`
}
