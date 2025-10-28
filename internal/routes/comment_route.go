package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// CommentRoute sets up the routes for comment-related operations
func CommentRoute(router fiber.Router) {
	commentRepo := repository.NewCommentRepository(config.DB)
	taskRepo := repository.NewTaskRepository(config.DB)
	activityLogRepo := repository.NewActivityLogRepository(config.DB)
	notificationRepo := repository.NewNotificationRepository(config.DB)

	notificationService := services.NewNotificationService(notificationRepo)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	commentService := services.NewCommentService(commentRepo, taskRepo, activityLogService, notificationService)
	commentHandler := handlers.NewCommentHandler(commentService)

	commentRoutes := router.Group("/tasks/:taskId/comments", middlewares.AuthMiddleware)
	commentRoutes.Post("/", commentHandler.CreateMainComment)
	commentRoutes.Get("/", commentHandler.GetComments)
	commentRoutes.Delete("/:id", commentHandler.DeleteComment)

	replyRoutes := router.Group("/comments", middlewares.AuthMiddleware)
	replyRoutes.Post("/:commentId/replies", commentHandler.CreateReply)
}
