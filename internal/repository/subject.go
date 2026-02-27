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
	ErrSubjectNotFound      = errors.New("subject not found")
	ErrSubjectConflict      = errors.New("subject already exists for this country/board/medium/grade with this title")
	ErrSubjectHasDependents = errors.New("subject cannot be deleted: it is used by books, user context, or generated content")
)

// SubjectRepo handles persistence for subjects.
type SubjectRepo struct {
	pool *pgxpool.Pool
}

// NewSubjectRepo returns a new SubjectRepo.
func NewSubjectRepo(pool *pgxpool.Pool) *SubjectRepo {
	return &SubjectRepo{pool: pool}
}

// ListVisible returns visible subjects filtered by country, board, medium, and grade (for public selection).
func (r *SubjectRepo) ListVisible(ctx context.Context, countryID int64, boardID int64, mediumID *int64, gradeID *int64) ([]domain.Subject, error) {
	return r.list(ctx, &countryID, &boardID, mediumID, gradeID, true)
}

// ListAll returns all subjects (admin). Optional country, board, medium, grade filters.
func (r *SubjectRepo) ListAll(ctx context.Context, countryID *int64, boardID *int64, mediumID *int64, gradeID *int64) ([]domain.Subject, error) {
	return r.list(ctx, countryID, boardID, mediumID, gradeID, false)
}

func (r *SubjectRepo) list(ctx context.Context, countryID *int64, boardID *int64, mediumID *int64, gradeID *int64, visibleOnly bool) ([]domain.Subject, error) {
	query := `
		SELECT id, country_id, board_id, medium_id, grade_id,
		       title, subject_type, is_visible, created_at, updated_at
		FROM subjects
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
	if mediumID != nil {
		query += " AND medium_id = $" + strconv.Itoa(n)
		args = append(args, *mediumID)
		n++
	}
	if gradeID != nil {
		query += " AND grade_id = $" + strconv.Itoa(n)
		args = append(args, *gradeID)
		n++
	}
	if visibleOnly {
		query += " AND is_visible = true"
	}
	query += " ORDER BY sequence_order NULLS LAST, title"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(
			&s.ID,
			&s.CountryID,
			&s.BoardID,
			&s.MediumID,
			&s.GradeID,
			&s.Title,
			&s.SubjectType,
			&s.IsVisible,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Create inserts a subject. Returns ErrSubjectConflict if duplicate.
func (r *SubjectRepo) Create(ctx context.Context, in domain.CreateSubjectInput) (domain.Subject, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM subjects
		WHERE country_id = $1
		  AND board_id = $2
		  AND medium_id = $3
		  AND grade_id = $4
		  AND title = $5
	`, in.CountryID, in.BoardID, in.MediumID, in.GradeID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Subject{}, ErrSubjectConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Subject{}, err
	}

	if in.SubjectType == "" {
		in.SubjectType = "core"
	}

	var s domain.Subject
	err = r.pool.QueryRow(ctx, `
		INSERT INTO subjects (
			country_id, board_id, medium_id, grade_id,
			title, subject_type, is_visible
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, country_id, board_id, medium_id, grade_id,
		          title, subject_type, is_visible, created_at, updated_at
	`,
		in.CountryID,
		in.BoardID,
		in.MediumID,
		in.GradeID,
		in.Title,
		in.SubjectType,
		in.IsVisible,
	).Scan(
		&s.ID,
		&s.CountryID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
		&s.SubjectType,
		&s.IsVisible,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return domain.Subject{}, err
	}
	return s, nil
}

// GetByID returns a subject by id.
func (r *SubjectRepo) GetByID(ctx context.Context, id int64) (domain.Subject, error) {
	var s domain.Subject
	err := r.pool.QueryRow(ctx, `
		SELECT id, country_id, board_id, medium_id, grade_id,
		       title, subject_type, is_visible, created_at, updated_at
		FROM subjects WHERE id = $1
	`, id).Scan(
		&s.ID,
		&s.CountryID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
		&s.SubjectType,
		&s.IsVisible,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subject{}, ErrSubjectNotFound
		}
		return domain.Subject{}, err
	}
	return s, nil
}

// Update applies partial updates to a subject.
func (r *SubjectRepo) Update(ctx context.Context, id int64, in domain.UpdateSubjectInput) (domain.Subject, error) {
	var s domain.Subject
	err := r.pool.QueryRow(ctx, `
		UPDATE subjects
		SET
			title = COALESCE($2, title),
			subject_type = COALESCE($3, subject_type),
			is_visible = COALESCE($4, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, country_id, board_id, medium_id, grade_id,
		          title, subject_type, is_visible, created_at, updated_at
	`,
		id,
		in.Title,
		in.SubjectType,
		in.IsVisible,
	).Scan(
		&s.ID,
		&s.CountryID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
		&s.SubjectType,
		&s.IsVisible,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subject{}, ErrSubjectNotFound
		}
		return domain.Subject{}, err
	}
	return s, nil
}

// Delete removes a subject. Returns ErrSubjectNotFound if id does not exist.
// Returns ErrSubjectHasDependents if any books, user_context, or generated_content reference this subject.
func (r *SubjectRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM books WHERE subject_id = $1) +
			(SELECT COUNT(*) FROM user_context WHERE current_subject_id = $1) +
			(SELECT COUNT(*) FROM generated_content WHERE subject_id = $1)
		) AS dependents
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSubjectHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM subjects WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrSubjectNotFound
	}
	return nil
}

