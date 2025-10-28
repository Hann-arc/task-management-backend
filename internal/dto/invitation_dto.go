package dto

import "time"

type InvitationResponse struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token,omitempty"`
}

type CreateInvitationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}
