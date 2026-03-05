package domain

import "time"

// Board represents a curriculum authority (e.g. CBSE, ICSE).
type Board struct {
	ID            int64     `json:"id"`
	CountryID     int64     `json:"country_id"`
	Title         string    `json:"title"`
	GradeMethodID *int64    `json:"grade_method_id,omitempty"`
	IsVisible     bool      `json:"is_visible"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateBoardInput is the payload for creating a board.
type CreateBoardInput struct {
	CountryID     int64  `json:"country_id"`
	Title         string `json:"title"`
	GradeMethodID *int64 `json:"grade_method_id,omitempty"`
	IsVisible     bool   `json:"is_visible"`
}

// UpdateBoardInput is the payload for partially updating a board.
type UpdateBoardInput struct {
	Title         *string `json:"title,omitempty"`
	GradeMethodID *int64  `json:"grade_method_id,omitempty"`
	IsVisible     *bool   `json:"is_visible,omitempty"`
}
