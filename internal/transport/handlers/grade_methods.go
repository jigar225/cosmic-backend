package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListGradeMethods handles GET /grade-methods (public: visible only).
func (h *Handlers) ListGradeMethods(c *fiber.Ctx) error {
	list, err := h.GradeMethodRepo.ListVisible(c.Context())
	if err != nil {
		log.Printf("grade_methods list visible: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list grade methods"})
	}
	if list == nil {
		list = []domain.GradeMethod{}
	}
	return c.JSON(list)
}

// ListAllGradeMethods handles GET /admin/grade-methods.
func (h *Handlers) ListAllGradeMethods(c *fiber.Ctx) error {
	list, err := h.GradeMethodRepo.ListAll(c.Context())
	if err != nil {
		log.Printf("grade_methods list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list grade methods"})
	}
	if list == nil {
		list = []domain.GradeMethod{}
	}
	return c.JSON(list)
}

// CreateGradeMethod handles POST /admin/grade-methods.
func (h *Handlers) CreateGradeMethod(c *fiber.Ctx) error {
	var in domain.CreateGradeMethodInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	g, err := h.GradeMethodRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrGradeMethodConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a grade method with this title already exists"})
		}
		log.Printf("grade_methods create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create grade method"})
	}
	return c.Status(fiber.StatusCreated).JSON(g)
}

// UpdateGradeMethod handles PATCH /admin/grade-methods/:id.
func (h *Handlers) UpdateGradeMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var in domain.UpdateGradeMethodInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	g, err := h.GradeMethodRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrGradeMethodNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "grade method not found"})
		}
		log.Printf("grade_methods update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update grade method"})
	}
	return c.JSON(g)
}

// DeleteGradeMethod handles DELETE /admin/grade-methods/:id.
func (h *Handlers) DeleteGradeMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.GradeMethodRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrGradeMethodNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "grade method not found"})
		}
		if errors.Is(err, repository.ErrGradeMethodInUse) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "grade method cannot be deleted: it is used by boards or grades"})
		}
		log.Printf("grade_methods delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete grade method"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
