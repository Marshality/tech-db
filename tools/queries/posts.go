package queries

const (
	InsertIntoPosts = "INSERT INTO posts (author, forum, thread, message, parent, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	SelectPostByThreadAndID = "SELECT * FROM posts WHERE thread = $1 AND id = $2"

	SelectPostsByThreadFlat = "SELECT id, author, forum, thread, message, parent, created_at " +
		"FROM posts WHERE thread = $1 %s ORDER BY id %s LIMIT %s"
)
