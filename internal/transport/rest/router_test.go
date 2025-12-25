package rest_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"itkdemo/internal/transport/rest"
	"itkdemo/internal/usecase"
)

func TestNewRouter_Registration(t *testing.T) {
	e := echo.New()

	u := &usecase.WalletUseCase{}

	rest.NewRouter(e, u)
	routes := e.Routes()

	assert.Greater(t, len(routes), 0, "router is empty")

	var foundV1 bool
	for _, route := range routes {
		if route.Path == "/api/v1" || (len(route.Path) >= 7 && route.Path[:7] == "/api/v1") {
			foundV1 = true
			break
		}
	}
	assert.True(t, foundV1, "group /api/v1 is not registered")
}
