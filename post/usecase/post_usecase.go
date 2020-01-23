package usecase

import (
	"errors"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"strconv"
	"time"
)

type PostUsecase struct {
	postRepo    post.Repository
	threadUcase thread.Usecase
}

func NewPostUsecase(pr post.Repository, tUc thread.Usecase) post.Usecase {
	return &PostUsecase{
		postRepo:    pr,
		threadUcase: tUc,
	}
}

func (pu *PostUsecase) CreatePosts(posts *models.Posts, slugOrID string) error {
	id, err := strconv.Atoi(slugOrID)

	t := &models.Thread{}

	if err != nil {
		t, err = pu.threadUcase.GetBySlug(slugOrID)
	} else {
		t, err = pu.threadUcase.GetByID(uint64(id))
	}

	if err != nil {
		return err
	}

	forumID := t.Forum
	threadID := t.ID
	createdAt := time.Now().Format(time.RFC3339)

	for _, p := range *posts {
		if p.Parent != 0 {
			if _, err := pu.postRepo.SelectByThreadAndID(p.Parent, threadID); err != nil {
				return errors.New("conflict")
			}
		}

		p.Forum = forumID
		p.Thread = threadID
		p.CreatedAt = createdAt

		if err := pu.postRepo.Insert(p); err != nil {
			return err
		}
	}

	return nil
}

func (pu *PostUsecase) GetPostsByThread(slugOrID string, since *uint64, limit uint64, sort string, desc bool) ([]*models.Post, error) {
	id, err := strconv.Atoi(slugOrID)

	t := &models.Thread{}

	if err != nil {
		t, err = pu.threadUcase.GetBySlug(slugOrID)
	} else {
		t, err = pu.threadUcase.GetByID(uint64(id))
	}

	if err != nil {
		return nil, err
	}

	threadID := t.ID

	switch sort {
	case FLAT_SORT, "":
		return pu.postRepo.GetPostsByThreadFlat(threadID, since, limit, desc)
	case TREE_SORT:
		return pu.postRepo.GetPostsByThreadTree(threadID, since, limit, desc)
	case PARENT_TREE_SORT:
		return pu.postRepo.GetPostsByThreadParentTree(threadID, since, limit, desc)
	}

	return nil, ErrNotFound
}
