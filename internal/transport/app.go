package transport

import (
	"back_testing/internal/transport/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// NewApp creates the Fiber app, wires handlers, sets middleware, and registers routes.
func NewApp(h *handlers.Handlers) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	RegisterRoutes(app, h)
	return app
}
