package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
	"log"
	"sync"
)

type PostRepository struct {
	db *sql.DB
	mux sync.Mutex
}

func NewPostRepository(conn *sql.DB) post.Repository {
	return &PostRepository{
		db: conn,
	}
}

func (pr *PostRepository) Insert(posts *models.Posts, forum string, threadID uint64, createdAt string) error {
	query := queries.InsertIntoPosts

	for i, p := range *posts {
		u := &models.User{}

		if err := pr.db.QueryRow(queries.SelectUserWhereNickname, p.Author).Scan(&u.ID, &u.Nickname, &u.Fullname, &u.Email, &u.About); err != nil {
			return err
		}

		p.Forum = forum
		p.Thread = threadID
		p.CreatedAt = createdAt
		p.Author = u.Nickname

		if p.Parent != 0 {
			if _, err := pr.SelectByThreadAndID(p.Parent, threadID); err != nil {
				return errors.New("conflict")
			}

			query += fmt.Sprintf("(nextval('post_id_seq'::regclass)," +
				"'%s', '%s', %d, '%s', %d, '%s', (SELECT path FROM posts WHERE id = %d AND thread = %d) || " +
				"currval(pg_get_serial_sequence('posts', 'id'))::bigint) ",
				p.Author, p.Forum, p.Thread, p.Message, p.Parent, p.CreatedAt, p.Parent, p.Thread)
		} else {
			query += fmt.Sprintf("(nextval('post_id_seq'::regclass)," +
				"'%s', '%s', %d, '%s', %d, '%s', ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint]) ",
				p.Author, p.Forum, p.Thread, p.Message, p.Parent, p.CreatedAt)
		}


		if i < len(*posts) - 1 {
			query += ", "
		}
	}

	query += "RETURNING id"

	pr.mux.Lock()
	rows, err := pr.db.Query(query)
	pr.mux.Unlock()

	if err != nil {
		return err
	}

	if _, err := pr.db.Exec("UPDATE forums SET posts = posts + $1 WHERE slug = $2", len(*posts), forum); err != nil {
		return err
	}

	var lastInsertedID uint64
	for rows.Next() {
		if err := rows.Scan(&lastInsertedID); err != nil {
			return err
		}
	}

	query = "INSERT INTO user_forum(forum_slug, user_id) VALUES "
	for i, p := range *posts {
		p.ID = lastInsertedID - uint64(len(*posts) - 1 - i)

		if i != 0 {
			query += " , "
		}

		query += fmt.Sprintf(" ('%s', (SELECT id FROM users WHERE nickname = '%s')) ", forum, p.Author)
	}

	query += " ON CONFLICT DO NOTHING "

	if _, err := pr.db.Exec(query); err != nil {
		return err
	}

	return nil
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

// related[0] - user
// related[1] - forum
// related[2] - thread
func (pr *PostRepository) SelectPostWhereID(id uint64, related ...string) (
	*models.Post, *models.Thread, *models.Forum, *models.User, error,
	) {

	p, t, f, u := &models.Post{}, &models.Thread{}, &models.Forum{}, &models.User{}

	var userFlag, threadFlag, forumFlag bool
	for _, value := range related {
		if value == "forum" {
			forumFlag = true
		}

		if value == "thread" {
			threadFlag = true
		}

		if value == "user" {
			userFlag = true
		}
	}

	if !userFlag {
		u = nil
	}

	if !threadFlag {
		t = nil
	}

	if !forumFlag {
		f = nil
	}

	helper, query := queries.QM.SelectPostWhereID(userFlag, threadFlag, forumFlag, related...)

	rows, err := pr.db.Query(query, id)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	columns, _ := rows.Columns()
	columnsCount := len(columns)

	for rows.Next() {
		cols := make([]interface{}, columnsCount)

		for i := 0; i < columnsCount; i++ {
			cols[i] = helper(columns[i], u, f, t, p)
		}

		err = rows.Scan(cols...)

		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return p, t, f, u, nil
}

func (pr *PostRepository) SelectPost(id uint64) (*models.Post, error) {
	p := &models.Post{}

	if err := pr.db.QueryRow(queries.SelectPost, id).Scan(
		&p.ID, &p.Forum, &p.Thread, &p.Author, &p.Message, &p.Parent, &p.IsEdited, &p.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return p, nil
}

func (pr *PostRepository) UpdatePost(p *models.Post) error {
	if err := pr.db.QueryRow(queries.UpdatePostWhereID, p.Message, p.ID).Scan(
		&p.Author, &p.CreatedAt, &p.ID, &p.IsEdited, &p.Message, &p.Parent, &p.Thread, &p.Forum); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}
