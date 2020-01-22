package queries

const (
	InsertIntoThreads = "INSERT INTO threads (slug, author, message, title, created) VALUES ($1, $2, $3, $4, $5) RETURNING id"
)
