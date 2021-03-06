package repository

import (
	"database/sql"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/tools/queries"
	"log"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(conn *sql.DB) thread.Repository {
	return &ThreadRepository{
		db: conn,
	}
}

func (tr *ThreadRepository) UpdateVote(v *models.Vote) error {
	_, err := tr.db.Exec(queries.UpdateVote, v.Voice, v.Nickname, v.Thread)
	return err
}

func (tr *ThreadRepository) SelectVote(nickname string, thread uint64) (*models.Vote, error) {
	v := &models.Vote{}

	if err := tr.db.QueryRow(queries.SelectVote, nickname, thread).Scan(&v.ID, &v.Nickname, &v.Thread, &v.Voice); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return v, nil
}

func (tr *ThreadRepository) UpdateVoteCount(v *models.Vote) error {
	_, err := tr.db.Exec(queries.UpdateThreadVotes, v.Voice, v.Thread)
	return err
}

func (tr *ThreadRepository) InsertVote(v *models.Vote) error {
	if err := tr.db.QueryRow(queries.InsertVote, v.Thread, v.Nickname, v.Voice).Scan(&v.ID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
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

func (tr *ThreadRepository) SelectByID(id uint64) (*models.Thread, error) {
	t := &models.Thread{}

	if err := tr.db.QueryRow(queries.SelectThreadWhereID, id).Scan(
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

func (tr *ThreadRepository) SelectThreadsWhereForum(slug string, limit uint64, desc bool) ([]*models.Thread, error) {
	var threads []*models.Thread

	rows, err := tr.db.Query(queries.QM.SelectThreadsWhereForum(desc), slug, limit)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		t := &models.Thread{}

		if err := rows.Scan(&t.ID, &t.Slug, &t.Author, &t.Forum, &t.Message, &t.Title, &t.Votes, &t.CreatedAt); err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	return threads, nil
}

func (tr *ThreadRepository) SelectThreadsWhereForumAndCreated(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	var threads []*models.Thread

	rows, err := tr.db.Query(queries.QM.SelectThreadsWhereForumAndCreated(desc), slug, since, limit)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for rows.Next() {
		t := &models.Thread{}

		if err := rows.Scan(&t.ID, &t.Slug, &t.Author, &t.Forum, &t.Message, &t.Title, &t.Votes, &t.CreatedAt); err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	return threads, nil
}

func (tr *ThreadRepository) Update(t *models.Thread) error {
	if err := tr.db.QueryRow(queries.UpdateThreadsWhere, t.Message, t.Title, t.ID).Scan(
		&t.Author, &t.Forum, &t.Slug, &t.Votes, &t.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				return ErrNotFound
			}

			return err
	}

	return nil
}
