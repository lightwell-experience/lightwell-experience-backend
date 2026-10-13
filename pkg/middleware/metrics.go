package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/config"
	handler_utils "github.com/lightwell-experience/lightwell-experience-backend/pkg/handler/utils"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/instrumentation"
)

type MetricsConfig struct {
	Skipper echo_middleware.Skipper
	Metrics *instrumentation.Metrics
}

func MetricsMiddleware(metrics *instrumentation.Metrics) echo.MiddlewareFunc {
	return MetricsMiddlewareWithConfig(MetricsConfig{Metrics: metrics})
}

func MetricsMiddlewareWithConfig(config MetricsConfig) echo.MiddlewareFunc {
	if config.Metrics == nil {
		panic("config.Metrics cannot be nil")
	}
	if config.Skipper == nil {
		config.Skipper = metricsSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			status := mapStatus(statusCode(c, err))
			config.Metrics.HTTPStatusHistogram.WithLabelValues(
				status,
				c.Request().Method,
				matchedPath(c),
			).Observe(time.Since(start).Seconds())
			return err
		}
	}
}

func mapStatus(status int) string {
	switch {
	case status >= 100 && status < 200:
		return "1xx"
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500 && status < 600:
		return "5xx"
	default:
		return ""
	}
}

func statusCode(c echo.Context, err error) int {
	if err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return he.Code
		}
		return http.StatusInternalServerError
	}
	status := c.Response().Status
	if status == 0 {
		return http.StatusOK
	}
	return status
}

// https://echo.labstack.com/middleware/prometheus/#skipping-certain-urls
func metricsSkipper(c echo.Context) bool {
	path := c.Request().URL.Path
	switch path {
	case "/ping", "/ping/":
		return true
	case "/metrics", "/metrics/":
		return true
	}
	pathItemsWithoutPrefixes := handler_utils.NewPathWithString(path).RemovePrefixes(config.Get().API.Path)
	return pathItemsWithoutPrefixes.StartWithResources(
		[]string{"ping"},
	)
}

func matchedPath(c echo.Context) string {
	if path := c.Path(); path != "" {
		return path
	}
	return c.Request().URL.Path
}
