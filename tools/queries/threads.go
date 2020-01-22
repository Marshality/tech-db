package queries

const (
	InsertIntoThreads = "INSERT INTO threads (slug, author, forum, message, title, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	SelectThreadWhereSlug = "SELECT id, slug, author, forum, message, title, votes, created_at FROM threads WHERE slug = $1"
)
