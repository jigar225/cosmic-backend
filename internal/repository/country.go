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
		SELECT id, country_code, title, phone_code, signup_methods,
		       have_board, is_visible, created_at
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
			&c.PhoneCode,
			&c.SignupMethods,
			&c.HaveBoard,
			&c.IsVisible,
			&c.CreatedAt,
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
		SELECT id, country_code, title, phone_code, signup_methods,
		       have_board, is_visible, created_at
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
			&c.PhoneCode,
			&c.SignupMethods,
			&c.HaveBoard,
			&c.IsVisible,
			&c.CreatedAt,
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

	signupMethods := in.SignupMethods
	if signupMethods == nil {
		signupMethods = []string{"email"}
	}
	var c domain.Country
	err = r.pool.QueryRow(ctx, `
		INSERT INTO countries (
			country_code, title, phone_code, signup_methods,
			have_board, is_visible
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, country_code, title, phone_code, signup_methods,
		          have_board, is_visible, created_at
	`,
		in.CountryCode,
		in.Title,
		in.PhoneCode,
		signupMethods,
		in.HaveBoard,
		in.IsVisible,
	).Scan(
		&c.ID,
		&c.CountryCode,
		&c.Title,
		&c.PhoneCode,
		&c.SignupMethods,
		&c.HaveBoard,
		&c.IsVisible,
		&c.CreatedAt,
	)
	if err != nil {
		return domain.Country{}, err
	}
	return c, nil
}

// Update applies partial updates to a country row and returns the updated row.
func (r *CountryRepo) Update(ctx context.Context, id int64, in domain.UpdateCountryInput) (domain.Country, error) {
	// When SignupMethods is nil, pass nil so COALESCE keeps existing value.
	var signupMethodsArg interface{} = nil
	if in.SignupMethods != nil {
		signupMethodsArg = *in.SignupMethods
	}
	var c domain.Country
	err := r.pool.QueryRow(ctx, `
		UPDATE countries
		SET
			title = COALESCE($2, title),
			phone_code = COALESCE($3, phone_code),
			signup_methods = COALESCE($4::text[], signup_methods),
			have_board = COALESCE($5, have_board),
			is_visible = COALESCE($6, is_visible)
		WHERE id = $1
		RETURNING id, country_code, title, phone_code, signup_methods,
		          have_board, is_visible, created_at
	`,
		id,
		in.Title,
		in.PhoneCode,
		signupMethodsArg,
		in.HaveBoard,
		in.IsVisible,
	).Scan(
		&c.ID,
		&c.CountryCode,
		&c.Title,
		&c.PhoneCode,
		&c.SignupMethods,
		&c.HaveBoard,
		&c.IsVisible,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Country{}, ErrCountryNotFound
		}
		return domain.Country{}, err
	}
	return c, nil
}

