package factory

import (
	"warehouse/internal/handlers/apiv1/products"
	"warehouse/internal/handlers/apiv1/reservations"
	"warehouse/internal/handlers/apiv1/warehouses"

	"github.com/google/wire"
)

var handlerSet = wire.NewSet(
	reservations.NewHandler,
	products.NewHandler,
	warehouses.NewHandler,
)
