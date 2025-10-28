package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AttachmentRoutes sets up the routes for attachment operations
func AttachmentRoutes(router fiber.Router) {
	attachmentRepo := repository.NewAttachmentRepository(config.DB)
	attachmentService := services.NewAttachmentService(attachmentRepo)
	attachmentHandler := handlers.NewAttachmentHandler(attachmentService)

	attachmentRoutes := router.Group("/tasks/:taskId/attachments", middlewares.AuthMiddleware)
	attachmentRoutes.Post("/", attachmentHandler.UploadAttachment)
	attachmentRoutes.Get("/", attachmentHandler.GetAttachments)
	attachmentRoutes.Delete("/:id", attachmentHandler.DeleteAttachment)
}
