package user

import "github.com/Marshality/tech-db/models"

type Usecase interface {
	Store(user *models.User) ([]*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	EditUser(user *models.User) error
}
