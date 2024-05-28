package factory

import (
	"warehouse/internal/presenters"

	"github.com/go-playground/validator"
	"github.com/google/wire"
)

var interfacesSet = wire.NewSet(
	validator.New,
	presenters.NewProductsApiPresenterImpl,
	wire.Bind(new(presenters.ProductsApiPresenter), new(*presenters.ProductsApiPresenterImpl)),
	presenters.NewReservationsApiPresenterImpl,
	wire.Bind(new(presenters.ReservationsApiPresenter), new(*presenters.ReservationsApiPresenterImpl)),
)
