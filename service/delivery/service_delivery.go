package delivery

import (
	"github.com/Marshality/tech-db/service"
	. "github.com/Marshality/tech-db/tools"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ServiceHandler struct {
	serviceUcase service.Usecase
}

func ConfigureServiceHandler(e *echo.Echo, sUc service.Usecase) {
	handler := &ServiceHandler{
		serviceUcase: sUc,
	}

	e.GET("/api/service/status", handler.Status())
	e.POST("/api/service/clear", handler.Clear())
}

func (sh *ServiceHandler) Clear() echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := sh.serviceUcase.DoClear(); err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		return c.String(http.StatusOK, "Success")
	}
}

func (sh *ServiceHandler) Status() echo.HandlerFunc {
	type Response struct {
		Forum uint64 `json:"forum"`
		Thread uint64 `json:"thread"`
		User uint64 `json:"user"`
		Post uint64 `json:"post"`
	}

	return func(c echo.Context) error {
		f, t, p, u, err := sh.serviceUcase.GetStatus()

		if err != nil {
			logrus.Info(err.Error())
			return c.JSON(http.StatusInternalServerError, Error{
				Message: err.Error(),
			})
		}

		response := &Response{
			Forum:  f,
			Thread: t,
			User:   u,
			Post:   p,
		}

		return c.JSON(http.StatusOK, response)
	}
}
