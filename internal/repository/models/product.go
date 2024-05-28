package models

import "time"

type ProductCollection []*Product

type Product struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"column:name"`
	Size      string    `gorm:"column:size"`
	SKU       string    `gorm:"uniqueIndex;column:sku"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
