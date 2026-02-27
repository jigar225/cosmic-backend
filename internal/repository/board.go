package repository

import (
	"context"
	"errors"
	"fmt"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBoardNotFound       = errors.New("board not found")
	ErrBoardConflict       = errors.New("board already exists for this country")
	ErrBoardHasDependents   = errors.New("board cannot be deleted: it has states, mediums, subjects, or other data linked to it")
)

// BoardRepo handles boards persistence.
type BoardRepo struct {
	pool *pgxpool.Pool
}

// NewBoardRepo returns a new BoardRepo.
func NewBoardRepo(pool *pgxpool.Pool) *BoardRepo {
	return &BoardRepo{pool: pool}
}

// List returns all boards that are visible (for public API).
func (r *BoardRepo) List(ctx context.Context) ([]domain.Board, error) {
	return r.list(ctx, nil, true)
}

// ListByCountryID returns all boards for a country (admin: visible + hidden).
func (r *BoardRepo) ListByCountryID(ctx context.Context, countryID int64) ([]domain.Board, error) {
	return r.list(ctx, &countryID, false)
}

// ListAll returns all boards; if countryID is non-nil, filter by that country.
func (r *BoardRepo) ListAll(ctx context.Context, countryID *int64) ([]domain.Board, error) {
	return r.list(ctx, countryID, false)
}

func (r *BoardRepo) list(ctx context.Context, countryID *int64, visibleOnly bool) ([]domain.Board, error) {
	query := `
		SELECT id, country_id, title, grade_method_id, is_visible, created_at, updated_at
		FROM boards
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1
	if countryID != nil {
		query += fmt.Sprintf(" AND country_id = $%d", n)
		args = append(args, *countryID)
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
	var out []domain.Board
	for rows.Next() {
		var b domain.Board
		if err := rows.Scan(&b.ID, &b.CountryID, &b.Title, &b.GradeMethodID, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// Create inserts a board and returns it. Returns ErrBoardConflict if same country_id+title exists.
func (r *BoardRepo) Create(ctx context.Context, in domain.CreateBoardInput) (domain.Board, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM boards WHERE country_id = $1 AND title = $2`, in.CountryID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Board{}, ErrBoardConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Board{}, err
	}

	var b domain.Board
	err = r.pool.QueryRow(ctx, `
		INSERT INTO boards (country_id, title, grade_method_id, is_visible)
		VALUES ($1, $2, $3, $4)
		RETURNING id, country_id, title, grade_method_id, is_visible, created_at, updated_at
	`, in.CountryID, in.Title, in.GradeMethodID, in.IsVisible).Scan(
		&b.ID, &b.CountryID, &b.Title, &b.GradeMethodID, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return domain.Board{}, err
	}
	return b, nil
}

// GetByID returns a board by id. Returns ErrBoardNotFound if not found.
func (r *BoardRepo) GetByID(ctx context.Context, id int64) (domain.Board, error) {
	var b domain.Board
	err := r.pool.QueryRow(ctx, `
		SELECT id, country_id, title, grade_method_id, is_visible, created_at, updated_at
		FROM boards WHERE id = $1
	`, id).Scan(&b.ID, &b.CountryID, &b.Title, &b.GradeMethodID, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Board{}, ErrBoardNotFound
		}
		return domain.Board{}, err
	}
	return b, nil
}

// Delete removes a board. Returns ErrBoardNotFound if id does not exist.
// Returns ErrBoardHasDependents if any states, mediums, subjects, user_context, or generated_content reference this board.
func (r *BoardRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM states WHERE default_board_id = $1) +
			(SELECT COUNT(*) FROM mediums WHERE board_id = $1) +
			(SELECT COUNT(*) FROM subjects WHERE board_id = $1) +
			(SELECT COUNT(*) FROM user_context WHERE current_board_id = $1) +
			(SELECT COUNT(*) FROM generated_content WHERE board_id = $1)
		) AS dependents
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrBoardHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrBoardNotFound
	}
	return nil
}

// Update applies partial updates to a board. Returns ErrBoardNotFound if id does not exist.
func (r *BoardRepo) Update(ctx context.Context, id int64, in domain.UpdateBoardInput) (domain.Board, error) {
	var b domain.Board
	err := r.pool.QueryRow(ctx, `
		UPDATE boards
		SET
			title = COALESCE($2, title),
			grade_method_id = COALESCE($3, grade_method_id),
			is_visible = COALESCE($4, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, country_id, title, grade_method_id, is_visible, created_at, updated_at
	`, id, in.Title, in.GradeMethodID, in.IsVisible).Scan(
		&b.ID, &b.CountryID, &b.Title, &b.GradeMethodID, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Board{}, ErrBoardNotFound
		}
		return domain.Board{}, err
	}
	return b, nil
}
