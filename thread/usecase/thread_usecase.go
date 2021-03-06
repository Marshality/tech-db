package usecase

import (
	"errors"
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
	"strconv"
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

func (tu *ThreadUsecase) GetByID(id uint64) (*models.Thread, error) {
	u, err := tu.threadRepo.SelectByID(id)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (tu *ThreadUsecase) Create(t *models.Thread) error {
	usr, err := tu.userUcase.GetByNickname(t.Author)

	if err != nil {
		return err
	}

	if t.Author != usr.Nickname { // на случай, если не совпадает регистр букв
		t.Author = usr.Nickname
	}

	frm, err := tu.forumUcase.GetBySlug(t.Forum)

	if err != nil {
		return err
	}

	if t.Forum != frm.Slug {
		t.Forum = frm.Slug
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

func (tu *ThreadUsecase) GetThreadsByForum(slug string, since string, limit uint64, desc bool) ([]*models.Thread, error) {
	f, err := tu.forumUcase.GetBySlug(slug)

	if err != nil {
		return nil, err
	}

	var threads []*models.Thread

	switch since {
	case "":
		threads, err = tu.threadRepo.SelectThreadsWhereForum(f.Slug, limit, desc)
	default:
		threads, err = tu.threadRepo.SelectThreadsWhereForumAndCreated(f.Slug, limit, since, desc)
	}

	if err != nil {
		return nil, err
	}

	if threads == nil {
		threads = []*models.Thread{}
	}

	return threads, nil
}

func (tu *ThreadUsecase) Vote(v *models.Vote, slugOrID string) (*models.Thread, error) {
	if _, err := tu.userUcase.GetByNickname(v.Nickname); err != nil {
		return nil, err
	}

	t, err := tu.GetThread(slugOrID)

	if err != nil {
		return nil, err
	}

	v.Thread = t.ID

	oldVote, err := tu.threadRepo.SelectVote(v.Nickname, t.ID)

	if err != nil && err == ErrNotFound {
		if err = tu.threadRepo.InsertVote(v); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	if oldVote != nil {
		if err = tu.threadRepo.UpdateVote(v); err != nil {
			return nil, err
		}

		if v.Voice == 1 && oldVote.Voice == -1 {
			v.Voice = 2
		} else if v.Voice == -1 && oldVote.Voice == 1 {
			v.Voice = -2
		} else {
			v.Voice = 0
		}
	}

	if v.Voice != 0 {
		if err := tu.threadRepo.UpdateVoteCount(v); err != nil {
			return nil, err
		}
	}

	if t, err = tu.GetByID(t.ID); err != nil {
		return nil, err
	}

	return t, nil
}

func (tu *ThreadUsecase) GetThread(slugOrID string) (*models.Thread, error) {
	id, err := strconv.Atoi(slugOrID)

	t := &models.Thread{}

	if err != nil {
		t, err = tu.GetBySlug(slugOrID)
	} else {
		t, err = tu.GetByID(uint64(id))
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tu *ThreadUsecase) EditThread(slugOrID string, t *models.Thread) error {
	founded, err := tu.GetThread(slugOrID)

	if err != nil {
		if err == ErrNotFound {
			return errors.New("user not found")
		}

		return err
	}

	t.ID = founded.ID

	if t.Message == "" {
		t.Message = founded.Message
	}

	if t.Title == "" {
		t.Title = founded.Title
	}

	err = tu.threadRepo.Update(t)

	if err != nil {
		if err == ErrNotFound {
			return errors.New("thread not found")
		}

		return err
	}

	return nil
}
