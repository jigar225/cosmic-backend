package handlers

import (
	"time"

	"back_testing/internal/auth"
	"back_testing/internal/domain"
	"back_testing/internal/repository"
	"back_testing/internal/transport/middleware"

	"github.com/gofiber/fiber/v2"
)

// AuthConfig holds token signing and expiry for auth handlers.
type AuthConfig struct {
	JWTSecret         []byte
	AccessExpiry      time.Duration
	RefreshExpiry     time.Duration
}

// SignupRequest is the JSON body for POST /auth/signup.
type SignupRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest is the JSON body for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest is the JSON body for POST /auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokensResponse is the JSON response for login/signup/refresh.
type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// RequireAdmin ensures the authenticated user has role=admin.
// It assumes RequireAuth has already run and set user_id in Locals.
func (h *Handlers) RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if h.UserRepo == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
		}
		userID := middleware.UserIDFromContext(c)
		if userID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		user, err := h.UserRepo.GetByID(c.Context(), userID)
		if err != nil {
			if err == repository.ErrUserNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin access required"})
		}
		return c.Next()
	}
}

func (h *Handlers) Signup(c *fiber.Ctx) error {
	if h.UserRepo == nil || h.RefreshTokenRepo == nil || h.AuthConfig == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	var req SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}
	if len(req.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password must be at least 8 characters"})
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}
	user, err := h.UserRepo.Create(c.Context(), domain.UserCreate{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, passwordHash)
	if err != nil {
		if err == repository.ErrUserConflict {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.issueTokens(c, user.ID, nil)
}

func (h *Handlers) Login(c *fiber.Ctx) error {
	if h.UserRepo == nil || h.RefreshTokenRepo == nil || h.AuthConfig == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}

	user, err := h.UserRepo.GetByEmail(c.Context(), req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "account disabled"})
	}
	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	_ = h.UserRepo.UpdateLastLoginAt(c.Context(), user.ID)
	return h.issueTokens(c, user.ID, deviceInfo(c))
}

func (h *Handlers) Refresh(c *fiber.Ctx) error {
	if h.RefreshTokenRepo == nil || h.AuthConfig == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	tokenHash, err := auth.ValidateAndHashRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token required"})
	}

	rt, err := h.RefreshTokenRepo.GetByTokenHash(c.Context(), tokenHash)
	if err != nil {
		if err == repository.ErrRefreshTokenNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired refresh token"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Rotate: delete old token, create new one
	_ = h.RefreshTokenRepo.DeleteByTokenHash(c.Context(), tokenHash)
	return h.issueTokens(c, rt.UserID, deviceInfo(c))
}

func (h *Handlers) Logout(c *fiber.Ctx) error {
	if h.RefreshTokenRepo == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	var req LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	tokenHash, err := auth.ValidateAndHashRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token required"})
	}
	_ = h.RefreshTokenRepo.RevokeByTokenHash(c.Context(), tokenHash)
	return c.JSON(fiber.Map{"message": "logged out"})
}

func (h *Handlers) LogoutAll(c *fiber.Ctx) error {
	if h.RefreshTokenRepo == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	userID := middleware.UserIDFromContext(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if err := h.RefreshTokenRepo.RevokeAllForUser(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "logged out from all devices"})
}

func (h *Handlers) Me(c *fiber.Ctx) error {
	if h.UserRepo == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}
	userID := middleware.UserIDFromContext(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	user, err := h.UserRepo.GetByID(c.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(userResponse(user))
}

// issueTokens creates access + refresh tokens, stores refresh token hash, returns both.
func (h *Handlers) issueTokens(c *fiber.Ctx, userID int64, device *deviceInfoOpts) error {
	accessToken, err := auth.CreateAccessToken(h.AuthConfig.JWTSecret, userID, h.AuthConfig.AccessExpiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create access token"})
	}
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create refresh token"})
	}
	tokenHash := auth.HashRefreshToken(refreshToken)
	expiresAt := time.Now().Add(h.AuthConfig.RefreshExpiry)

	var deviceName, deviceType, ipAddress *string
	if device != nil {
		deviceName = device.DeviceName
		deviceType = device.DeviceType
		ipAddress = device.IPAddress
	}
	_, err = h.RefreshTokenRepo.Create(c.Context(), userID, tokenHash, expiresAt, deviceName, deviceType, ipAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to store session"})
	}

	return c.JSON(TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.AuthConfig.AccessExpiry.Seconds()),
	})
}

type deviceInfoOpts struct {
	DeviceName *string
	DeviceType *string
	IPAddress  *string
}

func deviceInfo(c *fiber.Ctx) *deviceInfoOpts {
	dn := c.Get("X-Device-Name")
	dt := c.Get("X-Device-Type")
	ip := c.IP()
	var d deviceInfoOpts
	if dn != "" {
		d.DeviceName = &dn
	}
	if dt != "" {
		d.DeviceType = &dt
	}
	if ip != "" {
		d.IPAddress = &ip
	}
	return &d
}

func userResponse(u domain.User) fiber.Map {
	return fiber.Map{
		"id":          u.ID,
		"uuid":        u.UUID,
		"email":       u.Email,
		"first_name":  u.FirstName,
		"last_name":   u.LastName,
		"profile_photo": u.ProfilePhoto,
		"role":        u.Role,
		"is_active":   u.IsActive,
		"created_at":  u.CreatedAt,
	}
}
