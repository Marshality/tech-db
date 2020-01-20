package user

import "github.com/Marshality/tech-db/models"

type Repository interface {
	GetUsersWhere(nickname, email string) ([]*models.User, error)
	Create(user *models.User) error
}
