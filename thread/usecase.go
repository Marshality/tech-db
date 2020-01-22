package thread

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	GetBySlug(slug string) (*models.Thread, error)
	Create(t *models.Thread) error
}
