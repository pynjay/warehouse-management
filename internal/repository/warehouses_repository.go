package repository

import (
	"warehouse/internal/repository/models"
)

type WarehousesRepository interface {
	Create(warehouse models.Warehouse) (*models.Warehouse, error)
}
