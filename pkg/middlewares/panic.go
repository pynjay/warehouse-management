package middlewares

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

const (
	CtxErrorStack = "error_stack"
)

func PanicMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					var ok bool

					if err, ok = r.(error); !ok {
						err = fmt.Errorf("%v", r)
					}

					req := c.Request()
					reqCtx := req.Context()
					reqCtx = context.WithValue(reqCtx, CtxErrorStack, debug.Stack())
					req = req.WithContext(reqCtx)
					c.SetRequest(req)

					c.Error(err)
				}
			}()

			return h(c)
		}
	}
}
