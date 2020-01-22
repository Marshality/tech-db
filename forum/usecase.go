package forum

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	GetBySlug(slug string) (*models.Forum, error)
	Create(forum *models.Forum) error
}
