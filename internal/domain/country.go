package domain

import "time"

// Country represents a nation configuration in the curriculum system.
type Country struct {
	ID          int64     `json:"id"`
	CountryCode string    `json:"country_code"`
	Title       string    `json:"title"`
	IconPath    *string   `json:"icon_path,omitempty"`
	PhoneCode   *string   `json:"phone_code,omitempty"`
	SignupMethod string   `json:"signup_method"`
	HaveBoard   bool      `json:"have_board"`
	HasStates   bool      `json:"has_states"`
	IsVisible   bool      `json:"is_visible"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCountryInput is the payload for creating a country from admin.
type CreateCountryInput struct {
	CountryCode  string  `json:"country_code"`
	Title        string  `json:"title"`
	IconPath     *string `json:"icon_path,omitempty"`
	PhoneCode    *string `json:"phone_code,omitempty"`
	SignupMethod string  `json:"signup_method"`
	HaveBoard    bool    `json:"have_board"`
	HasStates    bool    `json:"has_states"`
	IsVisible    bool    `json:"is_visible"`
}

// UpdateCountryInput is the payload for partially updating a country.
// All fields are optional; only non-nil pointers are applied.
type UpdateCountryInput struct {
	Title        *string `json:"title,omitempty"`
	IconPath     *string `json:"icon_path,omitempty"`
	PhoneCode    *string `json:"phone_code,omitempty"`
	SignupMethod *string `json:"signup_method,omitempty"`
	HaveBoard    *bool   `json:"have_board,omitempty"`
	HasStates    *bool   `json:"has_states,omitempty"`
	IsVisible    *bool   `json:"is_visible,omitempty"`
}

