package usecase

import (
	"errors"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	"github.com/Marshality/tech-db/thread"
	"strconv"
	"time"
)

type PostUsecase struct {
	postRepo post.Repository
	threadUcase thread.Usecase
}

func NewPostUsecase(pr post.Repository, tUc thread.Usecase) post.Usecase {
	return &PostUsecase{
		postRepo: pr,
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
			if _, err := pu.postRepo.SelectByThreadAndID(threadID, p.Parent); err != nil {
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
