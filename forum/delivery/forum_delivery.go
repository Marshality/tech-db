package delivery

import (
	"github.com/Marshality/tech-db/forum"
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ForumHandler struct {
	ucase forum.Usecase
}

func ConfigureForumHandler(e *echo.Echo, ucase forum.Usecase) {
	handler := &ForumHandler{
		ucase: ucase,
	}

	e.POST("/api/forum/create", handler.CreateForum())
	e.GET("/api/forum/:slug/details", handler.GetForumDetails())
	e.POST("/api/forum/:slug/create", handler.CreateThread())
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

		err := fh.ucase.Create(f)

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

		f, err := fh.ucase.GetBySlug(slug)

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
	return func(c echo.Context) error {
		return nil
	}
}
