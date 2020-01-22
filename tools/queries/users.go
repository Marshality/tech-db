package queries

const (
	SelectUsersWhereNicknameOrEmail = "SELECT id, nickname, fullname, email, about FROM users WHERE nickname = $1 OR email = $2"
	InsertIntoUsers                 = "INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING id"
	SelectUserWhereNickname         = "SELECT id, nickname, fullname, email, about FROM users WHERE nickname = $1"
	SelectUserWhereEmail            = "SELECT id, nickname, fullname, email, about FROM users WHERE email = $1"
	UpdateUsers                     = "UPDATE users SET about = $1, fullname = $2, email = $3 WHERE nickname = $4"
)
