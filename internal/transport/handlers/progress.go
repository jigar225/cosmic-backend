package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetProgressTree handles GET /admin/progress/tree?board_id=&medium_id=&grade_id=.
func (h *Handlers) GetProgressTree(c *fiber.Ctx) error {
	boardStr := c.Query("board_id")
	mediumStr := c.Query("medium_id")
	if boardStr == "" || mediumStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "board_id and medium_id are required"})
	}

	boardID, err := strconv.ParseInt(boardStr, 10, 64)
	if err != nil || boardID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid board_id"})
	}
	mediumID, err := strconv.ParseInt(mediumStr, 10, 64)
	if err != nil || mediumID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medium_id"})
	}

	var gradeID *int64
	if gs := c.Query("grade_id"); gs != "" {
		id, err := strconv.ParseInt(gs, 10, 64)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid grade_id"})
		}
		gradeID = &id
	}

	tree, err := h.ProgressRepo.GetTree(c.Context(), boardID, mediumID, gradeID)
	if err != nil {
		log.Printf("progress tree: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to build progress tree"})
	}

	return c.JSON(tree)
}
