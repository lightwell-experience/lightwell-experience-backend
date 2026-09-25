package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func serveRouter(req *http.Request) (int, []byte, error) {
	router := echo.New()
	RegisterPing(router)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	response := rr.Result()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	return response.StatusCode, body, err
}

func TestPing(t *testing.T) {
	paths := []string{"/ping", "/ping/"}
	for _, path := range paths {
		req, _ := http.NewRequest("GET", path, nil)
		code, body, err := serveRouter(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, code)

		expected := "{\"message\":\"pong\"}\n"
		assert.Equal(t, expected, string(body))
	}
}
