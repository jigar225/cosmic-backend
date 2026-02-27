package handlers

import "back_testing/internal/repository"

// Handlers holds dependencies for HTTP handlers (DB, uploader, etc.).
type Handlers struct {
	BoardRepo       *repository.BoardRepo
	CountryRepo     *repository.CountryRepo
	GradeMethodRepo *repository.GradeMethodRepo
	GradeRepo       *repository.GradeRepo
	MediumRepo      *repository.MediumRepo
	LanguageRepo    *repository.LanguageRepo
	SubjectRepo     *repository.SubjectRepo
}
