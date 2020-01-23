package usecase

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
	"github.com/lib/pq"
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

	v.Thread = t.ID

	if err = tu.threadRepo.InsertVote(v); err != nil {
		psqlError, ok := err.(*pq.Error)

		if !ok {
			return nil, err
		}

		if psqlError.Code != "23505" {
			return nil, err
		}

		if err := tu.threadRepo.UpdateVote(v); err != nil {
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
