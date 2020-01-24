package queries

const (
	InsertIntoPosts = "INSERT INTO posts (author, forum, thread, message, parent, created_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	SelectPostByThreadAndID = "SELECT id, forum, thread, author, message, parent, is_edited, created_at " +
		"FROM posts WHERE thread = $1 AND id = $2"

	SelectPostsByThreadFlat = "SELECT id, author, forum, thread, message, parent, is_edited, created_at " +
		"FROM posts WHERE thread = $1 %s ORDER BY id %s LIMIT %s"

	SelectPostsByThreadTree = "SELECT posts.id, posts.author, posts.forum, posts.thread, " +
		"posts.message, posts.parent, posts.is_edited, posts.created_at " +
		"FROM posts %s posts.thread = $1 ORDER BY posts.path %s LIMIT %s"

	SelectPostsByThreadParentTree = "WITH roots AS (" +
		"SELECT id, author, forum, thread, message, parent, is_edited, created_at, path, " +
		"dense_rank() OVER (ORDER BY path[1] %s) as root " +
		"FROM posts WHERE thread = $1" +
		") SELECT roots.id, roots.author, roots.forum, roots.thread, roots.message, " +
		"roots.parent, roots.is_edited, roots.created_at FROM roots %s"
)
