package post

import "github.com/Marshality/tech-db/models"

type Repository interface {
	Insert(post *models.Post) error
	SelectByThreadAndID(id, thread uint64) (*models.Post, error)
}
