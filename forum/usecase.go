package forum

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	GetBySlug(slug string) (*models.Forum, error)
	Create(forum *models.Forum) error
	GetForumUsers(slug string, limit uint64, since string, desc bool) ([]*models.User, error)
}
