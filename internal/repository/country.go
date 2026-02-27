package repository

import (
	"context"
	"errors"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Domain-level errors for country operations.
var (
	ErrCountryNotFound = errors.New("country not found")
	ErrCountryConflict = errors.New("country already exists")
)

// CountryRepo handles persistence for countries.
type CountryRepo struct {
	pool *pgxpool.Pool
}

// NewCountryRepo returns a new CountryRepo.
func NewCountryRepo(pool *pgxpool.Pool) *CountryRepo {
	return &CountryRepo{pool: pool}
}

// ListVisible returns only countries marked as visible, ordered by title.
func (r *CountryRepo) ListVisible(ctx context.Context) ([]domain.Country, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, country_code, title, icon_path, phone_code, signup_method,
		       have_board, has_states, is_visible, created_at, updated_at
		FROM countries
		WHERE is_visible = TRUE
		ORDER BY title
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Country
	for rows.Next() {
		var c domain.Country
		if err := rows.Scan(
			&c.ID,
			&c.CountryCode,
			&c.Title,
			&c.IconPath,
			&c.PhoneCode,
			&c.SignupMethod,
			&c.HaveBoard,
			&c.HasStates,
			&c.IsVisible,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListAll returns all countries (visible and hidden), ordered by title.
func (r *CountryRepo) ListAll(ctx context.Context) ([]domain.Country, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, country_code, title, icon_path, phone_code, signup_method,
		       have_board, has_states, is_visible, created_at, updated_at
		FROM countries
		ORDER BY title
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Country
	for rows.Next() {
		var c domain.Country
		if err := rows.Scan(
			&c.ID,
			&c.CountryCode,
			&c.Title,
			&c.IconPath,
			&c.PhoneCode,
			&c.SignupMethod,
			&c.HaveBoard,
			&c.HasStates,
			&c.IsVisible,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Create inserts a new country row and returns it.
func (r *CountryRepo) Create(ctx context.Context, in domain.CreateCountryInput) (domain.Country, error) {
	// Prevent duplicates by country_code.
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM countries WHERE country_code = $1
	`, in.CountryCode).Scan(&dummy)
	if err == nil {
		return domain.Country{}, ErrCountryConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Country{}, err
	}

	var c domain.Country
	err = r.pool.QueryRow(ctx, `
		INSERT INTO countries (
			country_code, title, icon_path, phone_code, signup_method,
			have_board, has_states, is_visible
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, country_code, title, icon_path, phone_code, signup_method,
		          have_board, has_states, is_visible, created_at, updated_at
	`,
		in.CountryCode,
		in.Title,
		in.IconPath,
		in.PhoneCode,
		in.SignupMethod,
		in.HaveBoard,
		in.HasStates,
		in.IsVisible,
	).Scan(
		&c.ID,
		&c.CountryCode,
		&c.Title,
		&c.IconPath,
		&c.PhoneCode,
		&c.SignupMethod,
		&c.HaveBoard,
		&c.HasStates,
		&c.IsVisible,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return domain.Country{}, err
	}
	return c, nil
}

// Update applies partial updates to a country row and returns the updated row.
func (r *CountryRepo) Update(ctx context.Context, id int64, in domain.UpdateCountryInput) (domain.Country, error) {
	// Build a dynamic update using COALESCE-like semantics by passing through existing values.
	var c domain.Country
	err := r.pool.QueryRow(ctx, `
		UPDATE countries
		SET
			title = COALESCE($2, title),
			icon_path = COALESCE($3, icon_path),
			phone_code = COALESCE($4, phone_code),
			signup_method = COALESCE($5, signup_method),
			have_board = COALESCE($6, have_board),
			has_states = COALESCE($7, has_states),
			is_visible = COALESCE($8, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, country_code, title, icon_path, phone_code, signup_method,
		          have_board, has_states, is_visible, created_at, updated_at
	`,
		id,
		in.Title,
		in.IconPath,
		in.PhoneCode,
		in.SignupMethod,
		in.HaveBoard,
		in.HasStates,
		in.IsVisible,
	).Scan(
		&c.ID,
		&c.CountryCode,
		&c.Title,
		&c.IconPath,
		&c.PhoneCode,
		&c.SignupMethod,
		&c.HaveBoard,
		&c.HasStates,
		&c.IsVisible,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Country{}, ErrCountryNotFound
		}
		return domain.Country{}, err
	}
	return c, nil
}

