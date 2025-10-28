package dto

import "time"

type CommentResponse struct {
	ID        string      `json:"id"`
	TaskID    string      `json:"task_id"`
	UserID    string      `json:"user_id"`
	ParentID  *string     `json:"parent_id,omitempty"`
	Content   string      `json:"content"`
	User      CommentUser `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CommentUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}
