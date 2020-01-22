package forum

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectBySlug(slug string) (*models.Forum, error)
	InsertThread(thread *models.Thread) error
	Insert(forum *models.Forum) error
}
