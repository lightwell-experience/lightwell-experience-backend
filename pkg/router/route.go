package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/handler"
)

func ConfigureEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	handler.RegisterPing(e)

	return e
}
