package usecase

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
)

type ForumUsecase struct {
	forumRepo   forum.Repository
	userUcase   user.Usecase
}

func NewForumUsecase(fr forum.Repository, uUc user.Usecase) forum.Usecase {
	return &ForumUsecase{
		forumRepo:   fr,
		userUcase:   uUc,
	}
}

func (fu *ForumUsecase) GetBySlug(slug string) (*models.Forum, error) {
	u, err := fu.forumRepo.SelectBySlug(slug)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (fu *ForumUsecase) Create(f *models.Forum) error {
	frm, err := fu.GetBySlug(f.Slug)

	if err != nil && err != ErrNotFound {
		return err
	}

	if frm != nil {
		*f = *frm
		return ErrAlreadyExists
	}

	usr, err := fu.userUcase.GetByNickname(f.User)

	if err != nil {
		return err
	}

	if f.User != usr.Nickname { // на случай, если не совпадает регистр букв
		f.User = usr.Nickname
	}

	if err := fu.forumRepo.Insert(f); err != nil {
		return err
	}

	return nil
}
