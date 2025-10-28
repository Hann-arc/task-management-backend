package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// ProjectRoutes sets up the routes for project operations
func ProjectRoutes(router fiber.Router) {
	projectRepo := repository.NewProjectRepository(config.DB)
	userRepo := repository.NewUserRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)

	activityLogService := services.NewActivityLogService(activityLogRepo)
	projectService := services.NewProjectService(projectRepo, userRepo, activityLogService)
	projectHandler := handlers.NewProjectHandler(projectService)

	projectRoute := router.Group("/projects", middlewares.AuthMiddleware)
	projectRoute.Post("/", projectHandler.CreateProject)
	projectRoute.Get("/", projectHandler.ListProjects)
	projectRoute.Get("/:id", projectHandler.GetProject)
	projectRoute.Patch("/:id", projectHandler.UpdateProject)
	projectRoute.Delete("/:id", projectHandler.DeleteProject)
}
