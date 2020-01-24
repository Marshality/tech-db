package forum

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	Insert(forum *models.Forum) error
	SelectForumUsers(slug string, limit uint64, since string, desc bool) ([]*models.User, error)
}
