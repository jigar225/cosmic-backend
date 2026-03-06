package repository

import (
	"context"

	"back_testing/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ProgressRepo provides read-only queries for the progress tree.
type ProgressRepo struct {
	pool *pgxpool.Pool
}

// NewProgressRepo returns a new ProgressRepo.
func NewProgressRepo(pool *pgxpool.Pool) *ProgressRepo {
	return &ProgressRepo{pool: pool}
}

// GetTree returns subjects → books with chapter counts for a given board + medium.
// Optional grade_id filter.
func (r *ProgressRepo) GetTree(ctx context.Context, boardID, mediumID int64, gradeID *int64) (domain.ProgressTree, error) {
	query := `
		SELECT
			s.id          AS subject_id,
			s.title       AS subject_title,
			s.grade_id,
			b.id          AS book_id,
			b.title       AS book_title,
			COALESCE(ch.cnt, 0) AS chapter_count
		FROM subjects s
		LEFT JOIN books b ON b.subject_id = s.id
		LEFT JOIN (
			SELECT book_id, COUNT(*) AS cnt FROM chapters GROUP BY book_id
		) ch ON ch.book_id = b.id
		WHERE s.board_id  = $1
		  AND s.medium_id = $2
	`
	args := []interface{}{boardID, mediumID}
	n := 3
	if gradeID != nil {
		query += " AND s.grade_id = $3"
		args = append(args, *gradeID)
		n = 4
		_ = n // keep compiler happy
	}
	query += " ORDER BY s.title, b.title"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return domain.ProgressTree{}, err
	}
	defer rows.Close()

	subjectMap := map[int64]*domain.ProgressSubject{}
	var subjectOrder []int64

	for rows.Next() {
		var (
			sID, gradeID2 int64
			sTitle        string
			bID           *int64
			bTitle        *string
			chapCount     int
		)
		if err := rows.Scan(&sID, &sTitle, &gradeID2, &bID, &bTitle, &chapCount); err != nil {
			return domain.ProgressTree{}, err
		}

		sub, ok := subjectMap[sID]
		if !ok {
			sub = &domain.ProgressSubject{ID: sID, Title: sTitle, GradeID: gradeID2, Books: []domain.ProgressBook{}}
			subjectMap[sID] = sub
			subjectOrder = append(subjectOrder, sID)
		}
		if bID != nil && bTitle != nil {
			sub.Books = append(sub.Books, domain.ProgressBook{
				ID:           *bID,
				Title:        *bTitle,
				ChapterCount: chapCount,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return domain.ProgressTree{}, err
	}

	subjects := make([]domain.ProgressSubject, 0, len(subjectOrder))
	for _, id := range subjectOrder {
		subjects = append(subjects, *subjectMap[id])
	}
	return domain.ProgressTree{Subjects: subjects}, nil
}
