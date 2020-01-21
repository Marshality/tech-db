package usecase

import (
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
)

type UserUsecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) user.Usecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (ur *UserUsecase) Store(user *models.User) ([]*models.User, error) {
	users, err := ur.repo.SelectWhere(user.Nickname, user.Email)

	if err != nil {
		return nil, err
	}

	if users != nil {
		return users, ErrAlreadyExists
	}

	if err := ur.repo.Create(user); err != nil {
		return nil, err
	}

	return nil, nil
}

func (ur *UserUsecase) GetByNickname(nickname string) (*models.User, error) {
	u, err := ur.repo.SelectByNickname(nickname)

	if err == ErrNotFound {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserUsecase) EditUser(u *models.User) error {
	err := ur.repo.Update(u)

	if err != nil {
		return err
	}

	return nil
}
