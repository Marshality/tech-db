package forum

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	Store(forum *models.Forum) error
}
