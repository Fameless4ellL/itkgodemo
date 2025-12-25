package middleware

import (
	"time"

	logger "itkdemo/pkg/log"

	"github.com/labstack/echo/v4"
)

func Logging(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		logger.Log.Infof("%s %s %s", c.Request().Method, c.Request().URL.Path, time.Since(start))
		if err := h(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
