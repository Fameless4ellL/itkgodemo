package v1

import (
	"itkdemo/internal/usecase"
	"itkdemo/pkg/config"

	"github.com/labstack/echo/v4"
)

func NewWalletRoutes(V1 *echo.Group, t *usecase.WalletUseCase) {

	h := NewHandler(t)

	if config.Debug {
		V1.POST("/wallet", h.CreateWallet)
		V1.DELETE("/wallet", h.DeleteWallet)
	}
	V1.GET("/wallet", h.Operation)
	V1.POST("/wallets/:id", h.GetBalance)
}
