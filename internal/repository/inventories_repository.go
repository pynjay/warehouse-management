package repository

import (
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
)

const (
	TableInventories = "inventories"
)

type ReserveProductsParams struct {
	InventoryId       uint
	QuantityToReserve uint
}

type AvailableQuantities map[uint]struct {
	WarehouseId uint
	Quantity    uint
}

type InventoriesRepository interface {
	Create(warehouse models.Inventory) (*models.Inventory, error)
	SumAvailableQuantityByWarehouseId(warehouseId uint) (int64, error)
	GetAvailableQuantities(productId, quantity uint, tx ...*gorm.DB) (quantitiesByWarehouseId AvailableQuantities, err error)
	ReserveProducts(reserveProductsParams []ReserveProductsParams, tx ...*gorm.DB) error
	FreeReservedQuantity(productId, warehouseId, quantity uint, tx ...*gorm.DB) error
	UpdateQuantityClaimed(productId, warehouseId, quantity uint, tx ...*gorm.DB) error
}
