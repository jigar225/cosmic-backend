package repository

import (
	"context"
	"errors"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrGradeMethodNotFound = errors.New("grade method not found")
	ErrGradeMethodConflict = errors.New("grade method with this title already exists")
	ErrGradeMethodInUse   = errors.New("grade method cannot be deleted: it is used by boards or grades")
)

// GradeMethodRepo handles persistence for grade_methods.
type GradeMethodRepo struct {
	pool *pgxpool.Pool
}

// NewGradeMethodRepo returns a new GradeMethodRepo.
func NewGradeMethodRepo(pool *pgxpool.Pool) *GradeMethodRepo {
	return &GradeMethodRepo{pool: pool}
}

// ListVisible returns only visible grade methods (for public/dropdowns).
func (r *GradeMethodRepo) ListVisible(ctx context.Context) ([]domain.GradeMethod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, is_visible, created_at
		FROM grade_methods
		WHERE is_visible = true
		ORDER BY title
	`)
	if err != nil {
		return nil, err
	}
	return scanGradeMethods(rows)
}

// ListAll returns all grade methods (admin).
func (r *GradeMethodRepo) ListAll(ctx context.Context) ([]domain.GradeMethod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, is_visible, created_at
		FROM grade_methods
		ORDER BY title
	`)
	if err != nil {
		return nil, err
	}
	return scanGradeMethods(rows)
}

func scanGradeMethods(rows pgx.Rows) ([]domain.GradeMethod, error) {
	defer rows.Close()
	var out []domain.GradeMethod
	for rows.Next() {
		var g domain.GradeMethod
		if err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.IsVisible, &g.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// Create inserts a grade method. Returns ErrGradeMethodConflict if title already exists.
func (r *GradeMethodRepo) Create(ctx context.Context, in domain.CreateGradeMethodInput) (domain.GradeMethod, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM grade_methods WHERE title = $1`, in.Title).Scan(&dummy)
	if err == nil {
		return domain.GradeMethod{}, ErrGradeMethodConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.GradeMethod{}, err
	}

	var g domain.GradeMethod
	err = r.pool.QueryRow(ctx, `
		INSERT INTO grade_methods (title, description, is_visible)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, is_visible, created_at
	`, in.Title, in.Description, in.IsVisible).Scan(
		&g.ID, &g.Title, &g.Description, &g.IsVisible, &g.CreatedAt,
	)
	if err != nil {
		return domain.GradeMethod{}, err
	}
	return g, nil
}

// Delete removes a grade method. Returns ErrGradeMethodNotFound if id does not exist.
// Returns ErrGradeMethodInUse if any boards or grades reference this grade method.
func (r *GradeMethodRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM boards WHERE grade_method_id = $1) +
			(SELECT COUNT(*) FROM grades WHERE grade_method_id = $1)
		) AS in_use
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrGradeMethodInUse
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM grade_methods WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrGradeMethodNotFound
	}
	return nil
}

// GetByID returns a grade method by id. Returns ErrGradeMethodNotFound if not found.
func (r *GradeMethodRepo) GetByID(ctx context.Context, id int64) (domain.GradeMethod, error) {
	var g domain.GradeMethod
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, is_visible, created_at
		FROM grade_methods WHERE id = $1
	`, id).Scan(&g.ID, &g.Title, &g.Description, &g.IsVisible, &g.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.GradeMethod{}, ErrGradeMethodNotFound
		}
		return domain.GradeMethod{}, err
	}
	return g, nil
}

// Update applies partial updates. Returns ErrGradeMethodNotFound if id does not exist.
func (r *GradeMethodRepo) Update(ctx context.Context, id int64, in domain.UpdateGradeMethodInput) (domain.GradeMethod, error) {
	var g domain.GradeMethod
	err := r.pool.QueryRow(ctx, `
		UPDATE grade_methods
		SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			is_visible = COALESCE($4, is_visible)
		WHERE id = $1
		RETURNING id, title, description, is_visible, created_at
	`, id, in.Title, in.Description, in.IsVisible).Scan(
		&g.ID, &g.Title, &g.Description, &g.IsVisible, &g.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.GradeMethod{}, ErrGradeMethodNotFound
		}
		return domain.GradeMethod{}, err
	}
	return g, nil
}
