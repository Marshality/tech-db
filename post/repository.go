package post

import "github.com/Marshality/tech-db/models"

type Repository interface {
	Insert(post *models.Post) error
	SelectByThreadAndID(id, thread uint64) (*models.Post, error)
	GetPostsByThread(thread uint64, since uint64, limit uint64, desc bool, sort string) ([]*models.Post, error)
}
