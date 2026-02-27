package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListSubjects handles GET /subjects?country_id=&board_id=&medium_id=&grade_id= (public: visible only).
func (h *Handlers) ListSubjects(c *fiber.Ctx) error {
	countryStr := c.Query("country_id")
	boardStr := c.Query("board_id")
	if countryStr == "" || boardStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_id and board_id are required"})
	}
	countryID, err := strconv.ParseInt(countryStr, 10, 64)
	if err != nil || countryID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid country_id"})
	}
	boardID, err := strconv.ParseInt(boardStr, 10, 64)
	if err != nil || boardID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid board_id"})
	}

	var mediumID *int64
	if m := c.Query("medium_id"); m != "" {
		id, err := strconv.ParseInt(m, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medium_id"})
		}
		mediumID = &id
	}
	var gradeID *int64
	if g := c.Query("grade_id"); g != "" {
		id, err := strconv.ParseInt(g, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid grade_id"})
		}
		gradeID = &id
	}

	list, err := h.SubjectRepo.ListVisible(c.Context(), countryID, boardID, mediumID, gradeID)
	if err != nil {
		log.Printf("subjects list visible: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list subjects"})
	}
	if list == nil {
		list = []domain.Subject{}
	}
	return c.JSON(list)
}

// ListAllSubjects handles GET /admin/subjects with optional filters.
func (h *Handlers) ListAllSubjects(c *fiber.Ctx) error {
	var countryID *int64
	var boardID *int64
	var mediumID *int64
	var gradeID *int64

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
	if ms := c.Query("medium_id"); ms != "" {
		id, err := strconv.ParseInt(ms, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medium_id"})
		}
		mediumID = &id
	}
	if gs := c.Query("grade_id"); gs != "" {
		id, err := strconv.ParseInt(gs, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid grade_id"})
		}
		gradeID = &id
	}

	list, err := h.SubjectRepo.ListAll(c.Context(), countryID, boardID, mediumID, gradeID)
	if err != nil {
		log.Printf("subjects list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list subjects"})
	}
	if list == nil {
		list = []domain.Subject{}
	}
	return c.JSON(list)
}

// CreateSubject handles POST /admin/subjects.
func (h *Handlers) CreateSubject(c *fiber.Ctx) error {
	var in domain.CreateSubjectInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.CountryID <= 0 || in.BoardID <= 0 || in.MediumID <= 0 || in.GradeID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_id, board_id, medium_id, and grade_id are required and must be positive"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	if in.SubjectType == "" {
		in.SubjectType = "core"
	}
	s, err := h.SubjectRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrSubjectConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a subject with this title already exists for this country/board/medium/grade"})
		}
		log.Printf("subjects create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create subject"})
	}
	return c.Status(fiber.StatusCreated).JSON(s)
}

// UpdateSubject handles PATCH /admin/subjects/:id.
func (h *Handlers) UpdateSubject(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var in domain.UpdateSubjectInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	s, err := h.SubjectRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrSubjectNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
		}
		log.Printf("subjects update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update subject"})
	}
	return c.JSON(s)
}

// DeleteSubject handles DELETE /admin/subjects/:id.
func (h *Handlers) DeleteSubject(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.SubjectRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSubjectNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
		}
		if errors.Is(err, repository.ErrSubjectHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "subject cannot be deleted: it is used by books, user context, or generated content"})
		}
		log.Printf("subjects delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete subject"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

