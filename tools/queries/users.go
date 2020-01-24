package queries

const (
	SelectUsersWhereNicknameOrEmail = "SELECT id, nickname, fullname, email, about FROM users WHERE nickname = $1 OR email = $2"
	InsertIntoUsers                 = "INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING id"
	SelectUserWhereNickname         = "SELECT id, nickname, fullname, email, about FROM users WHERE nickname = $1"
	SelectUserWhereEmail            = "SELECT id, nickname, fullname, email, about FROM users WHERE email = $1"
	UpdateUsers                     = "UPDATE users SET about = $1, fullname = $2, email = $3 WHERE nickname = $4"

	SelectUsersWhereForumSlug = "SELECT U.id, U.nickname, U.fullname, U.email, U.about " +
		"FROM user_forum UF JOIN users U ON UF.user_id = U.id " +
		"WHERE UF.forum_slug = $1 %s ORDER BY U.nickname %s"
)
