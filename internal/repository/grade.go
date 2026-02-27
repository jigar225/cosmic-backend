package repository

import (
	"context"
	"errors"
	"strconv"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrGradeNotFound     = errors.New("grade not found")
	ErrGradeConflict     = errors.New("grade already exists for this grade method")
	ErrGradeHasDependents = errors.New("grade cannot be deleted: it is used by subjects, books, user context, or generated content")
)

// GradeRepo handles persistence for grades.
type GradeRepo struct {
	pool *pgxpool.Pool
}

// NewGradeRepo returns a new GradeRepo.
func NewGradeRepo(pool *pgxpool.Pool) *GradeRepo {
	return &GradeRepo{pool: pool}
}

// ListByGradeMethod returns all grades for a grade method (visible only by default).
func (r *GradeRepo) ListByGradeMethod(ctx context.Context, gradeMethodID int64, visibleOnly bool) ([]domain.Grade, error) {
	query := `
		SELECT id, grade_method_id, title, display_order, numeric_equivalent,
		       age_range_start, age_range_end, academic_stage, is_visible, created_at
		FROM grades
		WHERE grade_method_id = $1
	`
	args := []interface{}{gradeMethodID}
	if visibleOnly {
		query += " AND is_visible = true"
	}
	query += " ORDER BY display_order, id"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Grade
	for rows.Next() {
		var g domain.Grade
		if err := rows.Scan(
			&g.ID,
			&g.GradeMethodID,
			&g.Title,
			&g.DisplayOrder,
			&g.NumericEquivalent,
			&g.AgeRangeStart,
			&g.AgeRangeEnd,
			&g.AcademicStage,
			&g.IsVisible,
			&g.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// ListAll returns all grades (admin). Optional filter by gradeMethodID when >0.
func (r *GradeRepo) ListAll(ctx context.Context, gradeMethodID *int64) ([]domain.Grade, error) {
	query := `
		SELECT id, grade_method_id, title, display_order, numeric_equivalent,
		       age_range_start, age_range_end, academic_stage, is_visible, created_at
		FROM grades
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1
	if gradeMethodID != nil {
		query += " AND grade_method_id = $" + strconv.Itoa(n)
		args = append(args, *gradeMethodID)
		n++
	}
	query += " ORDER BY grade_method_id, display_order, id"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Grade
	for rows.Next() {
		var g domain.Grade
		if err := rows.Scan(
			&g.ID,
			&g.GradeMethodID,
			&g.Title,
			&g.DisplayOrder,
			&g.NumericEquivalent,
			&g.AgeRangeStart,
			&g.AgeRangeEnd,
			&g.AcademicStage,
			&g.IsVisible,
			&g.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// Create inserts a grade. Returns ErrGradeConflict if same (grade_method_id, title) exists.
func (r *GradeRepo) Create(ctx context.Context, in domain.CreateGradeInput) (domain.Grade, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM grades WHERE grade_method_id = $1 AND title = $2
	`, in.GradeMethodID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Grade{}, ErrGradeConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Grade{}, err
	}

	var g domain.Grade
	err = r.pool.QueryRow(ctx, `
		INSERT INTO grades (
			grade_method_id, title, display_order, numeric_equivalent,
			age_range_start, age_range_end, academic_stage, is_visible
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, grade_method_id, title, display_order, numeric_equivalent,
		          age_range_start, age_range_end, academic_stage, is_visible, created_at
	`,
		in.GradeMethodID,
		in.Title,
		in.DisplayOrder,
		in.NumericEquivalent,
		in.AgeRangeStart,
		in.AgeRangeEnd,
		in.AcademicStage,
		in.IsVisible,
	).Scan(
		&g.ID,
		&g.GradeMethodID,
		&g.Title,
		&g.DisplayOrder,
		&g.NumericEquivalent,
		&g.AgeRangeStart,
		&g.AgeRangeEnd,
		&g.AcademicStage,
		&g.IsVisible,
		&g.CreatedAt,
	)
	if err != nil {
		return domain.Grade{}, err
	}
	return g, nil
}

// GetByID returns a grade by id.
func (r *GradeRepo) GetByID(ctx context.Context, id int64) (domain.Grade, error) {
	var g domain.Grade
	err := r.pool.QueryRow(ctx, `
		SELECT id, grade_method_id, title, display_order, numeric_equivalent,
		       age_range_start, age_range_end, academic_stage, is_visible, created_at
		FROM grades WHERE id = $1
	`, id).Scan(
		&g.ID,
		&g.GradeMethodID,
		&g.Title,
		&g.DisplayOrder,
		&g.NumericEquivalent,
		&g.AgeRangeStart,
		&g.AgeRangeEnd,
		&g.AcademicStage,
		&g.IsVisible,
		&g.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Grade{}, ErrGradeNotFound
		}
		return domain.Grade{}, err
	}
	return g, nil
}

// Update applies partial updates to a grade.
func (r *GradeRepo) Update(ctx context.Context, id int64, in domain.UpdateGradeInput) (domain.Grade, error) {
	var g domain.Grade
	err := r.pool.QueryRow(ctx, `
		UPDATE grades
		SET
			title = COALESCE($2, title),
			display_order = COALESCE($3, display_order),
			numeric_equivalent = COALESCE($4, numeric_equivalent),
			age_range_start = COALESCE($5, age_range_start),
			age_range_end = COALESCE($6, age_range_end),
			academic_stage = COALESCE($7, academic_stage),
			is_visible = COALESCE($8, is_visible)
		WHERE id = $1
		RETURNING id, grade_method_id, title, display_order, numeric_equivalent,
		          age_range_start, age_range_end, academic_stage, is_visible, created_at
	`,
		id,
		in.Title,
		in.DisplayOrder,
		in.NumericEquivalent,
		in.AgeRangeStart,
		in.AgeRangeEnd,
		in.AcademicStage,
		in.IsVisible,
	).Scan(
		&g.ID,
		&g.GradeMethodID,
		&g.Title,
		&g.DisplayOrder,
		&g.NumericEquivalent,
		&g.AgeRangeStart,
		&g.AgeRangeEnd,
		&g.AcademicStage,
		&g.IsVisible,
		&g.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Grade{}, ErrGradeNotFound
		}
		return domain.Grade{}, err
	}
	return g, nil
}

// Delete removes a grade. Returns ErrGradeNotFound if id does not exist.
// Returns ErrGradeHasDependents if any subjects, books, user_context, or generated_content reference this grade.
func (r *GradeRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM subjects WHERE grade_id = $1) +
			(SELECT COUNT(*) FROM books WHERE grade_id = $1) +
			(SELECT COUNT(*) FROM user_context WHERE current_grade_id = $1) +
			(SELECT COUNT(*) FROM generated_content WHERE grade_id = $1)
		) AS dependents
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrGradeHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM grades WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrGradeNotFound
	}
	return nil
}

