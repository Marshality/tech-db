package queries

const (
	InsertIntoThreads = "INSERT INTO threads (slug, author, forum, message, title, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	UpdateThreadsWhere = "UPDATE threads SET message = $1, title = $2 WHERE id = $3 " +
		"RETURNING author, forum, slug, votes, created_at"

	SelectThreadWhereSlug = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads WHERE slug = $1"

	SelectThreadWhereID = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads WHERE id = $1"

	SelectThreadsWhereForum = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads " +
		"WHERE forum = $1 ORDER BY created_at %s LIMIT $2"

	SelectThreadsWhereForumAndCreated = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads " +
		"WHERE forum = $1 AND created_at %s $2 ORDER BY created_at %s LIMIT $3"

	UpdateThreadVotes = "UPDATE threads SET votes = votes + $1 WHERE id = $2"
)
