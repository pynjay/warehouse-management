package models

type Product struct {
    ID          uint   `gorm:"primaryKey;column:id"`
    SKU         string `gorm:"uniqueIndex;column:sku"`
    Quantity    int `gorm:"column:quantity"`
    Reserved    bool `gorm:"column:quantity"`
    WarehouseID uint `gorm:"column:quantity"`
}
