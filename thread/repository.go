package thread

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectBySlug(slug string) (*models.Thread, error)
	SelectByID(id uint64) (*models.Thread, error)
	Insert(t *models.Thread) error
	SelectThreadsWhereForum(slug string, limit uint64, desc bool) ([]*models.Thread, error)
	SelectThreadsWhereForumAndCreated(slug string, limit uint64, since string, desc bool) ([]*models.Thread, error)
	InsertVote(v *models.Vote) error
	UpdateVote(v *models.Vote) error
	SelectVote(nickname string, thread uint64) (*models.Vote, error)
	UpdateVoteCount(v *models.Vote) error
	Update(t *models.Thread) error
}
