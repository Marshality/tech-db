package user

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	Store(user *models.User) ([]*models.User, error)
}
