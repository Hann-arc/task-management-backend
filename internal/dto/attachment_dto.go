package dto

type GetAttachmentResponse struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	FileUrl    string `json:"file_url"`
	UploadedBy string `json:"uploaded_by"`
	Uploader   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"uploader"`
}

type CreateAttachmentResponse struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	FileUrl    string `json:"file_url"`
	UploadedBy string `json:"uploaded_by"`
}

type UploadAttachmentRequest struct {
	FileUrl string `json:"file_url" validate:"required,url"`
}
