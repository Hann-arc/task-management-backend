package middlewares

import (
	"strings"

	"github.com/Hann-arc/task-management-backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// Middleware for authenticating requests using JWT tokens
func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token not provided",
		})
	}

	claims, err := utils.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}
