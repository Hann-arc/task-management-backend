package dto

import "time"

type ActivityLogResponse struct {
	ID        string                 `json:"id"`
	ProjectID string                 `json:"project_id"`
	UserID    string                 `json:"user_id"`
	User      ActivityLogUser        `json:"user"`
	Action    string                 `json:"action"`
	Details   map[string]interface{} `json:"details,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type ActivityLogUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
