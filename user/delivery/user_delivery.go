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
	e.GET("/api/user/:nickname/profile", handler.GetUserInfo())
	e.POST("/api/user/:nickname/profile", handler.EditUserInfo())
}

func (uh *UserHandler) CreateUser() echo.HandlerFunc {
	type Request struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
	}

	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		u := &models.User{
			Nickname: nickname,
			About:    request.About,
			Fullname: request.Fullname,
			Email:    request.Email,
		}

		users, err := uh.ucase.Create(u)

		if err != ErrAlreadyExists && err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		if users != nil {
			logrus.Info("already exists")
			return c.JSON(http.StatusConflict, users)
		}

		return c.JSON(http.StatusCreated, u)
	}
}

func (uh *UserHandler) GetUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		u, err := uh.ucase.GetByNickname(nickname)

		if err != ErrNotFound && err != nil {
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		if u == nil {
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, u)
	}
}

func (uh *UserHandler) EditUserInfo() echo.HandlerFunc {
	type Request struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
	}

	return func(c echo.Context) error {
		nickname := c.Param("nickname")

		request := &Request{}
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, Error{
				Message: ErrHTTPBadRequest.Error(),
			})
		}

		u := &models.User{
			Nickname: nickname,
			Email:    request.Email,
			Fullname: request.Fullname,
			About:    request.About,
		}

		err := uh.ucase.EditUser(u)

		if err != nil && err == ErrNotFound {
			logrus.Info(err.Error())
			return c.JSON(http.StatusNotFound, Error{
				Message: err.Error(),
			})
		}

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusConflict, Error{
				Message: ErrAlreadyExists.Error(),
			})
		}

		return c.JSON(http.StatusOK, u)
	}
}
