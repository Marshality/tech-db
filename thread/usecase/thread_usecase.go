package usecase

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
	"time"
)

type ThreadUsecase struct {
	threadRepo thread.Repository
	forumUcase forum.Usecase
	userUcase  user.Usecase
}

func NewThreadUsecase(tr thread.Repository, uUc user.Usecase, fUc forum.Usecase) thread.Usecase {
	return &ThreadUsecase{
		threadRepo: tr,
		forumUcase: fUc,
		userUcase:  uUc,
	}
}

func (tu *ThreadUsecase) GetBySlug(slug string) (*models.Thread, error) {
	u, err := tu.threadRepo.SelectBySlug(slug)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (tu *ThreadUsecase) Create(t *models.Thread) error {
	frm, err := tu.forumUcase.GetBySlug(t.Forum)

	if err != nil {
		return err
	}

	if t.Forum != frm.Slug {
		t.Forum = frm.Slug
	}

	usr, err := tu.userUcase.GetByNickname(t.Author)

	if err != nil {
		return err
	}

	if t.Author != usr.Nickname { // на случай, если не совпадает регистр букв
		t.Author = usr.Nickname
	}

	if t.Slug != "" {
		thr, err := tu.GetBySlug(t.Slug)

		if err != nil && err != ErrNotFound {
			return err
		}

		if thr != nil {
			*t = *thr
			return ErrAlreadyExists
		}
	}

	if t.CreatedAt == "" {
		t.CreatedAt = time.Now().Format(time.RFC3339)
	}

	if err := tu.threadRepo.Insert(t); err != nil {
		return err
	}

	return nil
}
