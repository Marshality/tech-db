package queries

const (
	StatusQuery = "SELECT " +
		"(SELECT COUNT(*) FROM forums) as forums_status, " +
		"(SELECT COUNT(*) FROM threads) as threads_status, " +
		"(SELECT COUNT(*) FROM posts) as posts_status, " +
		"(SELECT COUNT(*) FROM users) as users_status"

	ClearQuery = "TRUNCATE users, forums, user_forum, vote, posts, threads RESTART IDENTITY"
)
