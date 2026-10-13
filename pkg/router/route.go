package router

import (
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/handler"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/instrumentation"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/middleware"
)

func ConfigureEcho() *echo.Echo {
	e := echo.New()
	e.Use(echo_middleware.Recover())
	e.Use(echo_middleware.RequestLogger())

	handler.RegisterPing(e)

	return e
}

func ConfigureEchoWithMetrics(metrics *instrumentation.Metrics) *echo.Echo {
	e := ConfigureEcho()
	e.Use(middleware.MetricsMiddleware(metrics))
	return e
}
