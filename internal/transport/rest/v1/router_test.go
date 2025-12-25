package v1

import (
	"net/http"
	"testing"

	"itkdemo/internal/usecase"
	"itkdemo/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewWalletRoutes(t *testing.T) {
	e := echo.New()
	v1Group := e.Group("/api/v1")

	mockUseCase := &usecase.WalletUseCase{}

	t.Run("Check routes registration when Debug is true", func(t *testing.T) {
		config.Debug = true
		NewWalletRoutes(v1Group, mockUseCase)

		expectedRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodPost, "/api/v1/wallet"},
			{http.MethodDelete, "/api/v1/wallet"},
			{http.MethodGet, "/api/v1/wallet"},
			{http.MethodPost, "/api/v1/wallets/:id"},
		}

		for _, er := range expectedRoutes {
			found := false
			for _, r := range e.Routes() {
				if r.Method == er.method && r.Path == er.path {
					found = true
					break
				}
			}
			assert.True(t, found, "Route %s %s should be registered", er.method, er.path)
		}
	})

	t.Run("Check routes registration when Debug is false", func(t *testing.T) {
		e := echo.New()
		v1Group := e.Group("/api/v1")

		config.Debug = false
		NewWalletRoutes(v1Group, mockUseCase)

		for _, r := range e.Routes() {
			assert.NotEqual(t, "/api/v1/wallets", r.Path, "Post /wallets should not be registered in non-debug mode")
			assert.NotEqual(t, http.MethodDelete, r.Method)
		}
	})
}
