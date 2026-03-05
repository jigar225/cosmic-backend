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
	ErrBookNotFound      = errors.New("book not found")
	ErrBookConflict      = errors.New("book already exists for this subject with this title")
	ErrBookHasDependents = errors.New("book cannot be deleted: it is used by chapters or generated content")
)

// BookRepo handles persistence for books.
type BookRepo struct {
	pool *pgxpool.Pool
}

// NewBookRepo returns a new BookRepo.
func NewBookRepo(pool *pgxpool.Pool) *BookRepo {
	return &BookRepo{pool: pool}
}

// CreateDefaultForSubject inserts a single "unit" book for the given subject (created_by = NULL for admin/system).
// file_path is NULL — no whole-book file; chapters are uploaded separately.
func (r *BookRepo) CreateDefaultForSubject(ctx context.Context, subjectID int64, subjectTitle string) (int64, error) {
	title := "Unit book - " + subjectTitle
	if len(title) > 500 {
		title = title[:497] + "..."
	}
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO books (subject_id, title, is_visible, is_active, created_by)
		VALUES ($1, $2, true, true, NULL)
		RETURNING id
	`, subjectID, title).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ListVisible returns visible books filtered by subject (for public selection).
func (r *BookRepo) ListVisible(ctx context.Context, subjectID int64) ([]domain.Book, error) {
	return r.list(ctx, &subjectID, true)
}

// ListAll returns all books (admin). Optional subject_id filter.
func (r *BookRepo) ListAll(ctx context.Context, subjectID *int64) ([]domain.Book, error) {
	return r.list(ctx, subjectID, false)
}

func (r *BookRepo) list(ctx context.Context, subjectID *int64, visibleOnly bool) ([]domain.Book, error) {
	query := `
		SELECT id, subject_id, title, is_visible, is_active, file_path, created_by, created_at
		FROM books
		WHERE 1=1
	`
	args := []interface{}{}
	n := 1
	if subjectID != nil {
		query += " AND subject_id = $" + strconv.Itoa(n)
		args = append(args, *subjectID)
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

	var out []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(
			&b.ID,
			&b.SubjectID,
			&b.Title,
			&b.IsVisible,
			&b.IsActive,
			&b.FilePath,
			&b.CreatedBy,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// Create inserts a book. Returns ErrBookConflict if duplicate title for same subject.
func (r *BookRepo) Create(ctx context.Context, in domain.CreateBookInput) (domain.Book, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM books WHERE subject_id = $1 AND title = $2
	`, in.SubjectID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Book{}, ErrBookConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Book{}, err
	}

	var b domain.Book
	err = r.pool.QueryRow(ctx, `
		INSERT INTO books (subject_id, title, is_visible, is_active, file_path)
		VALUES ($1, $2, $3, true, $4)
		RETURNING id, subject_id, title, is_visible, is_active, file_path, created_by, created_at
	`, in.SubjectID, in.Title, in.IsVisible, in.FilePath).Scan(
		&b.ID,
		&b.SubjectID,
		&b.Title,
		&b.IsVisible,
		&b.IsActive,
		&b.FilePath,
		&b.CreatedBy,
		&b.CreatedAt,
	)
	if err != nil {
		return domain.Book{}, err
	}
	return b, nil
}

// GetByID returns a book by id.
func (r *BookRepo) GetByID(ctx context.Context, id int64) (domain.Book, error) {
	var b domain.Book
	err := r.pool.QueryRow(ctx, `
		SELECT id, subject_id, title, is_visible, is_active, file_path, created_by, created_at
		FROM books WHERE id = $1
	`, id).Scan(
		&b.ID,
		&b.SubjectID,
		&b.Title,
		&b.IsVisible,
		&b.IsActive,
		&b.FilePath,
		&b.CreatedBy,
		&b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return b, nil
}

// Update applies partial updates to a book.
func (r *BookRepo) Update(ctx context.Context, id int64, in domain.UpdateBookInput) (domain.Book, error) {
	var b domain.Book
	err := r.pool.QueryRow(ctx, `
		UPDATE books
		SET
			title = COALESCE($2, title),
			is_visible = COALESCE($3, is_visible),
			file_path = COALESCE($4, file_path)
		WHERE id = $1
		RETURNING id, subject_id, title, is_visible, is_active, file_path, created_by, created_at
	`, id, in.Title, in.IsVisible, in.FilePath).Scan(
		&b.ID,
		&b.SubjectID,
		&b.Title,
		&b.IsVisible,
		&b.IsActive,
		&b.FilePath,
		&b.CreatedBy,
		&b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return b, nil
}

// Delete removes a book. Returns ErrBookNotFound if id does not exist.
// Returns ErrBookHasDependents if any chapters or generated_content reference this book.
func (r *BookRepo) Delete(ctx context.Context, id int64) error {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT (
			(SELECT COUNT(*) FROM chapters WHERE book_id = $1) +
			(SELECT COUNT(*) FROM generated_content gc
			 JOIN chapters ch ON gc.chapter_id = ch.id
			 WHERE ch.book_id = $1)
		) AS dependents
	`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrBookHasDependents
	}
	cmd, err := r.pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrBookNotFound
	}
	return nil
}
