package handlers

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"back_testing/internal/domain"
	"back_testing/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxChapterPDFSize = 50 * 1024 * 1024 // 50 MB
const pdfContentType = "application/pdf"

// ListChapters handles GET /admin/books/:book_id/chapters.
func (h *Handlers) ListChapters(c *fiber.Ctx) error {
	bookIDStr := c.Params("book_id")
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil || bookID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid book_id"})
	}

	list, err := h.ChapterRepo.ListByBookID(c.Context(), bookID)
	if err != nil {
		log.Printf("chapters list: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list chapters"})
	}
	if list == nil {
		list = []domain.Chapter{}
	}
	return c.JSON(list)
}

// CreateChapter handles POST /admin/books/:book_id/chapters (multipart: chapter_title, file).
func (h *Handlers) CreateChapter(c *fiber.Ctx) error {
	bookIDStr := c.Params("book_id")
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil || bookID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid book_id"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid multipart form"})
	}

	title := strings.TrimSpace(strings.Join(form.Value["chapter_title"], ""))
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "chapter_title is required"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required (PDF)"})
	}
	fileHeader := files[0]
	if fileHeader.Size > maxChapterPDFSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file too large (max 50MB)"})
	}

	// Validate PDF by extension or content-type
	filename := strings.ToLower(fileHeader.Filename)
	ct := fileHeader.Header.Get("Content-Type")
	if !strings.HasSuffix(filename, ".pdf") && ct != pdfContentType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file must be a PDF"})
	}

	if h.S3Uploader == nil || h.S3Uploader.Bucket() == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "file storage (S3) is not configured"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("chapters create: open file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file"})
	}
	defer file.Close()

	// Limit body size to the actual file size (already validated against maxChapterPDFSize).
	size := fileHeader.Size
	if size > maxChapterPDFSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file too large (max 50MB)"})
	}
	limitedBody := io.LimitReader(file, size)

	key := "chapters/" + strconv.FormatInt(bookID, 10) + "/" + uuid.New().String() + ".pdf"
	if err := h.S3Uploader.Upload(c.Context(), key, pdfContentType, limitedBody, size); err != nil {
		log.Printf("chapters create: s3 upload: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to upload file to storage"})
	}

	in := domain.CreateChapterInput{
		BookID:       bookID,
		ChapterTitle: title,
		FilePath:     key,
		IsVisible:    true,
	}
	if orderStr := strings.TrimSpace(strings.Join(form.Value["display_order"], "")); orderStr != "" {
		if orderVal, err := strconv.Atoi(orderStr); err == nil && orderVal >= 0 {
			in.DisplayOrder = &orderVal
		}
	}
	ch, err := h.ChapterRepo.Create(c.Context(), in)
	if err != nil {
		log.Printf("chapters create: db insert: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create chapter"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"chapter": ch})
}

// DeleteChapter handles DELETE /admin/chapters/:id.
func (h *Handlers) DeleteChapter(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chapter id"})
	}
	if err := h.ChapterRepo.Delete(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrChapterNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "chapter not found"})
		}
		if errors.Is(err, repository.ErrChapterHasDependents) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "chapter cannot be deleted: it has generated content"})
		}
		log.Printf("chapters delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete chapter"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetChapterDownloadURL handles GET /admin/chapters/:id/download-url.
// Returns a presigned URL (valid 15 min) for viewing/downloading the chapter PDF.
func (h *Handlers) GetChapterDownloadURL(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chapter id"})
	}

	ch, err := h.ChapterRepo.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrChapterNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "chapter not found"})
		}
		log.Printf("chapters download-url: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get chapter"})
	}

	if ch.FilePath == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "chapter has no file"})
	}

	if h.S3Uploader == nil || h.S3Uploader.Bucket() == "" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "file storage (S3) is not configured"})
	}

	url, err := h.S3Uploader.PresignGetURL(c.Context(), ch.FilePath)
	if err != nil {
		log.Printf("chapters download-url: presign: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate download URL"})
	}

	return c.JSON(fiber.Map{"url": url})
}
