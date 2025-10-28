package routes

import (
	"os"

	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// InvitationRote sets up the routes for invitation operations
func InvitationRote(router fiber.Router) {
	invitationRepo := repository.NewInvitationRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	projectMemberRepo := repository.NewProjectMemberRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	notificationRepo := repository.NewNotificationRepository(config.DB)
	notificationService := services.NewNotificationService(notificationRepo)

	// have to change for production
	emailService := &services.NoopEmailService{}

	isDev := os.Getenv("ENV") != "production"

	invitationService := services.NewInvitationService(
		invitationRepo,
		projectRepo,
		userRepo,
		projectMemberRepo,
		emailService,
		activityLogService,
		notificationService,
		isDev,
	)

	invitationHandler := handlers.NewInvitationHandler(invitationService)

	invitationRoutes := router.Group("/invitations", middlewares.AuthMiddleware)

	invitationRoutes.Post("/projects/:projectId", invitationHandler.CreateInvitation)
	invitationRoutes.Patch("/accept", invitationHandler.AcceptInvitation)
}
