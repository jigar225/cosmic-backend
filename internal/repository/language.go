package repository

import (
	"context"
	"errors"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Domain-level errors for language operations.
var (
	ErrLanguageNotFound      = errors.New("language not found")
	ErrLanguageConflict      = errors.New("language already exists")
	ErrLanguageHasDependents = errors.New("language cannot be deleted: it is used by mediums")
)

// LanguageRepo handles persistence for languages.
type LanguageRepo struct {
	pool *pgxpool.Pool
}

// NewLanguageRepo returns a new LanguageRepo.
func NewLanguageRepo(pool *pgxpool.Pool) *LanguageRepo {
	return &LanguageRepo{pool: pool}
}

// ListAll returns all languages, ordered by name.
func (r *LanguageRepo) ListAll(ctx context.Context) ([]domain.Language, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, code, name, is_visible, created_at, updated_at
		FROM languages
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Language
	for rows.Next() {
		var l domain.Language
		if err := rows.Scan(
			&l.ID,
			&l.Code,
			&l.Name,
			&l.IsVisible,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// Create inserts a new language row and returns it.
func (r *LanguageRepo) Create(ctx context.Context, in domain.CreateLanguageInput) (domain.Language, error) {
	// Prevent duplicates by code.
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM languages WHERE code = $1
	`, in.Code).Scan(&dummy)
	if err == nil {
		return domain.Language{}, ErrLanguageConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Language{}, err
	}

	var l domain.Language
	err = r.pool.QueryRow(ctx, `
		INSERT INTO languages (code, name, is_visible)
		VALUES ($1, $2, $3)
		RETURNING id, code, name, is_visible, created_at, updated_at
	`,
		in.Code,
		in.Name,
		in.IsVisible,
	).Scan(
		&l.ID,
		&l.Code,
		&l.Name,
		&l.IsVisible,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return domain.Language{}, err
	}
	return l, nil
}

// Update applies partial updates to a language row and returns the updated row.
func (r *LanguageRepo) Update(ctx context.Context, id int64, in domain.UpdateLanguageInput) (domain.Language, error) {
	var l domain.Language
	err := r.pool.QueryRow(ctx, `
		UPDATE languages
		SET
			name = COALESCE($2, name),
			is_visible = COALESCE($3, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, code, name, is_visible, created_at, updated_at
	`,
		id,
		in.Name,
		in.IsVisible,
	).Scan(
		&l.ID,
		&l.Code,
		&l.Name,
		&l.IsVisible,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Language{}, ErrLanguageNotFound
		}
		return domain.Language{}, err
	}
	return l, nil
}

// Delete removes a language. Returns ErrLanguageNotFound if id does not exist.
// Returns ErrLanguageHasDependents if any mediums reference this language.
func (r *LanguageRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM mediums WHERE language_id = $1
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrLanguageHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM languages WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrLanguageNotFound
	}
	return nil
}


