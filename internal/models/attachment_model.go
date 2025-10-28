package models

import "github.com/google/uuid"

type Attachment struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TaskID     uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	FileUrl    string    `json:"file_url"`
	UploadedBy uuid.UUID `json:"uploaded_by" gorm:"type:uuid;not null"`

	// relationships

	Task     Task `json:"task" gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Uploader User `json:"uploader" gorm:"foreignKey:UploadedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
