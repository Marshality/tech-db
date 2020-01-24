package post

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	CreatePosts(posts *models.Posts, slugOrID string) error
	GetPostsByThread(slugOrID string, since uint64, limit uint64, sort string, desc bool) ([]*models.Post, error)
	GetPostByID(id uint64, related ...string) (*models.Post, *models.Forum, *models.Thread, *models.User, error)
	GetPost(id uint64) (*models.Post, error)
	EditPost(p *models.Post) error
}
