package rest

import (
	"itkdemo/internal/transport/rest/middleware"
	"itkdemo/internal/usecase"

	"github.com/labstack/echo/v4"

	v1 "itkdemo/internal/transport/rest/v1"
)

func NewRouter(app *echo.Echo, t *usecase.WalletUseCase) {
	app.Use(middleware.Logging)
	app.Use(middleware.Recover)

	apiV1Group := app.Group("/api/v1")
	{
		v1.NewWalletRoutes(apiV1Group, t)
	}
}
