package dto

import "time"

type ProjectMemberResponse struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"project_id"`
	User      UserBasicMember `json:"user"`
	RoleID    string          `json:"role_id"`
	InvitedAt time.Time       `json:"invited_at"`
	JoinedAt  time.Time       `json:"joined_at"`
}

type CreateProjectMemberResponse struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	InvitedAt time.Time `json:"invited_at"`
	JoinedAt  time.Time `json:"joined_at"`
}

type UserBasicMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AddMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
}
