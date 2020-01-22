package user

import "github.com/Marshality/tech-db/models"

type Repository interface {
	SelectWhere(nickname, email string) ([]*models.User, error)
	SelectByNickname(nickname string) (*models.User, error)
	SelectByEmail(email string) (*models.User, error)
	Insert(user *models.User) error
	Update(user *models.User) error
}
