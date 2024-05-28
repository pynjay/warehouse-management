package middlewares

import (
	"strings"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

func AcceptJsonMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")

			if strings.ToLower(contentType) != "application/json" {
				return errors.ErrUnsupportedContentType()
			}

			return h(c)
		}
	}
}
