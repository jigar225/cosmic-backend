package middleware

import (
	"strings"

	"back_testing/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// AuthConfig holds JWT secret and optional settings for the auth middleware.
type AuthConfig struct {
	JWTSecret []byte
}

// RequireAuth returns a Fiber handler that validates the Bearer access token and sets user_id in Locals.
// If missing or invalid, responds with 401 and does not call next.
func RequireAuth(cfg AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(cfg.JWTSecret) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "auth not configured"})
		}
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Bearer token required"})
		}
		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token required"})
		}
		userID, err := auth.ValidateAccessToken(cfg.JWTSecret, tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}

// UserIDFromContext returns the user_id set by RequireAuth. Must be used only after RequireAuth.
func UserIDFromContext(c *fiber.Ctx) int64 {
	v := c.Locals("user_id")
	if v == nil {
		return 0
	}
	id, _ := v.(int64)
	return id
}
