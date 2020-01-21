package forum

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	Create(forum *models.Forum) error
}
