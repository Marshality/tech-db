package delivery

import (
	"github.com/Marshality/tech-db/models"
	"github.com/Marshality/tech-db/post"
	"github.com/Marshality/tech-db/thread"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ThreadHandler struct {
	threadUcase thread.Usecase
	postUcase post.Usecase
}

func ConfigureThreadHandler(e *echo.Echo, tUc thread.Usecase, pUc post.Usecase) {
	handler := &ThreadHandler{
		threadUcase: tUc,
		postUcase: pUc,
	}

	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts())
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
