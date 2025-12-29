package v1

import (
	"itkdemo/internal/usecase"
	"itkdemo/pkg/config"
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

func NewWalletRoutes(V1 *echo.Group, t *usecase.WalletUseCase) {

	h := NewHandler(t)

	if config.Debug {
		V1.POST("/wallets", h.CreateWallet)
		V1.DELETE("/wallet", h.DeleteWallet)
	}
	V1.POST("/wallet", h.Operation)
	V1.GET("/wallets/:id", h.GetBalance)

	V1.GET("/debug/pprof/", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	V1.GET("/debug/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	V1.GET("/debug/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	V1.GET("/debug/pprof/cmd", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	V1.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	V1.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	V1.GET("/debug/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	V1.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
}
