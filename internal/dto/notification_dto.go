package dto

import "time"

type NotificationResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ActorID       string    `json:"actor_id"`
	ActorName     string    `json:"actor_name"`
	Action        string    `json:"action"`
	RelatedID     string    `json:"related_id"`
	ReferenceType string    `json:"reference_type"`
	Message       string    `json:"message"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}

type MarkAsReadRequest struct {
	IDs []string `json:"ids"`
}
