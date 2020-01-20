package delivery

import (
	"github.com/Marshality/tech-db/models"
	. "github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	ucase user.Usecase
}

func ConfigureUserHandler(e *echo.Echo, ucase user.Usecase) {
	handler := &UserHandler{
		ucase: ucase,
	}

	e.POST("/api/user/:nickname/create", handler.CreateUser())
	//TODO: e.GET("/api/user/:nickname/profile", handler.GetUserInfo())
}

func (uh *UserHandler) CreateUser() echo.HandlerFunc {
	type Request struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
	}

	type Body map[string]interface{}

	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Body{
				"error": "bad request",
			})
		}

		u := &models.User{
			Nickname: nickname,
			About:    request.About,
			Fullname: request.Fullname,
			Email:    request.Email,
		}

		users, err := uh.ucase.Store(u)

		if err != ErrAlreadyExists && err != nil {
			return c.JSON(http.StatusInternalServerError, Body{
				"error": err.Error(),
			})
		}

		if users != nil {
			logrus.Info("already exists")
			return c.JSON(http.StatusConflict, users)
		}

		return c.JSON(http.StatusCreated, u)
	}
}
