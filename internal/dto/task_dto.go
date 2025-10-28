package dto

import (
	"time"
)

type TaskResponse struct {
	ID          string     `json:"id"`
	BoardID     string     `json:"board_id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	Assignee    *UserBasic `json:"assignee,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Labels      []LabelDTO `json:"labels,omitempty"`
}

type UserBasic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LabelDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high urgent"`
	DueDate     *string    `json:"due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	AssigneeID  *string    `json:"assignee_id" validate:"omitempty,uuid"`
	Labels      []LabelDTO `json:"labels" validate:"dive"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Priority    *string     `json:"priority,omitempty"`
	DueDate     *string     `json:"due_date,omitempty"`
	AssigneeID  *string     `json:"assignee_id,omitempty"`
	Labels      *[]LabelDTO `json:"labels,omitempty"`
}
