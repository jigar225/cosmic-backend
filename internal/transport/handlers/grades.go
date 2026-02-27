package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListGradesByMethod handles GET /grade-methods/:id/grades (public: visible only).
func (h *Handlers) ListGradesByMethod(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid grade_method_id"})
	}
	list, err := h.GradeRepo.ListByGradeMethod(c.Context(), id, true)
	if err != nil {
		log.Printf("grades list by method (public): %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list grades"})
	}
	if list == nil {
		list = []domain.Grade{}
	}
	return c.JSON(list)
}

// ListAllGrades handles GET /admin/grades. Optional ?grade_method_id=1
func (h *Handlers) ListAllGrades(c *fiber.Ctx) error {
	var methodID *int64
	if q := c.Query("grade_method_id"); q != "" {
		id, err := strconv.ParseInt(q, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid grade_method_id"})
		}
		methodID = &id
	}
	list, err := h.GradeRepo.ListAll(c.Context(), methodID)
	if err != nil {
		log.Printf("grades list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list grades"})
	}
	if list == nil {
		list = []domain.Grade{}
	}
	return c.JSON(list)
}

// CreateGrade handles POST /admin/grades.
func (h *Handlers) CreateGrade(c *fiber.Ctx) error {
	var in domain.CreateGradeInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.GradeMethodID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "grade_method_id is required and must be positive"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	if in.DisplayOrder == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "display_order is required and must be non-zero"})
	}
	g, err := h.GradeRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrGradeConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a grade with this title already exists for this grade method"})
		}
		log.Printf("grades create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create grade"})
	}
	return c.Status(fiber.StatusCreated).JSON(g)
}

// UpdateGrade handles PATCH /admin/grades/:id.
func (h *Handlers) UpdateGrade(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var in domain.UpdateGradeInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	g, err := h.GradeRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrGradeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "grade not found"})
		}
		log.Printf("grades update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update grade"})
	}
	return c.JSON(g)
}

// DeleteGrade handles DELETE /admin/grades/:id.
func (h *Handlers) DeleteGrade(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.GradeRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrGradeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "grade not found"})
		}
		if errors.Is(err, repository.ErrGradeHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "grade cannot be deleted: it is used by subjects, books, user context, or generated content"})
		}
		log.Printf("grades delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete grade"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

