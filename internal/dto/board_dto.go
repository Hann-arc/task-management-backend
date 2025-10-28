package dto

import "time"

type BoardResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	OrderIndex int       `json:"order_index"`
	ProjectID  string    `json:"project_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
