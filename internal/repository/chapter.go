package repository

import (
	"context"
	"errors"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrChapterNotFound = errors.New("chapter not found")

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
	err := r.pool.QueryRow(ctx, `
		INSERT INTO chapters (book_id, chapter_title, file_path, content_summary, is_visible)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, book_id, chapter_title, file_path, content_summary, is_visible, created_at
	`, in.BookID, in.ChapterTitle, in.FilePath, in.ContentSummary, in.IsVisible).Scan(
		&ch.ID,
		&ch.BookID,
		&ch.ChapterTitle,
		&ch.FilePath,
		&ch.ContentSummary,
		&ch.IsVisible,
		&ch.CreatedAt,
	)
	if err != nil {
		return domain.Chapter{}, err
	}
	return ch, nil
}

// ListByBookID returns chapters for a book.
func (r *ChapterRepo) ListByBookID(ctx context.Context, bookID int64) ([]domain.Chapter, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, book_id, chapter_title, file_path, content_summary, is_visible, created_at
		FROM chapters
		WHERE book_id = $1
		ORDER BY id ASC
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
			&ch.IsVisible,
			&ch.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	return out, rows.Err()
}

// GetByID returns a chapter by id.
func (r *ChapterRepo) GetByID(ctx context.Context, id int64) (domain.Chapter, error) {
	var ch domain.Chapter
	err := r.pool.QueryRow(ctx, `
		SELECT id, book_id, chapter_title, file_path, content_summary, is_visible, created_at
		FROM chapters WHERE id = $1
	`, id).Scan(
		&ch.ID,
		&ch.BookID,
		&ch.ChapterTitle,
		&ch.FilePath,
		&ch.ContentSummary,
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
