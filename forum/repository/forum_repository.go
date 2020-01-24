package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
	"log"
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

func (fr *ForumRepository) SelectForumUsers(slug string, limit uint64, since string, desc bool) ([]*models.User, error) {
	var users []*models.User

	var rows *sql.Rows
	var err error

	if since != "" {
		rows, err = fr.db.Query(queries.QM.SelectForumUsersSince(desc, limit), slug, since)
	} else {
		rows, err = fr.db.Query(queries.QM.SelectForumUsers(desc, limit), slug)
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
		u := &models.User{}

		if err := rows.Scan(&u.ID, &u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
