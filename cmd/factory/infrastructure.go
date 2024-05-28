package factory

import (
	"warehouse/internal/config"
	"warehouse/internal/handlers/apiv1/products"
	"warehouse/internal/handlers/apiv1/reservations"
	"warehouse/internal/handlers/apiv1/warehouses"
	"warehouse/internal/infrastructure/http"
	log2 "warehouse/pkg/log"

	"github.com/google/wire"
)

var infrustructureSet = wire.NewSet(
	provideLogger,
	provideHTTPServer,
)

func provideLogger(config config.Config) (log2.Logger, error) {
	logger, err := log2.NewZapWrapper("stdout", config.IsDev)

	if err != nil {
		return nil, err
	}

	return logger, nil
}

func provideHTTPServer(
	log log2.Logger,
	config config.Config,
	reservationsHandler *reservations.Handler,
	productsHandler *products.Handler,
	warehousesHandler *warehouses.Handler,
) (*http.HTTPServer, func()) {
	server := http.NewHTTPServer(
		log,
		reservationsHandler,
		productsHandler,
		warehousesHandler,
	)

	server.Run(config.HttpListenaddr)

	return server, func() {
		server.Shutdown()
	}
}
