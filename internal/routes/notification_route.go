package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/Hann-arc/task-management-backend/internal/websocket"
	"github.com/gofiber/fiber/v2"
)

// NotificationRoutes sets up the routes for notification operations
func NotificationRoutes(router fiber.Router) {
	notificationRepo := repository.NewNotificationRepository(config.DB)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	notificationRoutes := router.Group("/notifications", middlewares.AuthMiddleware)
	notificationRoutes.Get("/", notificationHandler.GetNotifications)
	notificationRoutes.Patch("/read", notificationHandler.MarkAsRead)
	notificationRoutes.Patch("/read-all", notificationHandler.MarkAllAsRead)

	// WebSocket route
	router.Get("/ws/notifications", websocket.Route())
}
