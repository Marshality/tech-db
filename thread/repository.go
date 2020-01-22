package thread

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectBySlug(slug string) (*models.Thread, error)
	Insert(t *models.Thread) error
}
