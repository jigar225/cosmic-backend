package domain

import "time"

// Subject represents a school subject within a specific country/board/medium/grade.
type Subject struct {
	ID           int64     `json:"id"`
	CountryID    int64     `json:"country_id"`
	BoardID      int64     `json:"board_id"`
	MediumID     int64     `json:"medium_id"`
	GradeID      int64     `json:"grade_id"`
	Title        string    `json:"title"`
	SubjectType  string    `json:"subject_type"`
	IsVisible    bool      `json:"is_visible"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateSubjectInput is the payload for creating a subject.
type CreateSubjectInput struct {
	CountryID    int64   `json:"country_id"`
	BoardID      int64   `json:"board_id"`
	MediumID     int64   `json:"medium_id"`
	GradeID      int64   `json:"grade_id"`
	Title        string  `json:"title"`
	SubjectType  string  `json:"subject_type"`
	IsVisible    bool    `json:"is_visible"`
}

// UpdateSubjectInput is the payload for partially updating a subject.
type UpdateSubjectInput struct {
	Title        *string `json:"title,omitempty"`
	SubjectType  *string `json:"subject_type,omitempty"`
	IsVisible    *bool   `json:"is_visible,omitempty"`
}

