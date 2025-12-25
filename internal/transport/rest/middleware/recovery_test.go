package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	logger "itkdemo/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	e := echo.New()
	logger.Init()

	t.Run("execution", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}

		mw := Recover(handler)
		err := mw(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			panic("something went wrong")
		}

		mw := Recover(handler)
		err := mw(c)
		assert.NoError(t, err)
	})
}
