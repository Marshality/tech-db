package thread

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	GetBySlug(slug string) (*models.Thread, error)
	GetByID(id uint64) (*models.Thread, error)
	Create(t *models.Thread) error
	GetThreadsByForum(slug string, since string, limit uint64, desc bool) ([]*models.Thread, error)
	Vote(v *models.Vote, slugOrID string) (*models.Thread, error)
	GetThread(slugOrID string) (*models.Thread, error)
}
