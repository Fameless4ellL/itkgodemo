package middleware

import (
	"itkdemo/internal/domain"
	logger "itkdemo/pkg/log"

	"github.com/labstack/echo/v4"
)

func Recover(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Errorf("panic recovered: %v", err)
				c.Error(domain.ErrInternal)
			}
		}()
		return h(c)
	}
}
