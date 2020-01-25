package queries

const (
	InsertIntoPosts = "INSERT INTO posts (id, author, forum, thread, message, parent, created_at, path) VALUES "

	SelectPostByThreadAndID = "SELECT id, forum, thread, author, message, parent, is_edited, created_at " +
		"FROM posts WHERE thread = $1 AND id = $2"

	SelectPost = "SELECT id, forum, thread, author, message, parent, is_edited, created_at " +
		"FROM posts WHERE id = $1"

	SelectPostsByThreadFlat = "SELECT id, author, forum, thread, message, parent, is_edited, created_at " +
		"FROM posts WHERE thread = $1 %s ORDER BY id %s LIMIT %s"

	SelectPostsByThreadTree = "SELECT posts.id, posts.author, posts.forum, posts.thread, " +
		"posts.message, posts.parent, posts.is_edited, posts.created_at " +
		"FROM posts %s posts.thread = $1 ORDER BY posts.path[1] %s, posts.path %s LIMIT %s"

	SelectPostsByThreadParentTree = "WITH roots AS (" +
		"SELECT id, author, forum, thread, message, parent, is_edited, created_at, path, " +
		"dense_rank() OVER (ORDER BY path[1] %s) as root " +
		"FROM posts WHERE thread = $1" +
		") SELECT roots.id, roots.author, roots.forum, roots.thread, roots.message, " +
		"roots.parent, roots.is_edited, roots.created_at FROM roots %s"

	SelectParentTree = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.is_edited, p.created_at " +
		"FROM posts as p WHERE p.thread = $1 AND " +
		"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 %s %s %s"

	SelectPostWhereID = "SELECT P.author AS post_author, P.created_at AS post_created, " +
		"P.id AS post_id, P.is_edited AS post_is_edited, P.message AS post_message, " +
		"P.parent AS post_parent, P.thread AS post_thread, T.forum AS post_forum %s " +
		"FROM posts P JOIN threads T ON P.thread = T.id %s " +
		"WHERE P.id = $1"

	UpdatePostWhereID = "UPDATE posts SET message = COALESCE($1, posts.message), " +
		"is_edited = COALESCE($1, posts.message) <> posts.message " +
		"FROM threads WHERE threads.id = posts.thread AND posts.id = $2 " +
		"RETURNING posts.author, posts.created_at, posts.id, posts.is_edited, " +
		"posts.message, posts.parent, posts.thread, threads.forum"
)
