package transport

import (
	"back_testing/internal/transport/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// MaxRequestBodySize is 55MB for chapter PDF uploads (50MB file + form overhead).
const MaxRequestBodySize = 55 * 1024 * 1024

// NewApp creates the Fiber app, wires handlers, sets middleware, and registers routes.
func NewApp(h *handlers.Handlers) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: MaxRequestBodySize,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	RegisterRoutes(app, h)
	return app
}
