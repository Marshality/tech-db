package delivery

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ForumHandler struct {
	forumUcase  forum.Usecase
	threadUcase thread.Usecase
}

func ConfigureForumHandler(e *echo.Echo, fUc forum.Usecase, tUc thread.Usecase) {
	handler := &ForumHandler{
		forumUcase:  fUc,
		threadUcase: tUc,
	}

	e.POST("/api/forum/create", handler.CreateForum())
	e.GET("/api/forum/:slug/details", handler.GetForumDetails())
	e.POST("/api/forum/:slug/create", handler.CreateThread())
	e.GET("/api/forum/:slug/threads", handler.GetForumThreads())
}

func (fh *ForumHandler) CreateForum() echo.HandlerFunc {
	type Request struct {
		Slug  string `json:"slug"`
		User  string `json:"user"`
		Title string `json:"title"`
	}

	return func(c echo.Context) error {
		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		f := &models.Forum{
			Slug:  request.Slug,
			User:  request.User,
			Title: request.Title,
		}

		err := fh.forumUcase.Create(f)

		if err != nil && err == ErrNotFound {
			logrus.Info(err.Error())
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		if err != nil && err == ErrAlreadyExists {
			logrus.Info(err.Error())
			return c.JSON(http.StatusConflict, f)
		}

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, f)
	}
}

func (fh *ForumHandler) GetForumDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")

		f, err := fh.forumUcase.GetBySlug(slug)

		if err != ErrNotFound && err != nil {
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		if f == nil {
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, f)
	}
}

func (fh *ForumHandler) CreateThread() echo.HandlerFunc {
	type Request struct {
		Author  string `json:"author"`
		Created string `json:"created"`
		Title   string `json:"title"`
		Slug    string `json:"slug"`
		Message string `json:"message"`
	}

	return func(c echo.Context) error {
		forumSlug := c.Param("slug")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		t := &models.Thread{
			Forum:     forumSlug,
			Slug:      request.Slug,
			Author:    request.Author,
			Title:     request.Title,
			CreatedAt: request.Created,
			Message:   request.Message,
		}

		err := fh.threadUcase.Create(t)

		if err != nil && err == ErrNotFound {
			logrus.Info(err.Error())
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		if err != nil && err == ErrAlreadyExists {
			logrus.Info(err.Error())
			return c.JSON(http.StatusConflict, t)
		}

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, t)
	}
}

func (fh *ForumHandler) GetForumThreads() echo.HandlerFunc {
	type Request struct {
		Limit uint64 `json:"limit"`
		Desc  bool   `json:"desc"`
		Since string `json:"since"`
	}

	return func(c echo.Context) error {
		forumSlug := c.Param("slug")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		threads, err := fh.threadUcase.GetThreadsByForum(forumSlug, request.Since, request.Limit, request.Desc)

		if err != nil && err == ErrNotFound {
			logrus.Info(err.Error())
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, threads)
	}
}
