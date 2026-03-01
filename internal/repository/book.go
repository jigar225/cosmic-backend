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
	ErrBookConflict      = errors.New("book already exists for this subject/medium/grade with this title")
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

// ListVisible returns visible books filtered by subject and optional medium/grade (for public selection).
func (r *BookRepo) ListVisible(ctx context.Context, subjectID int64, mediumID *int64, gradeID *int64) ([]domain.Book, error) {
	return r.list(ctx, &subjectID, mediumID, gradeID, true)
}

// ListAll returns all books (admin). Optional subject, medium, grade filters.
func (r *BookRepo) ListAll(ctx context.Context, subjectID *int64, mediumID *int64, gradeID *int64) ([]domain.Book, error) {
	return r.list(ctx, subjectID, mediumID, gradeID, false)
}

func (r *BookRepo) list(ctx context.Context, subjectID *int64, mediumID *int64, gradeID *int64, visibleOnly bool) ([]domain.Book, error) {
	query := `
		SELECT id, book_type, subject_id, medium_id, grade_id,
		       title, author, publisher, edition, publication_year,
		       isbn, book_code, uploaded_by_user_id, is_public,
		       curriculum_version, status, original_file_path,
		       processed_file_path, cover_image_url, is_visible,
		       view_count, download_count, created_at, updated_at
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

	var out []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(
			&b.ID,
			&b.BookType,
			&b.SubjectID,
			&b.MediumID,
			&b.GradeID,
			&b.Title,
			&b.Author,
			&b.Publisher,
			&b.Edition,
			&b.PublicationYear,
			&b.ISBN,
			&b.BookCode,
			&b.UploadedByUserID,
			&b.IsPublic,
			&b.CurriculumVersion,
			&b.Status,
			&b.OriginalFilePath,
			&b.ProcessedFilePath,
			&b.CoverImageURL,
			&b.IsVisible,
			&b.ViewCount,
			&b.DownloadCount,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// Create inserts a book. Returns ErrBookConflict if duplicate.
func (r *BookRepo) Create(ctx context.Context, in domain.CreateBookInput) (domain.Book, error) {
	var dummy int
	err := r.pool.QueryRow(ctx, `
		SELECT 1 FROM books
		WHERE subject_id = $1
		  AND medium_id = $2
		  AND grade_id = $3
		  AND title = $4
	`, in.SubjectID, in.MediumID, in.GradeID, in.Title).Scan(&dummy)
	if err == nil {
		return domain.Book{}, ErrBookConflict
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.Book{}, err
	}

	if in.Status == "" {
		in.Status = "draft"
	}

	var b domain.Book
	err = r.pool.QueryRow(ctx, `
		INSERT INTO books (
			book_type, subject_id, medium_id, grade_id,
			title, author, publisher, edition, publication_year,
			isbn, book_code, uploaded_by_user_id, is_public,
			curriculum_version, status, original_file_path,
			processed_file_path, cover_image_url, is_visible
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,
		        $10, $11, $12, $13, $14, $15, $16,
		        $17, $18, $19)
		RETURNING id, book_type, subject_id, medium_id, grade_id,
		          title, author, publisher, edition, publication_year,
		          isbn, book_code, uploaded_by_user_id, is_public,
		          curriculum_version, status, original_file_path,
		          processed_file_path, cover_image_url, is_visible,
		          view_count, download_count, created_at, updated_at
	`,
		in.BookType,
		in.SubjectID,
		in.MediumID,
		in.GradeID,
		in.Title,
		in.Author,
		in.Publisher,
		in.Edition,
		in.PublicationYear,
		in.ISBN,
		in.BookCode,
		in.UploadedByUserID,
		in.IsPublic,
		in.CurriculumVersion,
		in.Status,
		in.OriginalFilePath,
		in.ProcessedFilePath,
		in.CoverImageURL,
		in.IsVisible,
	).Scan(
		&b.ID,
		&b.BookType,
		&b.SubjectID,
		&b.MediumID,
		&b.GradeID,
		&b.Title,
		&b.Author,
		&b.Publisher,
		&b.Edition,
		&b.PublicationYear,
		&b.ISBN,
		&b.BookCode,
		&b.UploadedByUserID,
		&b.IsPublic,
		&b.CurriculumVersion,
		&b.Status,
		&b.OriginalFilePath,
		&b.ProcessedFilePath,
		&b.CoverImageURL,
		&b.IsVisible,
		&b.ViewCount,
		&b.DownloadCount,
		&b.CreatedAt,
		&b.UpdatedAt,
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
		SELECT id, book_type, subject_id, medium_id, grade_id,
		       title, author, publisher, edition, publication_year,
		       isbn, book_code, uploaded_by_user_id, is_public,
		       curriculum_version, status, original_file_path,
		       processed_file_path, cover_image_url, is_visible,
		       view_count, download_count, created_at, updated_at
		FROM books WHERE id = $1
	`, id).Scan(
		&b.ID,
		&b.BookType,
		&b.SubjectID,
		&b.MediumID,
		&b.GradeID,
		&b.Title,
		&b.Author,
		&b.Publisher,
		&b.Edition,
		&b.PublicationYear,
		&b.ISBN,
		&b.BookCode,
		&b.UploadedByUserID,
		&b.IsPublic,
		&b.CurriculumVersion,
		&b.Status,
		&b.OriginalFilePath,
		&b.ProcessedFilePath,
		&b.CoverImageURL,
		&b.IsVisible,
		&b.ViewCount,
		&b.DownloadCount,
		&b.CreatedAt,
		&b.UpdatedAt,
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
			book_type = COALESCE($2, book_type),
			title = COALESCE($3, title),
			author = COALESCE($4, author),
			publisher = COALESCE($5, publisher),
			edition = COALESCE($6, edition),
			publication_year = COALESCE($7, publication_year),
			isbn = COALESCE($8, isbn),
			book_code = COALESCE($9, book_code),
			is_public = COALESCE($10, is_public),
			curriculum_version = COALESCE($11, curriculum_version),
			status = COALESCE($12, status),
			processed_file_path = COALESCE($13, processed_file_path),
			cover_image_url = COALESCE($14, cover_image_url),
			is_visible = COALESCE($15, is_visible),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, book_type, subject_id, medium_id, grade_id,
		          title, author, publisher, edition, publication_year,
		          isbn, book_code, uploaded_by_user_id, is_public,
		          curriculum_version, status, original_file_path,
		          processed_file_path, cover_image_url, is_visible,
		          view_count, download_count, created_at, updated_at
	`,
		id,
		in.BookType,
		in.Title,
		in.Author,
		in.Publisher,
		in.Edition,
		in.PublicationYear,
		in.ISBN,
		in.BookCode,
		in.IsPublic,
		in.CurriculumVersion,
		in.Status,
		in.ProcessedFilePath,
		in.CoverImageURL,
		in.IsVisible,
	).Scan(
		&b.ID,
		&b.BookType,
		&b.SubjectID,
		&b.MediumID,
		&b.GradeID,
		&b.Title,
		&b.Author,
		&b.Publisher,
		&b.Edition,
		&b.PublicationYear,
		&b.ISBN,
		&b.BookCode,
		&b.UploadedByUserID,
		&b.IsPublic,
		&b.CurriculumVersion,
		&b.Status,
		&b.OriginalFilePath,
		&b.ProcessedFilePath,
		&b.CoverImageURL,
		&b.IsVisible,
		&b.ViewCount,
		&b.DownloadCount,
		&b.CreatedAt,
		&b.UpdatedAt,
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
// Returns ErrBookHasDependents if any chapters or generated_content reference this book via chapters.
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

