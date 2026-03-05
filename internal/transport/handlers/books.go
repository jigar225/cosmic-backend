package handlers

import (
	"errors"
	"log"
	"strconv"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// ListBooks handles GET /books?subject_id= (public: visible only).
func (h *Handlers) ListBooks(c *fiber.Ctx) error {
	subjectStr := c.Query("subject_id")
	if subjectStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "subject_id is required"})
	}
	subjectID, err := strconv.ParseInt(subjectStr, 10, 64)
	if err != nil || subjectID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject_id"})
	}

	list, err := h.BookRepo.ListVisible(c.Context(), subjectID)
	if err != nil {
		log.Printf("books list visible: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list books"})
	}
	if list == nil {
		list = []domain.Book{}
	}
	return c.JSON(list)
}

// ListAllBooks handles GET /admin/books with optional subject_id filter.
func (h *Handlers) ListAllBooks(c *fiber.Ctx) error {
	var subjectID *int64
	if ss := c.Query("subject_id"); ss != "" {
		id, err := strconv.ParseInt(ss, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject_id"})
		}
		subjectID = &id
	}

	list, err := h.BookRepo.ListAll(c.Context(), subjectID)
	if err != nil {
		log.Printf("books list all: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list books"})
	}
	if list == nil {
		list = []domain.Book{}
	}
	return c.JSON(list)
}

// CreateBook handles POST /admin/books.
func (h *Handlers) CreateBook(c *fiber.Ctx) error {
	var in domain.CreateBookInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	if in.SubjectID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "subject_id is required and must be positive"})
	}
	if in.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title is required"})
	}

	b, err := h.BookRepo.Create(c.Context(), in)
	if err != nil {
		if errors.Is(err, repository.ErrBookConflict) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "a book with this title already exists for this subject"})
		}
		log.Printf("books create: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create book"})
	}
	return c.Status(fiber.StatusCreated).JSON(b)
}

// UpdateBook handles PATCH /admin/books/:id.
func (h *Handlers) UpdateBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var in domain.UpdateBookInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	b, err := h.BookRepo.Update(c.Context(), id, in)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
		}
		log.Printf("books update: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update book"})
	}
	return c.JSON(b)
}

// DeleteBook handles DELETE /admin/books/:id.
func (h *Handlers) DeleteBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.BookRepo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "book not found"})
		}
		if errors.Is(err, repository.ErrBookHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "book cannot be deleted: it is used by chapters or generated content"})
		}
		log.Printf("books delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete book"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

