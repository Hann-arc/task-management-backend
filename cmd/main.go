package main

import (
	"log"

	"github.com/Hann-arc/task-management-backend/config"
	"github.com/Hann-arc/task-management-backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnnDB()
	config.SetupCloudinary()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hallo, world!")
	})

	routes.MainRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
