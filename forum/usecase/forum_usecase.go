package usecase

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
)

type ForumUsecase struct {
	forumRepo forum.Repository
	userUcase user.Usecase
}

func NewForumUsecase(fr forum.Repository, uUc user.Usecase) forum.Usecase {
	return &ForumUsecase{
		forumRepo: fr,
		userUcase: uUc,
	}
}

func (fu *ForumUsecase) Store(f *models.Forum) error {
	if _, err := fu.userUcase.GetByNickname(f.User); err != nil {
		return err
	}

	frm, err := fu.forumRepo.SelectBySlug(f.Slug)

	if err != nil && err != ErrNotFound {
		return err
	}

	if frm != nil {
		*f = *frm
		return ErrAlreadyExists
	}

	if err := fu.forumRepo.Create(f); err != nil {
		return err
	}

	return nil
}
