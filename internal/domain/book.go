package domain

import "time"

// Book represents a curriculum book tied to a subject/medium/grade.
type Book struct {
	ID                int64     `json:"id"`
	BookType          string    `json:"book_type"`
	SubjectID         int64     `json:"subject_id"`
	MediumID          int64     `json:"medium_id"`
	GradeID           int64     `json:"grade_id"`
	Title             string    `json:"title"`
	Author            *string   `json:"author,omitempty"`
	Publisher         *string   `json:"publisher,omitempty"`
	Edition           *string   `json:"edition,omitempty"`
	PublicationYear   *int      `json:"publication_year,omitempty"`
	ISBN              *string   `json:"isbn,omitempty"`
	BookCode          *string   `json:"book_code,omitempty"`
	UploadedByUserID  *int64    `json:"uploaded_by_user_id,omitempty"`
	IsPublic          *bool     `json:"is_public,omitempty"`
	CurriculumVersion *string   `json:"curriculum_version,omitempty"`
	Status            string    `json:"status"`
	OriginalFilePath  string    `json:"original_file_path"`
	ProcessedFilePath *string   `json:"processed_file_path,omitempty"`
	CoverImageURL     *string   `json:"cover_image_url,omitempty"`
	IsVisible         bool      `json:"is_visible"`
	ViewCount         *int      `json:"view_count,omitempty"`
	DownloadCount     *int      `json:"download_count,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateBookInput is the payload for creating a book.
type CreateBookInput struct {
	BookType          string  `json:"book_type"`
	SubjectID         int64   `json:"subject_id"`
	MediumID          int64   `json:"medium_id"`
	GradeID           int64   `json:"grade_id"`
	Title             string  `json:"title"`
	Author            *string `json:"author,omitempty"`
	Publisher         *string `json:"publisher,omitempty"`
	Edition           *string `json:"edition,omitempty"`
	PublicationYear   *int    `json:"publication_year,omitempty"`
	ISBN              *string `json:"isbn,omitempty"`
	BookCode          *string `json:"book_code,omitempty"`
	UploadedByUserID  *int64  `json:"uploaded_by_user_id,omitempty"`
	IsPublic          *bool   `json:"is_public,omitempty"`
	CurriculumVersion *string `json:"curriculum_version,omitempty"`
	Status            string  `json:"status"`
	OriginalFilePath  string  `json:"original_file_path"`
	ProcessedFilePath *string `json:"processed_file_path,omitempty"`
	CoverImageURL     *string `json:"cover_image_url,omitempty"`
	IsVisible         bool    `json:"is_visible"`
}

// UpdateBookInput is the payload for partially updating a book.
type UpdateBookInput struct {
	BookType          *string `json:"book_type,omitempty"`
	Title             *string `json:"title,omitempty"`
	Author            *string `json:"author,omitempty"`
	Publisher         *string `json:"publisher,omitempty"`
	Edition           *string `json:"edition,omitempty"`
	PublicationYear   *int    `json:"publication_year,omitempty"`
	ISBN              *string `json:"isbn,omitempty"`
	BookCode          *string `json:"book_code,omitempty"`
	IsPublic          *bool   `json:"is_public,omitempty"`
	CurriculumVersion *string `json:"curriculum_version,omitempty"`
	Status            *string `json:"status,omitempty"`
	ProcessedFilePath *string `json:"processed_file_path,omitempty"`
	CoverImageURL     *string `json:"cover_image_url,omitempty"`
	IsVisible         *bool   `json:"is_visible,omitempty"`
}

