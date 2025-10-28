package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/middlewares"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// UserRoutes sets up user-related routes
func UserRoutes(router fiber.Router) {
	userRepo := repository.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	usersRoute := router.Group("/users", middlewares.AuthMiddleware)
	usersRoute.Get("/me", userHandler.GetProfile)
	usersRoute.Patch("/me", userHandler.UpdateProfile)

	// usersRoute.Get("/", userHandler.GetAllUsers)
}
