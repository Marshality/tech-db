package queries

import (
	"fmt"
)

const (
	InsertIntoThreads = "INSERT INTO threads (slug, author, forum, message, title, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	SelectThreadWhereSlug = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads WHERE slug = $1"

	SelectThreadsWhereForum = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads " +
		"WHERE forum = $1 ORDER BY created_at %s LIMIT $2"

	SelectThreadsWhereForumAndCreated = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads " +
		"WHERE forum = $1 AND created_at %s $2 ORDER BY created_at %s LIMIT $3"
)

type QueryManager struct{}
var QM *QueryManager

func (qm *QueryManager) SelectThreadsWhereForum(desc bool) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectThreadsWhereForum, "DESC")
	case false:
		query = fmt.Sprintf(SelectThreadsWhereForum, "ASC")
	}

	return query
}

func (qm *QueryManager) SelectThreadsWhereForumAndCreated(desc bool) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectThreadsWhereForumAndCreated, "<=", "DESC")
	case false:
		query = fmt.Sprintf(SelectThreadsWhereForumAndCreated, ">=", "ASC")
	}

	return query
}
