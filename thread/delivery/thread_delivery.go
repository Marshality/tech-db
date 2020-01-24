package delivery

import (
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type ThreadHandler struct {
	threadUcase thread.Usecase
	postUcase   post.Usecase
}

func ConfigureThreadHandler(e *echo.Echo, tUc thread.Usecase, pUc post.Usecase) {
	handler := &ThreadHandler{
		threadUcase: tUc,
		postUcase:   pUc,
	}

	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts())
	e.POST("/api/thread/:slug_or_id/vote", handler.Vote())
	e.GET("/api/thread/:slug_or_id/details", handler.GetThreadDetails())
	e.GET("/api/thread/:slug_or_id/posts", handler.GetThreadPosts())
	e.POST("/api/thread/:slug_od_id/details", handler.EditThread())
}

func (th *ThreadHandler) CreatePosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		var posts models.Posts
		if err := c.Bind(&posts); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		err := th.postUcase.CreatePosts(&posts, slugOrID)

		if err != nil && err == ErrNotFound {
			logrus.Info(err.Error())
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		if err != nil && err.Error() == "conflict" {
			logrus.Info(err.Error())
			return c.JSON(http.StatusConflict, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, posts)
	}
}

func (th *ThreadHandler) Vote() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Voice    int8   `json:"voice"`
	}

	voiceIsValid := func(value int8) bool {
		return value == 1 || value == -1
	}

	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		request := &Request{}
		if err := c.Bind(request); err != nil || !voiceIsValid(request.Voice) {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		v := &models.Vote{
			Nickname: request.Nickname,
			Voice:    request.Voice,
		}

		t, err := th.threadUcase.Vote(v, slugOrID)

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

		return c.JSON(http.StatusOK, t)
	}
}

func (th *ThreadHandler) GetThreadDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		t, err := th.threadUcase.GetThread(slugOrID)

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

		return c.JSON(http.StatusOK, t)
	}
}

func (th *ThreadHandler) GetThreadPosts() echo.HandlerFunc {
	type Request struct {
		Limit uint64 `query:"limit"`
		Since uint64 `query:"since"`
		Sort  string `query:"sort"`
		Desc  bool   `query:"desc"`
	}

	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		posts, err := th.postUcase.GetPostsByThread(slugOrID, request.Since, request.Limit, request.Sort, request.Desc)

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

		return c.JSON(http.StatusOK, posts)
	}
}

func (th *ThreadHandler) EditThread() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
		Title   string `json:"title"`
	}

	return func(c echo.Context) error {
		slugOrID := c.Param("slug_or_id")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		t := &models.Thread{
			Message: request.Message,
			Title: request.Title,
		}

		err := th.threadUcase.EditThread(slugOrID, t)

		if err != nil && strings.Contains(err.Error(), "not found") {
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

		return c.JSON(http.StatusOK, t)
	}
}
