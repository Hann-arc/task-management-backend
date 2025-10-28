package models

import "github.com/google/uuid"

type TaskLabel struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TaskID uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	Name   string    `json:"name" gorm:"not null"`
	Color  string    `json:"color" gorm:"not null"`

	// Relationships
	Task Task `json:"task" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
