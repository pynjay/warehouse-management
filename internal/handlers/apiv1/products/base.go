package products

import (
	"warehouse/internal/presenters"
	"warehouse/internal/repository"
	"warehouse/internal/usecases"
	"warehouse/pkg/log"

	"github.com/go-playground/validator"
)

type Handler struct {
	log                   log.Logger
	validate              *validator.Validate
	createProductsUsecase usecases.CreateProductsUseCase
	productsPresenter     presenters.ProductsApiPresenter
	inventoriesRepo       repository.InventoriesRepository
}

func NewHandler(
	log log.Logger,
	validate *validator.Validate,
	createProductsUsecase usecases.CreateProductsUseCase,
	productsPresenter presenters.ProductsApiPresenter,
	inventoriesRepo repository.InventoriesRepository,
) *Handler {
	return &Handler{
		log:                   log,
		validate:              validate,
		createProductsUsecase: createProductsUsecase,
		productsPresenter:     productsPresenter,
		inventoriesRepo:       inventoriesRepo,
	}
}
