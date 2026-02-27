package domain

import "time"

// GradeMethod represents a grade system (e.g. "India 1–12", "US K–12").
type GradeMethod struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateGradeMethodInput is the payload for creating a grade method.
type CreateGradeMethodInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	IsVisible   bool    `json:"is_visible"`
}

// UpdateGradeMethodInput is the payload for partially updating a grade method.
type UpdateGradeMethodInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsVisible   *bool   `json:"is_visible,omitempty"`
}
