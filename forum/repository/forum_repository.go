package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
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

	if err := fr.db.QueryRow("SELECT id, slug, posts, threads, title, usr FROM forums WHERE slug = $1",
		slug,
	).Scan(&f.ID, &f.Slug, &f.Posts, &f.Threads, &f.Title, &f.User); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return f, nil
}

func (fr *ForumRepository) Create(forum *models.Forum) error {
	return fr.db.QueryRow("INSERT INTO forums (slug, title, usr) VALUES ($1, $2, $3) RETURNING id",
		forum.Slug,
		forum.Title,
		forum.User,
	).Scan(&forum.ID)
}
