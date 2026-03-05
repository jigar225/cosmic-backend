package transport

import (
	"back_testing/internal/transport/handlers"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all HTTP routes on the Fiber app.
func RegisterRoutes(app *fiber.App, h *handlers.Handlers) {
	app.Get("/health", h.Health)

	// Public curriculum configuration
	app.Get("/countries", h.ListPublicCountries)
	app.Get("/boards", h.ListBoards)
	app.Get("/mediums", h.ListMediums)
	app.Get("/grade-methods", h.ListGradeMethods)
	app.Get("/grade-methods/:id/grades", h.ListGradesByMethod)
	app.Get("/subjects", h.ListSubjects)
	app.Get("/books", h.ListBooks)

	// Admin: countries
	app.Get("/admin/countries", h.ListAllCountries)
	app.Post("/admin/countries", h.CreateCountry)
	app.Patch("/admin/countries/:id", h.UpdateCountry)

	// Admin: boards (list by country or all; create; update)
	app.Get("/admin/boards", h.ListAllBoards)
	app.Get("/admin/countries/:id/boards", h.ListBoardsByCountry)
	app.Post("/admin/boards", h.CreateBoard)
	app.Patch("/admin/boards/:id", h.UpdateBoard)
	app.Delete("/admin/boards/:id", h.DeleteBoard)

	// Admin: grade methods
	app.Get("/admin/grade-methods", h.ListAllGradeMethods)
	app.Post("/admin/grade-methods", h.CreateGradeMethod)
	app.Patch("/admin/grade-methods/:id", h.UpdateGradeMethod)
	app.Delete("/admin/grade-methods/:id", h.DeleteGradeMethod)

	// Admin: languages
	app.Get("/admin/languages", h.ListAllLanguages)
	app.Post("/admin/languages", h.CreateLanguage)
	app.Patch("/admin/languages/:id", h.UpdateLanguage)
	app.Delete("/admin/languages/:id", h.DeleteLanguage)

	// Admin: grades
	app.Get("/admin/grades", h.ListAllGrades)
	app.Post("/admin/grades", h.CreateGrade)
	app.Patch("/admin/grades/:id", h.UpdateGrade)
	app.Delete("/admin/grades/:id", h.DeleteGrade)

	// Admin: mediums
	app.Get("/admin/mediums", h.ListAllMediums)
	app.Post("/admin/mediums", h.CreateMedium)
	app.Patch("/admin/mediums/:id", h.UpdateMedium)
	app.Delete("/admin/mediums/:id", h.DeleteMedium)

	// Admin: subjects
	app.Get("/admin/subjects", h.ListAllSubjects)
	app.Post("/admin/subjects", h.CreateSubject)
	app.Patch("/admin/subjects/:id", h.UpdateSubject)
	app.Delete("/admin/subjects/:id", h.DeleteSubject)

	// Admin: books
	app.Get("/admin/books", h.ListAllBooks)
	app.Post("/admin/books", h.CreateBook)
	app.Patch("/admin/books/:id", h.UpdateBook)
	app.Delete("/admin/books/:id", h.DeleteBook)

	// Admin: chapters (PDF upload to S3)
	app.Get("/admin/books/:book_id/chapters", h.ListChapters)
	app.Post("/admin/books/:book_id/chapters", h.CreateChapter)
	app.Get("/admin/chapters/:id/download-url", h.GetChapterDownloadURL)
}
