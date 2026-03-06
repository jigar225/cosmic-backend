package domain

// ProgressBook represents a book's chapter-upload progress.
type ProgressBook struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	ChapterCount int    `json:"chapter_count"`
}

// ProgressSubject groups books under a subject inside the progress tree.
type ProgressSubject struct {
	ID      int64          `json:"id"`
	Title   string         `json:"title"`
	GradeID int64          `json:"grade_id"`
	Books   []ProgressBook `json:"books"`
}

// ProgressTree is the top-level response for GET /admin/progress/tree.
type ProgressTree struct {
	Subjects []ProgressSubject `json:"subjects"`
}
