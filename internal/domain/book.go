package domain

import "time"

// Book represents a curriculum book tied to a subject. Simplified schema: subject_id, title, file_path.
type Book struct {
	ID         int64     `json:"id"`
	SubjectID  int64     `json:"subject_id"`
	Title      string    `json:"title"`
	IsVisible  bool      `json:"is_visible"`
	IsActive   bool      `json:"is_active"`
	FilePath   *string   `json:"file_path,omitempty"`
	CreatedBy  *int64    `json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateBookInput is the payload for creating a book.
type CreateBookInput struct {
	SubjectID int64   `json:"subject_id"`
	Title     string  `json:"title"`
	IsVisible bool    `json:"is_visible"`
	FilePath  *string `json:"file_path,omitempty"`
}

// UpdateBookInput is the payload for partially updating a book.
type UpdateBookInput struct {
	Title     *string `json:"title,omitempty"`
	IsVisible *bool   `json:"is_visible,omitempty"`
	FilePath  *string `json:"file_path,omitempty"`
}
