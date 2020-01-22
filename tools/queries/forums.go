package queries

const (
	SelectForumWhereSlug = "SELECT id, slug, posts, threads, title, usr FROM forums WHERE slug = $1"
	InsertIntoForums     = "INSERT INTO forums (slug, title, usr) VALUES ($1, $2, $3) RETURNING id"
)
