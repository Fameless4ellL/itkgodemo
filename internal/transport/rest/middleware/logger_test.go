package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	logger "itkdemo/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	e := echo.New()
	logger.Init()

	t.Run("Success request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		mw := Logging(handler)
		err := mw(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Handler returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/error", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		testErr := errors.New("something went wrong")
		handler := func(c echo.Context) error {
			return testErr
		}

		mw := Logging(handler)
		err := mw(c)

		assert.NoError(t, err)
	})
}
