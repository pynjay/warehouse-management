package usecases

import (
	"fmt"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"
	"warehouse/pkg/log"
)

type CreateWarehouseUseCase interface {
	Invoke(params *CreateWarehouseParams, logger log.Logger) (*models.Warehouse, error)
}

type CreateWarehouseParams struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	IsAvailable *bool  `json:"is_available"`
}

type CreateWarehousesUseCaseImpl struct {
	warehousesRepo repository.WarehousesRepository
}

func NewCreateWarehousesUseCase(
	warehousesRepo repository.WarehousesRepository,
) *CreateWarehousesUseCaseImpl {
	return &CreateWarehousesUseCaseImpl{
		warehousesRepo: warehousesRepo,
	}
}

func (u *CreateWarehousesUseCaseImpl) Invoke(
	params *CreateWarehouseParams,
	logger log.Logger,
) (*models.Warehouse, error) {
	logger.Info(
		fmt.Sprintf(
			"Creating a warehouse: %v",
			*params,
		),
	)

	IsAvailable := true

	if params.IsAvailable != nil {
		IsAvailable = *params.IsAvailable
	}

	warehouseModels, err := u.warehousesRepo.Create(models.Warehouse{Name: params.Name, IsAvailable: IsAvailable})

	if err != nil {
		return &models.Warehouse{}, fmt.Errorf("error inserting a new warehouses into the db: %w", err)
	}

	return warehouseModels, nil
}
