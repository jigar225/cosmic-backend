package repository

import (
	"context"
	"errors"
	"time"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// RefreshTokenRepo handles persistence for refresh tokens.
type RefreshTokenRepo struct {
	pool *pgxpool.Pool
}

// NewRefreshTokenRepo returns a new RefreshTokenRepo.
func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool: pool}
}

// Create inserts a new refresh token row and returns it.
func (r *RefreshTokenRepo) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time, deviceName, deviceType, ipAddress *string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.pool.QueryRow(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, device_name, device_type, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, token_hash, device_name, device_type, ip_address, expires_at, created_at, revoked
	`,
		userID, tokenHash, deviceName, deviceType, ipAddress, expiresAt,
	).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.DeviceName,
		&rt.DeviceType,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.Revoked,
	)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return rt, nil
}

// GetByTokenHash returns the token row if found, not revoked, and not expired.
func (r *RefreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, device_name, device_type, ip_address, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()
	`, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.DeviceName,
		&rt.DeviceType,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.Revoked,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshToken{}, ErrRefreshTokenNotFound
		}
		return domain.RefreshToken{}, err
	}
	return rt, nil
}

// RevokeByTokenHash marks the token as revoked (or deletes it; we delete for simplicity so rotation is clean).
func (r *RefreshTokenRepo) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	res, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	}
	return nil
}

// RevokeAllForUser deletes all refresh tokens for the user (logout everywhere).
func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

// DeleteByTokenHash removes the token row (used after rotation: delete old, insert new).
func (r *RefreshTokenRepo) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}
