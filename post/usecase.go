package post

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	CreatePosts(posts *models.Posts, slugOrID string) error
}
