package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListBoards handles GET /boards (public: visible boards only).
func (h *Handlers) ListBoards(c *fiber.Ctx) error {
	list, err := h.BoardRepo.List(c.Context())
	if err != nil {
		log.Printf("boards list: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list boards"})
	}
	if list == nil {
		list = []domain.Board{}
	}
	return c.JSON(list)
}

// ListAllBoards handles GET /admin/boards. Optional query: ?country_id=1
func (h *Handlers) ListAllBoards(c *fiber.Ctx) error {
	var countryID *int64
	if q := c.Query("country_id"); q != "" {
		id, err := strconv.ParseInt(q, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid country_id"})
		}
		countryID = &id
	}
	list, err := h.BoardRepo.ListAll(c.Context(), countryID)
	if err != nil {
		log.Printf("boards list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list boards"})
	}
	if list == nil {
		list = []domain.Board{}
	}
	return c.JSON(list)
}

// ListBoardsByCountry handles GET /admin/countries/:id/boards.
func (h *Handlers) ListBoardsByCountry(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid country id"})
	}
	list, err := h.BoardRepo.ListByCountryID(c.Context(), id)
	if err != nil {
		log.Printf("boards list by country: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list boards"})
	}
	if list == nil {
		list = []domain.Board{}
	}
	return c.JSON(list)
}

// CreateBoard handles POST /admin/boards.
func (h *Handlers) CreateBoard(c *fiber.Ctx) error {
	var in domain.CreateBoardInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}
	if in.CountryID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "country_id is required and must be positive"})
	}
	board, err := h.BoardRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrBoardConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a board with this title already exists for this country"})
		}
		log.Printf("boards create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create board"})
	}
	return c.Status(fiber.StatusCreated).JSON(board)
}

// UpdateBoard handles PATCH /admin/boards/:id.
func (h *Handlers) UpdateBoard(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var in domain.UpdateBoardInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	board, err := h.BoardRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "board not found"})
		}
		log.Printf("boards update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update board"})
	}
	return c.JSON(board)
}

// DeleteBoard handles DELETE /admin/boards/:id.
func (h *Handlers) DeleteBoard(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.BoardRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "board not found"})
		}
		if errors.Is(err, repository.ErrBoardHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "board cannot be deleted: it has states, mediums, subjects, or other data linked to it"})
		}
		log.Printf("boards delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete board"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
