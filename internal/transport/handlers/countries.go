package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListPublicCountries handles GET /countries for non-admin consumers.
func (h *Handlers) ListPublicCountries(c *fiber.Ctx) error {
	list, err := h.CountryRepo.ListVisible(c.Context())
	if err != nil {
		log.Printf("countries list visible: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list countries"})
	}
	if list == nil {
		list = []domain.Country{}
	}
	return c.JSON(list)
}

// ListAllCountries handles GET /admin/countries for admins (visible + hidden).
func (h *Handlers) ListAllCountries(c *fiber.Ctx) error {
	list, err := h.CountryRepo.ListAll(c.Context())
	if err != nil {
		log.Printf("countries list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list countries"})
	}
	if list == nil {
		list = []domain.Country{}
	}
	return c.JSON(list)
}

// CreateCountry handles POST /admin/countries.
func (h *Handlers) CreateCountry(c *fiber.Ctx) error {
	var in domain.CreateCountryInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.CountryCode == "" || in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_code and title are required"})
	}
	if in.SignupMethod == "" {
		in.SignupMethod = "email"
	}
	country, err := h.CountryRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrCountryConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "country with this country_code already exists"})
		}
		log.Printf("countries create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create country"})
	}
	return c.Status(fiber.StatusCreated).JSON(country)
}

// UpdateCountry handles PATCH /admin/countries/:id to update fields including is_visible.
func (h *Handlers) UpdateCountry(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var in domain.UpdateCountryInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	country, err := h.CountryRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrCountryNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "country not found"})
		}
		log.Printf("countries update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update country"})
	}
	return c.JSON(country)
}

