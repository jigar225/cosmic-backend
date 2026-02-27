package domain

import "time"

// Language represents a global language definition (e.g. Hindi, English).
type Language struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsVisible bool      `json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateLanguageInput is the payload for creating a language.
type CreateLanguageInput struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsVisible bool   `json:"is_visible"`
}

// UpdateLanguageInput is the payload for partially updating a language.
type UpdateLanguageInput struct {
	Name      *string `json:"name,omitempty"`
	IsVisible *bool   `json:"is_visible,omitempty"`
}

