package queries

const (
	InsertVote = "INSERT INTO vote (thread, nickname, voice) VALUES ($1, $2, $3) RETURNING id"

	//VoteQuery = "INSERT INTO vote (nickname, thread, voice) " +
	//	"VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT unique_vote " +
	//	"DO UPDATE SET voice = $3 WHERE vote.thread = $2 AND vote.nickname = $1 RETURNING vote.id"

	UpdateVote = "UPDATE vote SET voice = $1 WHERE nickname = $2 AND thread = $3"

	SelectVote = "SELECT id, nickname, thread, voice FROM vote WHERE nickname = $1 AND thread = $2"
)
