package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
	"log"
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

func (pr *PostRepository) GetPostsByThread(thread uint64, since uint64, limit uint64, desc bool, sort string) ([]*models.Post, error) {
	var posts []*models.Post

	var rows *sql.Rows
	var err error

	if since != 0 {
		rows, err = pr.db.Query(queries.QM.SelectPostsByThreadSince(desc, sort), thread, since, limit)
	} else {
		rows, err = pr.db.Query(queries.QM.SelectPostsByThread(desc, sort), thread, limit)
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		p := &models.Post{}

		if err := rows.Scan(&p.ID, &p.Author, &p.Forum, &p.Thread, &p.Message, &p.Parent, &p.IsEdited, &p.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}
