package models

type Inventory struct {
	ID                uint `gorm:"primaryKey"`
	WarehouseId       uint `gorm:"column:warehouse_id"`
	ProductId         uint `gorm:"column:product_id"`
	Quantity          uint `gorm:"column:quantity"`
	QuantityReserved  uint `gorm:"column:reserved_quantity"`
	QuantityAvailable uint `gorm:"->"`
}
