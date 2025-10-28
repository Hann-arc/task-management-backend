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

type InvitationService struct {
	InvitationRepo      *repository.InvitationRepository
	ProjectRepo         *repository.ProjectRepository
	UserRepo            *repository.UserRepository
	ProjectMemberRepo   *repository.ProjectMemberRepository
	EmailService        EmailService
	ActivityLogService  *ActivityLogService
	NotificationService *NotificationService
	IsDev               bool
}

// NewInvitationService creates a new instance of InvitationService
func NewInvitationService(
	invitationRepo *repository.InvitationRepository,
	projectRepo *repository.ProjectRepository,
	userRepo *repository.UserRepository,
	projectMemberRepo *repository.ProjectMemberRepository,
	emailService EmailService,
	activityLogService *ActivityLogService,
	notificationService *NotificationService,
	isDev bool,
) *InvitationService {
	return &InvitationService{
		InvitationRepo:      invitationRepo,
		ProjectRepo:         projectRepo,
		UserRepo:            userRepo,
		ProjectMemberRepo:   projectMemberRepo,
		EmailService:        emailService,
		IsDev:               isDev,
		ActivityLogService:  activityLogService,
		NotificationService: notificationService,
	}
}

// CreateInvitation creates a new invitation for a user to join a project
func (s *InvitationService) CreateInvitation(projectID, inviterID uuid.UUID, email string) (*dto.InvitationResponse, error) {

	// only owner can invite
	isOwner, err := s.ProjectRepo.IsOwner(projectID, inviterID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperrors.ErrUnauthorizedOwnerOnly
	}

	inviter, err := s.UserRepo.FindByID(inviterID)
	if err != nil {
		return nil, err
	}
	if inviter.Email == email {
		return nil, apperrors.ErrCannotInviteSelf
	}

	user, _ := s.UserRepo.FindByEmail(email)
	if user != nil {
		isMember, err := s.InvitationRepo.IsMember(projectID, user.ID)
		if err != nil {
			return nil, err
		}
		if isMember {
			return nil, apperrors.ErrAlreadyMember
		}
	}

	// _, err = s.InvitationRepo.FindPendingByProjectAndEmail(projectID, email)
	// if err == nil {
	//
	//
	// }

	token := uuid.New().String()

	invitation := &models.Invitation{
		ID:        uuid.New(),
		ProjectID: projectID,
		Email:     email,
		InviterID: inviterID,
		Token:     token,
		Status:    "pending",
	}

	if err := s.InvitationRepo.Create(invitation); err != nil {
		return nil, err
	}

	// for dev purpose, return token in response
	resp := &dto.InvitationResponse{
		ID:        invitation.ID.String(),
		ProjectID: invitation.ProjectID.String(),
		Email:     invitation.Email,
		Status:    invitation.Status,
		CreatedAt: invitation.CreatedAt,
	}

	// send email invitation for production (it will not send email in dev mode)
	project, _ := s.ProjectRepo.FindByID(projectID)
	if project != nil {
		s.EmailService.SendInvitation(email, token, project.Name)
	}

	if s.IsDev {
		resp.Token = token
	}

	// Log activity
	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(projectID, inviterID, "invitation.sent", map[string]interface{}{
			"email": email,
		})
	}

	user, err = s.UserRepo.FindByEmail(email)
	if err == nil && s.NotificationService != nil {
		go s.NotificationService.CreateNotification(
			user.ID,
			inviterID,
			"invitation.sent",
			"project",
			projectID,
			"You have been invited to a project",
		)
	}

	return resp, nil
}

// AcceptInvitation allows a user to accept an invitation to join a project
func (s *InvitationService) AcceptInvitation(token string, userID uuid.UUID) error {
	invitation, err := s.InvitationRepo.FindByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrInvitationNotFound
		}
		return err
	}

	if invitation.Status != "pending" {
		return apperrors.ErrInvitationUsed
	}

	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user.Email != invitation.Email {
		return apperrors.ErrInvitationNotFound
	}

	isMember, err := s.InvitationRepo.IsMember(invitation.ProjectID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return apperrors.ErrAlreadyMember
	}

	defaultRole, err := s.getMemberRole()
	if err != nil {
		return err
	}

	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: invitation.ProjectID,
		UserID:    userID,
		RoleID:    defaultRole.ID,
		InvitedAt: time.Now(),
		JoinedAt:  time.Now(),
	}

	if err := s.ProjectMemberRepo.Create(member); err != nil {
		return err
	}

	// Log activity

	if s.ActivityLogService != nil {
		s.ActivityLogService.LogActivity(invitation.ProjectID, userID, "invitation.accepted", map[string]interface{}{
			"invitation_id": invitation.ID.String(),
			"email":         invitation.Email,
		})
	}

	//  Notification to inviter
	if s.NotificationService != nil {
		go s.NotificationService.CreateNotification(
			invitation.InviterID,
			userID,
			"invitation.accepted",
			"project",
			invitation.ProjectID,
			"Your invitation has been accepted",
		)
	}

	return s.InvitationRepo.UpdateStatus(invitation.ID, "accepted")
}

// getMemberRole retrieves the default member role
func (s *InvitationService) getMemberRole() (*models.Role, error) {
	var role models.Role
	if err := s.ProjectRepo.DB.Where("name = ?", "member").First(&role).Error; err != nil {
		role = models.Role{ID: uuid.New(), Name: "member"}
		if err := s.ProjectRepo.DB.Create(&role).Error; err != nil {
			return nil, err
		}
	}
	return &role, nil
}
