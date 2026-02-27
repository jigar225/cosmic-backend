package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListMediums handles GET /mediums?country_id=&board_id= (public: visible only).
func (h *Handlers) ListMediums(c *fiber.Ctx) error {
	countryStr := c.Query("country_id")
	if countryStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_id is required"})
	}
	countryID, err := strconv.ParseInt(countryStr, 10, 64)
	if err != nil || countryID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid country_id"})
	}
	var boardID *int64
	if b := c.Query("board_id"); b != "" {
		id, err := strconv.ParseInt(b, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid board_id"})
		}
		boardID = &id
	}
	list, err := h.MediumRepo.ListVisible(c.Context(), countryID, boardID)
	if err != nil {
		log.Printf("mediums list visible: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list mediums"})
	}
	if list == nil {
		list = []domain.Medium{}
	}
	return c.JSON(list)
}

// ListAllMediums handles GET /admin/mediums with optional country_id, board_id filters.
func (h *Handlers) ListAllMediums(c *fiber.Ctx) error {
	var countryID *int64
	var boardID *int64

	if cs := c.Query("country_id"); cs != "" {
		id, err := strconv.ParseInt(cs, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid country_id"})
		}
		countryID = &id
	}
	if bs := c.Query("board_id"); bs != "" {
		id, err := strconv.ParseInt(bs, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid board_id"})
		}
		boardID = &id
	}

	list, err := h.MediumRepo.ListAll(c.Context(), countryID, boardID)
	if err != nil {
		log.Printf("mediums list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list mediums"})
	}
	if list == nil {
		list = []domain.Medium{}
	}
	return c.JSON(list)
}

// CreateMedium handles POST /admin/mediums.
func (h *Handlers) CreateMedium(c *fiber.Ctx) error {
	var in domain.CreateMediumInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.CountryID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_id is required and must be positive"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	m, err := h.MediumRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrMediumConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a medium with this title already exists for this country/board"})
		}
		log.Printf("mediums create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create medium"})
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}

// UpdateMedium handles PATCH /admin/mediums/:id.
func (h *Handlers) UpdateMedium(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var in domain.UpdateMediumInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	m, err := h.MediumRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrMediumNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "medium not found"})
		}
		log.Printf("mediums update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update medium"})
	}
	return c.JSON(m)
}

// DeleteMedium handles DELETE /admin/mediums/:id.
func (h *Handlers) DeleteMedium(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.MediumRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrMediumNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "medium not found"})
		}
		if errors.Is(err, repository.ErrMediumHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "medium cannot be deleted: it is used by subjects, books, user context, or generated content"})
		}
		log.Printf("mediums delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete medium"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

