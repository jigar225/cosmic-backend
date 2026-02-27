package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListAllLanguages handles GET /admin/languages.
func (h *Handlers) ListAllLanguages(c *fiber.Ctx) error {
	list, err := h.LanguageRepo.ListAll(c.Context())
	if err != nil {
		log.Printf("languages list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list languages"})
	}
	if list == nil {
		list = []domain.Language{}
	}
	return c.JSON(list)
}

// CreateLanguage handles POST /admin/languages.
func (h *Handlers) CreateLanguage(c *fiber.Ctx) error {
	var in domain.CreateLanguageInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.Code == "" || in.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "code and name are required"})
	}
	lang, err := h.LanguageRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrLanguageConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "language with this code already exists"})
		}
		log.Printf("languages create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create language"})
	}
	return c.Status(fiber.StatusCreated).JSON(lang)
}

// UpdateLanguage handles PATCH /admin/languages/:id.
func (h *Handlers) UpdateLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var in domain.UpdateLanguageInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	lang, err := h.LanguageRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrLanguageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "language not found"})
		}
		log.Printf("languages update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update language"})
	}
	return c.JSON(lang)
}

// DeleteLanguage handles DELETE /admin/languages/:id.
func (h *Handlers) DeleteLanguage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.LanguageRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLanguageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "language not found"})
		}
		if errors.Is(err, repository.ErrLanguageHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "language cannot be deleted: it is used by mediums"})
		}
		log.Printf("languages delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete language"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}


