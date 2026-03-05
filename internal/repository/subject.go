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
	ErrSubjectHasDependents = errors.New("subject cannot be deleted: it is used by books, user context, or generated content") // unused when cascade delete is used
)

// SubjectRepo handles persistence for subjects.
type SubjectRepo struct {
	pool *pgxpool.Pool
}

// NewSubjectRepo returns a new SubjectRepo.
func NewSubjectRepo(pool *pgxpool.Pool) *SubjectRepo {
	return &SubjectRepo{pool: pool}
}

// ListVisible returns visible subjects filtered by board, medium, and grade (for public selection).
func (r *SubjectRepo) ListVisible(ctx context.Context, boardID int64, mediumID *int64, gradeID *int64) ([]domain.Subject, error) {
	return r.list(ctx, &boardID, mediumID, gradeID, true)
}

// ListAll returns all subjects (admin). Optional board, medium, grade filters.
func (r *SubjectRepo) ListAll(ctx context.Context, boardID *int64, mediumID *int64, gradeID *int64) ([]domain.Subject, error) {
	return r.list(ctx, boardID, mediumID, gradeID, false)
}

func (r *SubjectRepo) list(ctx context.Context, boardID *int64, mediumID *int64, gradeID *int64, visibleOnly bool) ([]domain.Subject, error) {
	query := `
		SELECT id, board_id, medium_id, grade_id, title, is_visible, created_at, updated_at
		FROM subjects
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1
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
	query += " ORDER BY title"

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
			&s.BoardID,
			&s.MediumID,
			&s.GradeID,
			&s.Title,
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

// Create inserts a subject. Returns ErrSubjectConflict if duplicate (same board/medium/grade/title).
func (r *SubjectRepo) Create(ctx context.Context, in domain.CreateSubjectInput) (domain.Subject, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM subjects
		WHERE board_id = $1 AND medium_id = $2 AND grade_id = $3 AND title = $4
	`, in.BoardID, in.MediumID, in.GradeID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Subject{}, ErrSubjectConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Subject{}, err
	}

	var s domain.Subject
	err = r.pool.QueryRow(ctx, `
		INSERT INTO subjects (board_id, medium_id, grade_id, title, is_visible)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, board_id, medium_id, grade_id, title, is_visible, created_at, updated_at
	`,
		in.BoardID,
		in.MediumID,
		in.GradeID,
		in.Title,
		in.IsVisible,
	).Scan(
		&s.ID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
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
		SELECT id, board_id, medium_id, grade_id, title, is_visible, created_at, updated_at
		FROM subjects WHERE id = $1
	`, id).Scan(
		&s.ID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
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
			is_visible = COALESCE($3, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, board_id, medium_id, grade_id, title, is_visible, created_at, updated_at
	`,
		id,
		in.Title,
		in.IsVisible,
	).Scan(
		&s.ID,
		&s.BoardID,
		&s.MediumID,
		&s.GradeID,
		&s.Title,
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
// Delete removes a subject and all its dependents (cascade): generated_content for this subject,
// chapters of books for this subject, then books for this subject, then the subject.
// So a wrong/empty subject (including its default unit book) can always be deleted in one step.
func (r *SubjectRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Verify subject exists
	var dummy int64
	err = tx.QueryRow(ctx, `SELECT id FROM subjects WHERE id = $1`, id).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSubjectNotFound
		}
		return err
	}

	// Cascade delete order: generated_content (by subject_id) → chapters (by book_id) → books → subject
	_, err = tx.Exec(ctx, `DELETE FROM generated_content WHERE subject_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM chapters WHERE book_id IN (SELECT id FROM books WHERE subject_id = $1)`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM books WHERE subject_id = $1`, id)
	if err != nil {
		return err
	}
	cmd, err := tx.Exec(ctx, `DELETE FROM subjects WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrSubjectNotFound
	}
	return tx.Commit(ctx)
}

