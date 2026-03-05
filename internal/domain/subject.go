package domain

import "time"

// Subject represents a school subject within a specific board/medium/grade.
type Subject struct {
	ID        int64     `json:"id"`
	BoardID   int64     `json:"board_id"`
	MediumID  int64     `json:"medium_id"`
	GradeID   int64     `json:"grade_id"`
	Title     string    `json:"title"`
	IsVisible bool      `json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSubjectInput is the payload for creating a subject.
type CreateSubjectInput struct {
	BoardID   int64  `json:"board_id"`
	MediumID  int64  `json:"medium_id"`
	GradeID   int64  `json:"grade_id"`
	Title     string `json:"title"`
	IsVisible bool   `json:"is_visible"`
}

// UpdateSubjectInput is the payload for partially updating a subject.
type UpdateSubjectInput struct {
	Title     *string `json:"title,omitempty"`
	IsVisible *bool   `json:"is_visible,omitempty"`
}

