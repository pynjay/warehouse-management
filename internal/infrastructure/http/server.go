package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"warehouse/internal/handlers/apiv1/products"
	"warehouse/internal/handlers/apiv1/reservations"
	"warehouse/internal/handlers/apiv1/warehouses"
	"warehouse/pkg/errors"
	"warehouse/pkg/log"
	"warehouse/pkg/middlewares"

	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	e                   *echo.Echo
	log                 log.Logger
	reservationsHandler *reservations.Handler
	productsHandler     *products.Handler
	warehousesHandler   *warehouses.Handler
}

func NewHTTPServer(
	log log.Logger,
	reservationsHandler *reservations.Handler,
	productsHandler *products.Handler,
	warehousesHandler *warehouses.Handler,
) *HTTPServer {
	return &HTTPServer{
		echo.New(),
		log,
		reservationsHandler,
		productsHandler,
		warehousesHandler,
	}
}

func (server *HTTPServer) Run(port int) {
	server.e.Debug = false
	server.e.HTTPErrorHandler = middlewares.ErrorHandler(server.log)
	echo.NotFoundHandler = func(c echo.Context) error {
		return errors.ErrRouteNotFound()
	}
	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		return errors.ErrMethodNotAllowed()
	}

	server.BuildRouter()

	go func() {
		if err := server.e.Start(fmt.Sprintf(":%d", port)); err != nil {
			server.log.Err(err)
		}
	}()
}

func (server *HTTPServer) Shutdown() {
	err := server.e.Shutdown(context.TODO())
	if err != nil {
		server.log.Err(err)
	}
}

func (server *HTTPServer) handlePing(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Response string `json:"response"`
		Time     string `json:"time"`
	}{
		Response: "ok",
		Time:     time.Now().Format(time.RFC850),
	})
}
