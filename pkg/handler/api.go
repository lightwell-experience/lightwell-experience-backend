package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterPing(engine *echo.Echo) {
	engine.GET("/ping", ping)
	engine.GET("/ping/", ping)
	engine.GET("/api/lightwell-next/v1.0/status", ping)
}

func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "pong",
	})
}
