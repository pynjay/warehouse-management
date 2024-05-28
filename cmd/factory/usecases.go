package factory

import (
	"warehouse/internal/usecases"

	"github.com/google/wire"
)

var usecaseSet = wire.NewSet(
	usecases.NewCreateProductsUseCase,
	wire.Bind(new(usecases.CreateProductsUseCase), new(*usecases.CreateProductsUseCaseImpl)),
	usecases.NewCreateWarehousesUseCase,
	wire.Bind(new(usecases.CreateWarehouseUseCase), new(*usecases.CreateWarehousesUseCaseImpl)),
	usecases.NewMakeReservationsUseCase,
	wire.Bind(new(usecases.MakeReservationsUseCase), new(*usecases.MakeReservationsUseCaseImpl)),
	usecases.NewUpdateReservationStatusUseCase,
	wire.Bind(new(usecases.UpdateReservationStatusUseCase), new(*usecases.UpdateReservationStatusUseCaseImpl)),
)
