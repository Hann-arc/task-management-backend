package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ActivityLogRoutes sets up the routes for activity log operations
func ActivityLogRoutes(router fiber.Router) {
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	activityLogHandler := handlers.NewActivityLogHandler(activityLogService)

	logRoutes := router.Group("/projects/:projectId/activity-logs", middlewares.AuthMiddleware)
	logRoutes.Get("/", activityLogHandler.GetActivityLogs)
}
