package routes

import (
	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/handlers"
	"github.com/Hann-arc/task-management-backend/internal/repository"
	"github.com/Hann-arc/task-management-backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes sets up authentication routes

func AuthRoutes(router fiber.Router) {

	userRepo := repository.NewUserRepository(config.DB)
	authService := services.NewAuthservice(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := router.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
}
