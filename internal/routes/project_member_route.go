package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ProjectMemberRouter sets up the routes for project member operations
func ProjectMemberRouter(router fiber.Router) {
	projectMemberRepo := repository.NewProjectMemberRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	notificationRepo := repository.NewNotificationRepository(config.DB)

	notificationService := services.NewNotificationService(notificationRepo)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	projectMemberService := services.NewProjectMemberService(projectMemberRepo, userRepo, projectRepo, activityLogService, notificationService)
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService)

	memberRoutes := router.Group("/projects/:projectId/members", middlewares.AuthMiddleware)

	memberRoutes.Post("/", projectMemberHandler.AddMember)
	memberRoutes.Get("/", projectMemberHandler.GetMembers)
	memberRoutes.Delete("/:userId", projectMemberHandler.RemoveMember)
}
