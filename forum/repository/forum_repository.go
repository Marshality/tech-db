package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(conn *sql.DB) forum.Repository {
	return &ForumRepository{
		db: conn,
	}
}

func (fr *ForumRepository) SelectBySlug(slug string) (*models.Forum, error) {
	f := &models.Forum{}

	if err := fr.db.QueryRow(queries.SelectForumWhereSlug, slug).Scan(
		&f.ID, &f.Slug, &f.Posts, &f.Threads, &f.Title, &f.User); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return f, nil
}

func (fr *ForumRepository) Insert(forum *models.Forum) error {
	return fr.db.QueryRow(queries.InsertIntoForums,
		forum.Slug,
		forum.Title,
		forum.User,
	).Scan(&forum.ID)
}

func (fr *ForumRepository) InsertThread(thread *models.Thread) error {
	return fr.db.QueryRow(queries.InsertIntoThreads,
		thread.Slug,
		thread.Author,
		thread.Message,
		thread.Title,
		thread.CreatedAt,
	).Scan(&thread.ID)
}
