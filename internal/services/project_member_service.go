package services

import (
	"errors"
	"time"

	"github.com/Hann-arc/task-management-backend/internal/dto"
	apperrors "github.com/Hann-arc/task-management-backend/internal/errors"
	"github.com/Hann-arc/task-management-backend/internal/models"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMemberService struct {
	ProjectMemberRepo   *repository.ProjectMemberRepository
	UserRepo            *repository.UserRepository
	ProjectRepo         *repository.ProjectRepository
	ActivityLogService  *ActivityLogService
	NotificationService *NotificationService
}

// NewProjectMemberService creates a new instance of ProjectMemberService
func NewProjectMemberService(
	projectMemberRepo *repository.ProjectMemberRepository,
	userRepo *repository.UserRepository,
	projectRepo *repository.ProjectRepository,
	activityLogService *ActivityLogService,
	notificationService *NotificationService,
) *ProjectMemberService {
	return &ProjectMemberService{
		ProjectMemberRepo:   projectMemberRepo,
		UserRepo:            userRepo,
		ProjectRepo:         projectRepo,
		ActivityLogService:  activityLogService,
		NotificationService: notificationService,
	}
}

// AddMember adds a new member to a project
func (s *ProjectMemberService) AddMember(projectID, ownerID uuid.UUID, req *dto.AddMemberRequest) (*dto.CreateProjectMemberResponse, error) {

	// only owner can add members
	isOwner, err := s.ProjectMemberRepo.IsOwner(projectID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrUnauthorizedOwnerOnly
	}

	invitee, err := s.ProjectMemberRepo.UserExistsByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	if invitee.ID == ownerID {
		return nil, apperrors.ErrInviteeIsOwner
	}

	// Check if already a member
	_, err = s.ProjectMemberRepo.FindByProjectAndUser(projectID, invitee.ID)
	if err == nil {
		return nil, apperrors.ErrAlreadyMember
	}

	// Get default role (it would be developed later)
	var defaultRole models.Role
	if err := s.ProjectMemberRepo.DB.Where("name = ?", "member").First(&defaultRole).Error; err != nil {
		defaultRole = models.Role{
			ID:   uuid.New(),
			Name: "member",
		}
		if err := s.ProjectMemberRepo.DB.Create(&defaultRole).Error; err != nil {
			return nil, err
		}
	}

	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    invitee.ID,
		RoleID:    defaultRole.ID,
		InvitedAt: time.Now(),
		JoinedAt:  time.Now(),
	}

	if err := s.ProjectMemberRepo.Create(member); err != nil {
		return nil, err
	}

	// Log activity and send notification

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, ownerID, "member.added", map[string]interface{}{
			"member_id": invitee.ID.String(),
			"email":     req.Email,
		})
	}

	if s.NotificationService != nil {
		go s.NotificationService.CreateNotification(
			invitee.ID,
			ownerID,
			"member.added",
			"project",
			projectID,
			"You have been added to a project",
		)
	}

	return s.buildCreateMemberResponse(member), nil
}

// GetMembers retrieves all members of a project
func (s *ProjectMemberService) GetMembers(projectID, userID uuid.UUID) ([]dto.ProjectMemberResponse, error) {
	isOwner, err := s.ProjectRepo.IsOwner(projectID, userID)
	if err != nil {
		return nil, err
	}
	isMember, err := s.ProjectRepo.IsMember(projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner && !isMember {
		return nil, apperrors.ErrUnauthorizedProject
	}

	members, err := s.ProjectMemberRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	var result []dto.ProjectMemberResponse
	for _, m := range members {
		result = append(result, *s.buildMemberResponse(&m))
	}
	return result, nil
}

// RemoveMember removes a member from a project
func (s *ProjectMemberService) RemoveMember(projectID, ownerID, targetUserID uuid.UUID) error {

	// only owner can remove members
	isOwner, err := s.ProjectMemberRepo.IsOwner(projectID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperrors.ErrUnauthorizedOwnerOnly
	}

	if ownerID == targetUserID {
		return apperrors.ErrCannotRemoveSelf
	}

	// Ensure target is a member
	_, err = s.ProjectMemberRepo.FindByProjectAndUser(projectID, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrProjectMemberNotFound
		}
		return err
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, ownerID, "member.removed", map[string]interface{}{
			"member_id": targetUserID.String(),
		})
	}

	return s.ProjectMemberRepo.Delete(projectID, targetUserID)
}

// buildMemberResponse builds a response for a project member
func (s *ProjectMemberService) buildMemberResponse(member *models.ProjectMember) *dto.ProjectMemberResponse {
	return &dto.ProjectMemberResponse{
		ID:        member.ID.String(),
		ProjectID: member.ProjectID.String(),
		User: dto.UserBasicMember{
			ID:    member.User.ID.String(),
			Name:  member.User.Name,
			Email: member.User.Email,
		},
		RoleID:    member.RoleID.String(),
		InvitedAt: member.InvitedAt,
		JoinedAt:  member.JoinedAt,
	}
}

// buildCreateMemberResponse builds a response for creating a project member
func (s *ProjectMemberService) buildCreateMemberResponse(member *models.ProjectMember) *dto.CreateProjectMemberResponse {
	return &dto.CreateProjectMemberResponse{
		ID:        member.ID.String(),
		ProjectID: member.ProjectID.String(),
		UserID:    member.UserID.String(),
		RoleID:    member.RoleID.String(),
		InvitedAt: member.InvitedAt,
		JoinedAt:  member.JoinedAt,
	}
}
