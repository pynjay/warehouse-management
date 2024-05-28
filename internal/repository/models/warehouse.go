package models

import "time"

type Warehouse struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"column:name"`
	IsAvailable bool      `gorm:"column:is_available"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
