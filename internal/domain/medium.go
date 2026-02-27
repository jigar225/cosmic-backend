package domain

import "time"

// Medium represents a medium of instruction (language) for a country/board.
type Medium struct {
	ID          int64     `json:"id"`
	CountryID   int64     `json:"country_id"`
	BoardID     *int64    `json:"board_id,omitempty"`
	Title       string    `json:"title"`
	LanguageCode *string   `json:"language_code,omitempty"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMediumInput is the payload for creating a medium.
type CreateMediumInput struct {
	CountryID   int64   `json:"country_id"`
	BoardID     *int64  `json:"board_id,omitempty"`
	Title       string  `json:"title"`
	LanguageCode *string `json:"language_code,omitempty"`
	IsVisible   bool    `json:"is_visible"`
}

// UpdateMediumInput is the payload for partially updating a medium.
type UpdateMediumInput struct {
	Title       *string `json:"title,omitempty"`
	LanguageCode *string `json:"language_code,omitempty"`
	IsVisible   *bool   `json:"is_visible,omitempty"`
}

