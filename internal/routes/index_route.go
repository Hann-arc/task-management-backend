package routes

import "github.com/gofiber/fiber/v2"

func MainRoutes(app *fiber.App) {
	api := app.Group("/v1/api")

	AuthRoutes(api)
	UserRoutes(api)
	ProjectRoutes(api)
	BoardRouter(api)
	TaskRouter(api)
	ProjectMemberRouter(api)
	InvitationRote(api)
	CommentRoute(api)
	ActivityLogRoutes(api)
	AttachmentRoutes(api)
	NotificationRoutes(api)
}
