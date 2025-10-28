package services

import (
	"mime/multipart"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/google/uuid"
)

type AttachmentService struct {
	AttachmentRepo repository.AttachmentRepository
}

// NewAttachmentService creates a new instance of AttachmentService
func NewAttachmentService(attachmentRepo *repository.AttachmentRepository) *AttachmentService {
	return &AttachmentService{AttachmentRepo: *attachmentRepo}
}

// UploadAttachment handles the uploading of an attachment to a task
func (s *AttachmentService) UploadAttachment(taskID, userID uuid.UUID, file *multipart.FileHeader) (*dto.CreateAttachmentResponse, error) {
	isMember, err := s.AttachmentRepo.IsTaskMember(taskID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	fileUrl, err := utils.UploadToCloudinary(file)
	if err != nil {
		return nil, err
	}

	attachment := &models.Attachment{
		TaskID:     taskID,
		FileUrl:    fileUrl,
		UploadedBy: userID,
	}

	if err := s.AttachmentRepo.Create(attachment); err != nil {
		return nil, err
	}

	return &dto.CreateAttachmentResponse{
		ID:         attachment.ID.String(),
		TaskID:     attachment.TaskID.String(),
		FileUrl:    attachment.FileUrl,
		UploadedBy: attachment.UploadedBy.String(),
	}, nil
}

// GetAttachments retrieves all attachments for a given task
func (s *AttachmentService) GetAttachments(taskID, userID uuid.UUID) ([]dto.GetAttachmentResponse, error) {
	isMember, err := s.AttachmentRepo.IsTaskMember(taskID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	attachments, err := s.AttachmentRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	var result []dto.GetAttachmentResponse
	for _, att := range attachments {
		result = append(result, dto.GetAttachmentResponse{
			ID:         att.ID.String(),
			TaskID:     att.TaskID.String(),
			FileUrl:    att.FileUrl,
			UploadedBy: att.UploadedBy.String(),
			Uploader: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   att.Uploader.ID.String(),
				Name: att.Uploader.Name,
			},
		})
	}

	return result, nil
}

// DeleteAttachment handles the deletion of an attachment
func (s *AttachmentService) DeleteAttachment(attachmentID, userID uuid.UUID) error {
	attachment, err := s.AttachmentRepo.FindByID(attachmentID)
	if err != nil {
		return err
	}

	if attachment.UploadedBy == userID {
		return s.deleteAttachmentInternal(attachmentID, attachment.FileUrl)
	}

	isOwner, err := s.AttachmentRepo.IsOwnerOfAttachmentProject(attachmentID, userID)
	if err != nil {
		return err
	}
	if isOwner {
		return s.deleteAttachmentInternal(attachmentID, attachment.FileUrl)
	}

	return apperrors.ErrUnauthorizedOwnerOnly
}

// deleteAttachmentInternal performs the actual deletion of the attachment from cloud storage and database
func (s *AttachmentService) deleteAttachmentInternal(attachmentID uuid.UUID, fileUrl string) error {

	if err := utils.DeleteFromCloudinary(fileUrl); err != nil {
		return err
	}

	return s.AttachmentRepo.Delete(attachmentID)
}
