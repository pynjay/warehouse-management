package repository

import (
	"strings"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoriesRepositoryPostgres struct {
	db *gorm.DB
}

func NewInventoriesRepositoryPostgres(db *gorm.DB) *InventoriesRepositoryPostgres {
	return &InventoriesRepositoryPostgres{db: db}
}

func (s *InventoriesRepositoryPostgres) Create(inventory models.Inventory) (*models.Inventory, error) {
	err := s.db.
		Clauses(clause.OnConflict{
			OnConstraint: "inventories_pkey",
			DoUpdates:    clause.Assignments(map[string]interface{}{"quantity": gorm.Expr("inventories.quantity + EXCLUDED.quantity")}),
		}).
		Create(&inventory).
		Error

	return &inventory, err
}

func (s *InventoriesRepositoryPostgres) SumAvailableQuantityByWarehouseId(warehouseId uint) (int64, error) {
	nResult := struct {
		N int64
	}{}

	err := s.db.Table(repository.TableInventories).
		Where("warehouse_id = ?", warehouseId).
		Select("SUM(available_quantity) as N").
		Scan(&nResult).
		Error

	return nResult.N, err
}

func (s *InventoriesRepositoryPostgres) GetAvailableQuantities(
	productId, quantity uint,
	tx ...*gorm.DB,
) (quantitiesByInventoryId repository.AvailableQuantities, err error) {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}

	quantitiesByInventoryId = repository.AvailableQuantities{}

	var getEnoughRowsToExceedQuantityQuery = `WITH cte AS (
        SELECT
        inventories.id,
        inventories.warehouse_id,
        inventories.available_quantity,
        SUM(available_quantity) OVER (ORDER BY inventories.id) AS RunningTotal
        FROM inventories
        JOIN warehouses ON inventories.warehouse_id=warehouses.id
        WHERE inventories.product_id=?
        AND warehouses.is_available=True
    ),
    cte2 AS (
        SELECT
        id,
        warehouse_id,
        available_quantity,
        LAG(RunningTotal) OVER (ORDER BY id) AS LagRunningTotal
        FROM cte
    )
    SELECT id, warehouse_id, available_quantity
    FROM cte2
    WHERE COALESCE(LagRunningTotal, 0) <= ?`

	rows, err := connection.Raw(getEnoughRowsToExceedQuantityQuery, productId, quantity).Rows()
	var (
		Id                uint
		warehouseId       uint
		availableQuantity uint
	)

	for rows.Next() {
		if err = rows.Scan(&Id, &warehouseId, &availableQuantity); err != nil {
			return repository.AvailableQuantities{}, err
		}
		quantitiesByInventoryId[Id] = struct {
			WarehouseId uint
			Quantity    uint
		}{
			WarehouseId: warehouseId,
			Quantity:    availableQuantity,
		}
	}

	return quantitiesByInventoryId, err
}

func (s *InventoriesRepositoryPostgres) ReserveProducts(
	reserveProductsParams []repository.ReserveProductsParams,
	tx ...*gorm.DB,
) error {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}

	var params []interface{}
	sb := strings.Builder{}
	sb.WriteString("UPDATE inventories SET reserved_quantity = CASE ")

	for _, param := range reserveProductsParams {
		sb.WriteString("WHEN id = ? THEN reserved_quantity + ? ")
		params = append(params, param.InventoryId, param.QuantityToReserve)
	}
	sb.WriteString("ELSE reserved_quantity END ")

	return connection.Exec(sb.String(), params...).Error
}

func (s *InventoriesRepositoryPostgres) FreeReservedQuantity(productId, warehouseId, quantity uint, tx ...*gorm.DB) error {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}

	return connection.
		Table("inventories").
		Where("product_id = ?", productId).
		Where("warehouse_id = ?", warehouseId).
		UpdateColumn("reserved_quantity", gorm.Expr("reserved_quantity - ?", quantity)).
		Error
}

func (s *InventoriesRepositoryPostgres) UpdateQuantityClaimed(productId, warehouseId, quantity uint, tx ...*gorm.DB) error {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}

	return connection.
		Table("inventories").
		Where("product_id = ?", productId).
		Where("warehouse_id = ?", warehouseId).
		UpdateColumn("reserved_quantity", gorm.Expr("reserved_quantity - ?", quantity)).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", quantity)).
		Error
}
