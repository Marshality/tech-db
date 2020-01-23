package queries

const (
	InsertVote = "INSERT INTO vote (thread, nickname, voice) VALUES ($1, $2, $3) RETURNING id"
	UpdateVote = "UPDATE vote SET voice = $1 WHERE nickname = $2 AND thread = $3"
)
