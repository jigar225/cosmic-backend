package domain

import "time"

// Grade represents a concrete grade within a grade method (e.g. Std 8, Grade 5).
type Grade struct {
	ID            int64     `json:"id"`
	GradeMethodID int64     `json:"grade_method_id"`
	Title         string    `json:"title"`
	AgeRangeStart *int      `json:"age_range_start,omitempty"`
	AgeRangeEnd   *int      `json:"age_range_end,omitempty"`
	IsVisible     bool      `json:"is_visible"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateGradeInput is the payload for creating a grade.
type CreateGradeInput struct {
	GradeMethodID int64  `json:"grade_method_id"`
	Title         string `json:"title"`
	AgeRangeStart *int   `json:"age_range_start,omitempty"`
	AgeRangeEnd   *int   `json:"age_range_end,omitempty"`
	IsVisible     bool   `json:"is_visible"`
}

// UpdateGradeInput is the payload for partially updating a grade.
type UpdateGradeInput struct {
	Title         *string `json:"title,omitempty"`
	AgeRangeStart *int    `json:"age_range_start,omitempty"`
	AgeRangeEnd   *int    `json:"age_range_end,omitempty"`
	IsVisible     *bool   `json:"is_visible,omitempty"`
}

