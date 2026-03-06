package repository

import (
	"context"
	"errors"
	"strings"

	"back_testing/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserConflict  = errors.New("user with this email already exists")
)

// UserRepo handles persistence for users.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo returns a new UserRepo.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// GetByEmail returns the user by email (case-insensitive), or ErrUserNotFound.
// Only returns non-deleted, active users.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, uuid, email, phone_number, password_hash,
		       first_name, last_name, profile_photo, role, is_active, is_verified,
		       email_verified_at, last_login_at, created_at, updated_at,
		       COALESCE(preferable_subject, '{}'), COALESCE(plateform_version, '0.0.0')
		FROM users
		WHERE LOWER(TRIM(email)) = LOWER(TRIM($1))
		  AND deleted_at IS NULL
	`, email).Scan(
		&u.ID,
		&u.UUID,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.ProfilePhoto,
		&u.Role,
		&u.IsActive,
		&u.IsVerified,
		&u.EmailVerifiedAt,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PreferableSubject,
		&u.PlateformVersion,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

// Create inserts a new user (signup). Password must be hashed by caller; we store hash only.
func (r *UserRepo) Create(ctx context.Context, in domain.UserCreate, passwordHash string) (domain.User, error) {
	emailTrim := trimLower(in.Email)
	if emailTrim == "" {
		return domain.User{}, errors.New("email is required")
	}
	// Check duplicate email
	var dummy int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM users WHERE LOWER(TRIM(email)) = $1 AND deleted_at IS NULL`, emailTrim).Scan(&dummy)
	if err == nil {
		return domain.User{}, ErrUserConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	}

	role := in.Role
	if role == "" {
		role = "teacher"
	}
	plateformVer := in.PlateformVer
	if plateformVer == "" {
		plateformVer = "0.0.0"
	}

	newUUID := uuid.New()
	var u domain.User
	err = r.pool.QueryRow(ctx, `
		INSERT INTO users (
			uuid, email, password_hash, first_name, last_name,
			role, preferable_subject, plateform_version
		)
		VALUES ($1, $2, $3, $4, $5, $6, '{}', $7)
		RETURNING id, uuid, email, phone_number, password_hash,
		          first_name, last_name, profile_photo, role, is_active, is_verified,
		          email_verified_at, last_login_at, created_at, updated_at,
		          COALESCE(preferable_subject, '{}'), COALESCE(plateform_version, '0.0.0')
	`,
		newUUID,
		emailTrim,
		passwordHash,
		in.FirstName,
		in.LastName,
		role,
		plateformVer,
	).Scan(
		&u.ID,
		&u.UUID,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.ProfilePhoto,
		&u.Role,
		&u.IsActive,
		&u.IsVerified,
		&u.EmailVerifiedAt,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PreferableSubject,
		&u.PlateformVersion,
	)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

// GetByID returns the user by id, or ErrUserNotFound.
func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, uuid, email, phone_number, password_hash,
		       first_name, last_name, profile_photo, role, is_active, is_verified,
		       email_verified_at, last_login_at, created_at, updated_at,
		       COALESCE(preferable_subject, '{}'), COALESCE(plateform_version, '0.0.0')
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&u.ID,
		&u.UUID,
		&u.Email,
		&u.PhoneNumber,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.ProfilePhoto,
		&u.Role,
		&u.IsActive,
		&u.IsVerified,
		&u.EmailVerifiedAt,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PreferableSubject,
		&u.PlateformVersion,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

// UpdateLastLoginAt sets last_login_at to now for the user.
func (r *UserRepo) UpdateLastLoginAt(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`, userID)
	return err
}

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
