package domain

import "time"

// Chapter represents a chapter (PDF) under a book.
type Chapter struct {
	ID            int64     `json:"id"`
	BookID        int64     `json:"book_id"`
	ChapterTitle  string    `json:"chapter_title"`
	FilePath      string    `json:"file_path"` // S3 object key
	ContentSummary *string  `json:"content_summary,omitempty"`
	IsVisible     bool      `json:"is_visible"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateChapterInput is the payload for creating a chapter (title + file_path set after S3 upload).
type CreateChapterInput struct {
	BookID         int64
	ChapterTitle   string
	FilePath       string // S3 key after upload
	ContentSummary *string
	IsVisible      bool
}
