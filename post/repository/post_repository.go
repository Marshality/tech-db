package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(conn *sql.DB) post.Repository {
	return &PostRepository{
		db: conn,
	}
}

func (pr *PostRepository) Insert(p *models.Post) error {
	return pr.db.QueryRow(queries.InsertIntoPosts, p.Author, p.Forum, p.Thread, p.Message, p.Parent, p.CreatedAt).Scan(&p.ID)
}

func (pr *PostRepository) SelectByThreadAndID(id, thread uint64) (*models.Post, error) {
	p := &models.Post{}

	if err := pr.db.QueryRow(queries.SelectPostByThreadAndID, thread, id).Scan(
		&p.ID, &p.Forum, &p.Thread, &p.Author, &p.Message, &p.Parent, &p.IsEdited, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return p, nil
}
