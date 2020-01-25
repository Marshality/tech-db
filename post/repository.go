package post

import "github.com/Marshality/tech-db/models"

type Repository interface {
	Insert(posts *models.Posts, forum string, threadID uint64, createdAt string) error
	SelectByThreadAndID(id, thread uint64) (*models.Post, error)
	GetPostsByThread(thread uint64, since uint64, limit uint64, desc bool, sort string) ([]*models.Post, error)
	SelectPostWhereID(id uint64, related ...string) (*models.Post, *models.Thread, *models.Forum, *models.User, error)
	SelectPost(id uint64) (*models.Post, error)
	UpdatePost(p *models.Post) error
}
