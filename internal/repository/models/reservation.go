package models

import "time"

type ReservationCollection []*Reservation

const (
	ReservationStatusPending   = "pending"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusFulfilled = "fulfilled"
)

type Reservation struct {
	ID          uint      `gorm:"primaryKey"`
	WarehouseId uint      `gorm:"column:warehouse_id"`
	ProductId   uint      `gorm:"column:product_id"`
	OrderId     uint      `gorm:"column:order_id"`
	Quantity    uint      `gorm:"column:quantity"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
