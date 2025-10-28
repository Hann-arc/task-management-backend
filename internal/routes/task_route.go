package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// TaskRouter sets up the routes for task operations
func TaskRouter(router fiber.Router) {
	taskRepo := repository.NewTaskRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	notificationRepo := repository.NewNotificationRepository(config.DB)

	notificationService := services.NewNotificationService(notificationRepo)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	taskService := services.NewTaskService(taskRepo, projectRepo, userRepo, activityLogService, notificationService)
	taskHandler := handlers.NewTaskHandler(taskService)

	taskRoutes := router.Group("/boards/:boardId/tasks", middlewares.AuthMiddleware)

	taskRoutes.Post("/", taskHandler.CreateTask)
	taskRoutes.Get("/", taskHandler.GetTasksByBoard)

	taskRoutes.Patch("/:id", taskHandler.UpdateTask)
	taskRoutes.Delete("/:id", taskHandler.DeleteTask)

}
