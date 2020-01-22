package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(conn *sql.DB) thread.Repository {
	return &ThreadRepository{
		db: conn,
	}
}

func (tr *ThreadRepository) SelectBySlug(slug string) (*models.Thread, error) {
	t := &models.Thread{}

	if err := tr.db.QueryRow(queries.SelectThreadWhereSlug, slug).Scan(
		&t.ID, &t.Slug, &t.Author, &t.Forum, &t.Message, &t.Title, &t.Votes, &t.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return t, nil
}

func (tr *ThreadRepository) Insert(thread *models.Thread) error {
	return tr.db.QueryRow(queries.InsertIntoThreads,
		thread.Slug,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Title,
		thread.CreatedAt,
	).Scan(&thread.ID)
}
