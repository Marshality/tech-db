package delivery

import (
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	postUcase post.Usecase
}

func ConfigurePostHandler(e *echo.Echo, pUc post.Usecase) {
	handler := &PostHandler{
		postUcase: pUc,
	}

	e.GET("/api/post/:id/details", handler.GetPostDetails())
	e.POST("/api/post/:id/details", handler.EditPost())
}

func (ph *PostHandler) GetPostDetails() echo.HandlerFunc {
	type Response struct {
		Post   *models.Post   `json:"post,omitempty"`
		Thread *models.Thread `json:"thread,omitempty"`
		Forum  *models.Forum  `json:"forum,omitempty"`
		User   *models.User   `json:"author,omitempty"`
	}

	return func(c echo.Context) error {
		postID := c.Param("id")

		id, err := strconv.Atoi(postID)

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusBadRequest, Error{
				Message: "bad request",
			})
		}

		param := c.QueryParam("related")

		related := strings.Split(param, ",")

		p, f, t, u, err := ph.postUcase.GetPostByID(uint64(id), related...)

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

		response := &Response{
			Post:   p,
			Thread: t,
			Forum:  f,
			User:   u,
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (ph *PostHandler) EditPost() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}

	return func(c echo.Context) error {
		postID := c.Param("id")

		id, err := strconv.Atoi(postID)

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusBadRequest, Error{
				Message: "bad request",
			})
		}

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		p := &models.Post{
			ID: uint64(id),
			Message: request.Message,
		}

		err = ph.postUcase.EditPost(p)

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

		return c.JSON(http.StatusOK, p)
	}
}
