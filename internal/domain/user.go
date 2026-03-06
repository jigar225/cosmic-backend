package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is the persisted user entity (users table).
type User struct {
	ID                int64      `json:"id"`
	UUID              uuid.UUID  `json:"uuid"`
	Email             *string    `json:"email,omitempty"`
	PhoneNumber       *string    `json:"phone_number,omitempty"`
	PasswordHash      string     `json:"-"` // never expose
	FirstName         *string    `json:"first_name,omitempty"`
	LastName          *string    `json:"last_name,omitempty"`
	ProfilePhoto      *string    `json:"profile_photo,omitempty"`
	Role              string     `json:"role"`
	IsActive          bool       `json:"is_active"`
	IsVerified        *bool      `json:"is_verified,omitempty"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"-"`
	PreferableSubject []string   `json:"preferable_subject,omitempty"`
	PlateformVersion  string     `json:"plateform_version,omitempty"`
}

// UserCreate is input for creating a new user (signup).
type UserCreate struct {
	Email        string  `json:"email"`
	Password     string  `json:"password"` // plain; repo never stores plain
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Role         string  `json:"-"` // server default
	PlateformVer string  `json:"-"`
}

// RefreshToken is one stored session (refresh_tokens table).
type RefreshToken struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	TokenHash  string    `json:"-"`
	DeviceName *string   `json:"device_name,omitempty"`
	DeviceType *string   `json:"device_type,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	Revoked    bool      `json:"revoked"`
}
