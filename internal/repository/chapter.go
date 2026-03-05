package repository

import (
	"context"
	"errors"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrChapterNotFound      = errors.New("chapter not found")
	ErrChapterHasDependents = errors.New("chapter cannot be deleted: it is used by generated content")
)

// ChapterRepo handles persistence for chapters.
type ChapterRepo struct {
	pool *pgxpool.Pool
}

// NewChapterRepo returns a new ChapterRepo.
func NewChapterRepo(pool *pgxpool.Pool) *ChapterRepo {
	return &ChapterRepo{pool: pool}
}

// Create inserts a chapter.
func (r *ChapterRepo) Create(ctx context.Context, in domain.CreateChapterInput) (domain.Chapter, error) {
	var ch domain.Chapter

	// Resolve display_order: use provided value or auto-assign (max + 1 for this book).
	var displayOrder int
	if in.DisplayOrder != nil {
		displayOrder = *in.DisplayOrder
	} else {
		if err := r.pool.QueryRow(ctx,
			`SELECT COALESCE(MAX(display_order), 0) + 1 FROM chapters WHERE book_id = $1`,
			in.BookID,
		).Scan(&displayOrder); err != nil {
			displayOrder = 1
		}
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO chapters (book_id, chapter_title, file_path, content_summary, display_order, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, book_id, chapter_title, COALESCE(file_path, '') AS file_path,
		          content_summary, display_order, is_visible, created_at
	`, in.BookID, in.ChapterTitle, in.FilePath, in.ContentSummary, displayOrder, in.IsVisible).Scan(
		&ch.ID,
		&ch.BookID,
		&ch.ChapterTitle,
		&ch.FilePath,
		&ch.ContentSummary,
		&ch.DisplayOrder,
		&ch.IsVisible,
		&ch.CreatedAt,
	)
	if err != nil {
		return domain.Chapter{}, err
	}
	return ch, nil
}

// ListByBookID returns chapters for a book ordered by display_order.
func (r *ChapterRepo) ListByBookID(ctx context.Context, bookID int64) ([]domain.Chapter, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, book_id, chapter_title, COALESCE(file_path, '') AS file_path,
		       content_summary, display_order, is_visible, created_at
		FROM chapters
		WHERE book_id = $1
		ORDER BY display_order ASC, id ASC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Chapter
	for rows.Next() {
		var ch domain.Chapter
		if err := rows.Scan(
			&ch.ID,
			&ch.BookID,
			&ch.ChapterTitle,
			&ch.FilePath,
			&ch.ContentSummary,
			&ch.DisplayOrder,
			&ch.IsVisible,
			&ch.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	return out, rows.Err()
}

// Delete deletes a chapter by id. Returns ErrChapterNotFound or ErrChapterHasDependents if applicable.
func (r *ChapterRepo) Delete(ctx context.Context, id int64) error {
	var count int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM generated_content WHERE chapter_id = $1`, id,
	).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return ErrChapterHasDependents
	}

	tag, err := r.pool.Exec(ctx, `DELETE FROM chapters WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrChapterNotFound
	}
	return nil
}

// GetByID returns a chapter by id.
func (r *ChapterRepo) GetByID(ctx context.Context, id int64) (domain.Chapter, error) {
	var ch domain.Chapter
	err := r.pool.QueryRow(ctx, `
		SELECT id, book_id, chapter_title, COALESCE(file_path, '') AS file_path,
		       content_summary, display_order, is_visible, created_at
		FROM chapters WHERE id = $1
	`, id).Scan(
		&ch.ID,
		&ch.BookID,
		&ch.ChapterTitle,
		&ch.FilePath,
		&ch.ContentSummary,
		&ch.DisplayOrder,
		&ch.IsVisible,
		&ch.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Chapter{}, ErrChapterNotFound
		}
		return domain.Chapter{}, err
	}
	return ch, nil
}
