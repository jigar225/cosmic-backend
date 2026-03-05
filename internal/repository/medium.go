package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMediumNotFound       = errors.New("medium not found")
	ErrMediumConflict       = errors.New("medium already exists for this country/board with this title")
	ErrMediumHasDependents  = errors.New("medium cannot be deleted: it is used by subjects, books, user context, or generated content")
)

// MediumRepo handles persistence for mediums.
type MediumRepo struct {
	pool *pgxpool.Pool
}

// NewMediumRepo returns a new MediumRepo.
func NewMediumRepo(pool *pgxpool.Pool) *MediumRepo {
	return &MediumRepo{pool: pool}
}

// ListVisible returns visible mediums filtered by country and optional board (for public selection).
func (r *MediumRepo) ListVisible(ctx context.Context, countryID int64, boardID *int64) ([]domain.Medium, error) {
	return r.list(ctx, &countryID, boardID, true)
}

// ListAll returns all mediums (admin). Optional country and board filters.
func (r *MediumRepo) ListAll(ctx context.Context, countryID *int64, boardID *int64) ([]domain.Medium, error) {
	return r.list(ctx, countryID, boardID, false)
}

func (r *MediumRepo) list(ctx context.Context, countryID *int64, boardID *int64, visibleOnly bool) ([]domain.Medium, error) {
	query := `
		SELECT m.id,
		       m.country_id,
		       m.board_id,
		       m.title,
		       l.code AS language_code,
		       m.is_visible,
		       m.created_at,
		       m.updated_at
		FROM mediums m
		LEFT JOIN languages l ON m.language_id = l.id
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1
	if countryID != nil {
		query += " AND country_id = $" + strconv.Itoa(n)
		args = append(args, *countryID)
		n++
	}
	if boardID != nil {
		query += " AND board_id = $" + strconv.Itoa(n)
		args = append(args, *boardID)
		n++
	}
	if visibleOnly {
		query += " AND is_visible = true"
	}
	query += " ORDER BY title"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Medium
	for rows.Next() {
		var m domain.Medium
		if err := rows.Scan(
			&m.ID,
			&m.CountryID,
			&m.BoardID,
			&m.Title,
			&m.LanguageCode,
			&m.IsVisible,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// Create inserts a medium. Returns ErrMediumConflict if duplicate.
func (r *MediumRepo) Create(ctx context.Context, in domain.CreateMediumInput) (domain.Medium, error) {
	// Ensure there is a language row when a language_code is provided.
	var languageID sql.NullInt64
	if in.LanguageCode != nil {
		var id int64
		if err := r.pool.QueryRow(ctx, `
			INSERT INTO languages (code, name)
			VALUES ($1, $2)
			ON CONFLICT (code) DO UPDATE
				SET name = EXCLUDED.name,
				    updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, *in.LanguageCode, in.Title).Scan(&id); err != nil {
			return domain.Medium{}, err
		}
		languageID = sql.NullInt64{Int64: id, Valid: true}
	}

	var dummy int
	var err error
	if languageID.Valid {
		// Prevent duplicate medium for same (country, board, language).
		err = r.pool.QueryRow(ctx, `
			SELECT 1 FROM mediums
			WHERE country_id = $1
			  AND board_id IS NOT DISTINCT FROM $2
			  AND language_id IS NOT DISTINCT FROM $3
		`, in.CountryID, in.BoardID, languageID).Scan(&dummy)
	} else {
		// Fallback to previous behavior when no language_code is provided.
		err = r.pool.QueryRow(ctx, `
			SELECT 1 FROM mediums
			WHERE country_id = $1
			  AND board_id IS NOT DISTINCT FROM $2
			  AND title = $3
		`, in.CountryID, in.BoardID, in.Title).Scan(&dummy)
	}
	if err == nil {
		return domain.Medium{}, ErrMediumConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Medium{}, err
	}

	var id int64
	err = r.pool.QueryRow(ctx, `
		INSERT INTO mediums (country_id, board_id, title, language_id, is_visible)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		in.CountryID,
		in.BoardID,
		in.Title,
		languageID,
		in.IsVisible,
	).Scan(&id)
	if err != nil {
		return domain.Medium{}, err
	}
	return r.GetByID(ctx, id)
}

// GetByID returns a medium by id.
func (r *MediumRepo) GetByID(ctx context.Context, id int64) (domain.Medium, error) {
	var m domain.Medium
	err := r.pool.QueryRow(ctx, `
		SELECT m.id,
		       m.country_id,
		       m.board_id,
		       m.title,
		       l.code AS language_code,
		       m.is_visible,
		       m.created_at,
		       m.updated_at
		FROM mediums m
		LEFT JOIN languages l ON m.language_id = l.id
		WHERE m.id = $1
	`, id).Scan(
		&m.ID,
		&m.CountryID,
		&m.BoardID,
		&m.Title,
		&m.LanguageCode,
		&m.IsVisible,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Medium{}, ErrMediumNotFound
		}
		return domain.Medium{}, err
	}
	return m, nil
}

// Update applies partial updates to a medium.
func (r *MediumRepo) Update(ctx context.Context, id int64, in domain.UpdateMediumInput) (domain.Medium, error) {
	// If language_code is being updated, ensure we have/derive a language row.
	var languageID sql.NullInt64
	if in.LanguageCode != nil {
		// Fetch existing medium to get a reasonable default name.
		existing, err := r.GetByID(ctx, id)
		if err != nil {
			return domain.Medium{}, err
		}
		name := existing.Title
		if in.Title != nil && *in.Title != "" {
			name = *in.Title
		}

		var lid int64
		if err := r.pool.QueryRow(ctx, `
			INSERT INTO languages (code, name)
			VALUES ($1, $2)
			ON CONFLICT (code) DO UPDATE
				SET name = EXCLUDED.name,
				    updated_at = CURRENT_TIMESTAMP
			RETURNING id
		`, *in.LanguageCode, name).Scan(&lid); err != nil {
			return domain.Medium{}, err
		}
		languageID = sql.NullInt64{Int64: lid, Valid: true}
	}

	var err error
	if in.LanguageCode != nil {
		// Update including language linkage.
		cmd, execErr := r.pool.Exec(ctx, `
			UPDATE mediums
			SET
				title = COALESCE($2, title),
				language_id = $3,
				is_visible = COALESCE($4, is_visible),
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`,
			id,
			in.Title,
			languageID,
			in.IsVisible,
		)
		if execErr != nil {
			err = execErr
		} else if cmd.RowsAffected() == 0 {
			err = ErrMediumNotFound
		}
	} else {
		// No language change; keep existing language linkage.
		cmd, execErr := r.pool.Exec(ctx, `
			UPDATE mediums
			SET
				title = COALESCE($2, title),
				is_visible = COALESCE($3, is_visible),
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`,
			id,
			in.Title,
			in.IsVisible,
		)
		if execErr != nil {
			err = execErr
		} else if cmd.RowsAffected() == 0 {
			err = ErrMediumNotFound
		}
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Medium{}, ErrMediumNotFound
		}
		return domain.Medium{}, err
	}
	return r.GetByID(ctx, id)
}

// Delete removes a medium. Returns ErrMediumNotFound if id does not exist.
// Returns ErrMediumHasDependents if any subjects, books, user_context, or generated_content reference this medium.
func (r *MediumRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM subjects WHERE medium_id = $1) +
			(SELECT COUNT(*) FROM books WHERE medium_id = $1) +
			(SELECT COUNT(*) FROM user_default WHERE current_medium_id = $1) +
			(SELECT COUNT(*) FROM generated_content WHERE medium_id = $1)
		) AS dependents
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrMediumHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM mediums WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrMediumNotFound
	}
	return nil
}

