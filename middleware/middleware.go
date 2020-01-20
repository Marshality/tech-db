package middleware

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"time"
)

type GoMiddleware struct {}

func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func (m *GoMiddleware) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Println(c.Request().Method, c.Request().URL)

		start := time.Now()
		err := next(c)
		end := time.Now()

		logrus.Println(c.Response().Status, end.Sub(start))
		fmt.Println("###")

		return err
	}
}

func NewGoMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
