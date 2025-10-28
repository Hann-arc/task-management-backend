// internal/routes/board_route.go
package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// BoardRouter sets up the routes for board-related operations
func BoardRouter(router fiber.Router) {
	boardRepo := repository.NewBoardRepository(config.DB)
	projectRepo := repository.NewProjectRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	boardService := services.NewBoardService(config.DB, boardRepo, projectRepo, activityLogService)
	boardHandler := handlers.NewBoardHandler(boardService)

	boardRoutes := router.Group("/projects/:projectId/boards", middlewares.AuthMiddleware)
	boardRoutes.Get("/", boardHandler.GetBoards)
	boardRoutes.Post("/", boardHandler.CreateBoard)
	boardRoutes.Patch("/:boardId", boardHandler.UpdateBoard)
	boardRoutes.Delete("/:boardId", boardHandler.DeleteBoard)
}
