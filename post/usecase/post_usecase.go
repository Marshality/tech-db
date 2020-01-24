package usecase

import (
	"errors"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	"github.com/Marshality/tech-db/thread"
	"github.com/Marshality/tech-db/user"
	"github.com/sirupsen/logrus"
	"time"
)

type PostUsecase struct {
	postRepo    post.Repository
	threadUcase thread.Usecase
	userUcase user.Usecase
}

func NewPostUsecase(pr post.Repository, tUc thread.Usecase) post.Usecase {
	return &PostUsecase{
		postRepo:    pr,
		threadUcase: tUc,
	}
}

func (pu *PostUsecase) CreatePosts(posts *models.Posts, slugOrID string) error {
	t, err := pu.threadUcase.GetThread(slugOrID)

	if err != nil {
		return err
	}

	forumID := t.Forum
	threadID := t.ID
	createdAt := time.Now().Format(time.RFC3339)

	for _, p := range *posts {
		if p.Parent != 0 {
			if _, err := pu.postRepo.SelectByThreadAndID(p.Parent, threadID); err != nil {
				logrus.Info(err)
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

func (pu *PostUsecase) GetPostsByThread(slugOrID string, since uint64, limit uint64, sort string, desc bool) ([]*models.Post, error) {
	t, err := pu.threadUcase.GetThread(slugOrID)

	if err != nil {
		return nil, err
	}

	threadID := t.ID

	posts, err := pu.postRepo.GetPostsByThread(threadID, since, limit, desc, sort)

	if posts == nil {
		posts = []*models.Post{}
	}

	return posts, err
}

func (pu *PostUsecase) GetPost(id uint64) (*models.Post, error) {
	return pu.postRepo.SelectPost(id)
}

func (pu *PostUsecase) GetPostByID(id uint64, related ...string) (*models.Post, *models.Forum, *models.Thread, *models.User, error) {
	if _, err := pu.GetPost(id); err != nil {
		return nil, nil, nil, nil, err
	}

	p, t, f, u, err := pu.postRepo.SelectPostWhereID(id, related...)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	return p, f, t, u, err
}

func (pu *PostUsecase) EditPost(p *models.Post) error {
	founded, err := pu.postRepo.SelectPost(p.ID)

	if err != nil {
		return err
	}

	if p.Message == "" {
		p.Message = founded.Message
	}

	return pu.postRepo.UpdatePost(p)
}
