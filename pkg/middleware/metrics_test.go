package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lightwell-experience/lightwell-experience-backend/pkg/instrumentation"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMetricsMiddleware(t *testing.T) {
	metrics := instrumentation.NewMetrics(prometheus.NewRegistry())
	middleware := MetricsMiddleware(metrics)
	assert.NotNil(t, middleware)
}

func TestMetricsMiddlewareSkipper(t *testing.T) {
	testCases := []struct {
		Name     string
		Given    string
		Expected bool
	}{
		{Name: "empty", Given: "/", Expected: false},
		{Name: "ping", Given: "/ping", Expected: true},
		{Name: "ping with /", Given: "/ping/", Expected: true},
		{Name: "metrics", Given: "/metrics", Expected: true},
		{Name: "metrics with /", Given: "/metrics/", Expected: true},
		{Name: "ping as resource", Given: "/api/lightwell-next/v1/ping", Expected: true},
		{Name: "ping as resource with slash", Given: "/api/lightwell-next/v1.0/ping/", Expected: true},
		{Name: "status", Given: "/api/lightwell-next/v1.0/status", Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx := echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, testCase.Given, http.NoBody),
				httptest.NewRecorder())
			assert.Equal(t, testCase.Expected, metricsSkipper(ctx))
		})
	}
}

func TestMapStatus(t *testing.T) {
	testCases := []struct {
		Name     string
		Given    int
		Expected string
	}{
		{Name: "0", Given: 0, Expected: ""},
		{Name: "1xx", Given: http.StatusContinue, Expected: "1xx"},
		{Name: "2xx", Given: http.StatusOK, Expected: "2xx"},
		{Name: "3xx", Given: http.StatusMultipleChoices, Expected: "3xx"},
		{Name: "4xx", Given: http.StatusBadRequest, Expected: "4xx"},
		{Name: "5xx", Given: http.StatusInternalServerError, Expected: "5xx"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Expected, mapStatus(testCase.Given))
		})
	}
}

func TestMetricsMiddlewareWithConfig(t *testing.T) {
	assert.Panics(t, func() {
		MetricsMiddlewareWithConfig(MetricsConfig{})
	})

	config := MetricsConfig{
		Metrics: instrumentation.NewMetrics(prometheus.NewRegistry()),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/ping"
		},
	}

	require.NotPanics(t, func() {
		MetricsMiddlewareWithConfig(config)
	})

	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}

	e := echo.New()
	e.Use(MetricsMiddlewareWithConfig(config))
	path := "/api/lightwell-next/v1.0/status"
	e.Add(http.MethodGet, path, h)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	e.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Ok", resp.Body.String())
}
