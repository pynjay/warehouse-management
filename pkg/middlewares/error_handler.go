package middlewares

import (
	"fmt"
	"runtime/debug"
	"warehouse/pkg/errors"
	errors2 "warehouse/pkg/errors"
	"warehouse/pkg/log"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(logger log.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var stack []byte
		he, ok := err.(errors2.ResponseError)

		if !ok || he.ErrorStatus() == errors2.ErrCodeBaseInternalServer {
			_stack := c.Request().Context().Value(CtxErrorStack)

			if _stack != nil {
				stack = _stack.([]byte)
			} else {
				_, stacker := errors2.DeepestErrorWithStack(err)

				if stacker != nil {
					stack = []byte(fmt.Sprintf("%+v", stacker))
				} else {
					stack = debug.Stack()
				}
			}

			he = errors2.ErrInternalServer(err)
		}

		handleError(c, logger, err.Error(), stack)
		errors.SendJson(he, c)
	}
}

func handleError(c echo.Context, log log.Logger, err interface{}, stack []byte) {

	reqData := fmt.Sprintf(
		"%s <-> %s - %s %s, headers: %v",
		c.RealIP(),
		c.Request().RemoteAddr,
		c.Request().Method,
		c.Request().RequestURI,
		c.Request().Header,
	)

	if len(reqData) > 1024 {
		reqData = reqData[:1024]
	}

	log.Err(
		err,
		"request", reqData,
		"stack", string(stack),
	)
}
