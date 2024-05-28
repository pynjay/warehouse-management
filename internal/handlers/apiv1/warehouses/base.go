package warehouses

import (
	"warehouse/internal/repository"
	"warehouse/internal/usecases"
	"warehouse/pkg/log"

	"github.com/go-playground/validator"
)

type Handler struct {
	log                    log.Logger
	validate               *validator.Validate
	createWarehouseUseCase usecases.CreateWarehouseUseCase
	inventoriesRepo        repository.InventoriesRepository
}

func NewHandler(
	log log.Logger,
	validate *validator.Validate,
	createWarehouseUseCase usecases.CreateWarehouseUseCase,
	inventoriesRepo repository.InventoriesRepository,
) *Handler {
	return &Handler{
		log:                    log,
		validate:               validate,
		createWarehouseUseCase: createWarehouseUseCase,
		inventoriesRepo:        inventoriesRepo,
	}
}
