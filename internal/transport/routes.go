package transport

import (
	"back_testing/internal/transport/handlers"
	"back_testing/internal/transport/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all HTTP routes on the Fiber app.
func RegisterRoutes(app *fiber.App, h *handlers.Handlers) {
	app.Get("/health", h.Health)

	// Auth (public)
	app.Post("/auth/signup", h.Signup)
	app.Post("/auth/login", h.Login)
	app.Post("/auth/refresh", h.Refresh)
	app.Post("/auth/logout", h.Logout)

	// Auth (protected)
	authMw := middleware.RequireAuth(middleware.AuthConfig{JWTSecret: getJWTSecret(h)})
	app.Get("/me", authMw, h.Me)
	app.Post("/auth/logout-all", authMw, h.LogoutAll)

	// Public curriculum configuration
	app.Get("/countries", h.ListPublicCountries)
	app.Get("/boards", h.ListBoards)
	app.Get("/mediums", h.ListMediums)
	app.Get("/grade-methods", h.ListGradeMethods)
	app.Get("/grade-methods/:id/grades", h.ListGradesByMethod)
	app.Get("/subjects", h.ListSubjects)
	app.Get("/books", h.ListBooks)

	// Admin group: must be authenticated + admin role.
	admin := app.Group("/admin", authMw, h.RequireAdmin())

	// Admin: countries
	admin.Get("/countries", h.ListAllCountries)
	admin.Post("/countries", h.CreateCountry)
	admin.Patch("/countries/:id", h.UpdateCountry)

	// Admin: boards (list by country or all; create; update)
	admin.Get("/boards", h.ListAllBoards)
	admin.Get("/countries/:id/boards", h.ListBoardsByCountry)
	admin.Post("/boards", h.CreateBoard)
	admin.Patch("/boards/:id", h.UpdateBoard)
	admin.Delete("/boards/:id", h.DeleteBoard)

	// Admin: grade methods
	admin.Get("/grade-methods", h.ListAllGradeMethods)
	admin.Post("/grade-methods", h.CreateGradeMethod)
	admin.Patch("/grade-methods/:id", h.UpdateGradeMethod)
	admin.Delete("/grade-methods/:id", h.DeleteGradeMethod)

	// Admin: languages
	admin.Get("/languages", h.ListAllLanguages)
	admin.Post("/languages", h.CreateLanguage)
	admin.Patch("/languages/:id", h.UpdateLanguage)
	admin.Delete("/languages/:id", h.DeleteLanguage)

	// Admin: grades
	admin.Get("/grades", h.ListAllGrades)
	admin.Post("/grades", h.CreateGrade)
	admin.Patch("/grades/:id", h.UpdateGrade)
	admin.Delete("/grades/:id", h.DeleteGrade)

	// Admin: mediums
	admin.Get("/mediums", h.ListAllMediums)
	admin.Post("/mediums", h.CreateMedium)
	admin.Patch("/mediums/:id", h.UpdateMedium)
	admin.Delete("/mediums/:id", h.DeleteMedium)

	// Admin: subjects
	admin.Get("/subjects", h.ListAllSubjects)
	admin.Post("/subjects", h.CreateSubject)
	admin.Patch("/subjects/:id", h.UpdateSubject)
	admin.Delete("/subjects/:id", h.DeleteSubject)

	// Admin: books
	admin.Get("/books", h.ListAllBooks)
	admin.Post("/books", h.CreateBook)
	admin.Patch("/books/:id", h.UpdateBook)
	admin.Delete("/books/:id", h.DeleteBook)

	// Admin: progress tree
	app.Get("/admin/progress/tree", h.GetProgressTree)

	// Admin: chapters (PDF upload to S3)
	admin.Get("/books/:book_id/chapters", h.ListChapters)
	admin.Post("/books/:book_id/chapters", h.CreateChapter)
	admin.Patch("/chapters/:id", h.UpdateChapterPdf)
	admin.Delete("/chapters/:id", h.DeleteChapter)
	admin.Get("/chapters/:id/download-url", h.GetChapterDownloadURL)
}

// getJWTSecret returns the JWT secret for auth middleware. Nil if auth not configured.
func getJWTSecret(h *handlers.Handlers) []byte {
	if h.AuthConfig == nil {
		return nil
	}
	return h.AuthConfig.JWTSecret
}
