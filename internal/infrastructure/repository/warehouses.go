package repository

import (
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
)

type WarehousesRepositoryPostgres struct {
	db *gorm.DB
}

func NewWarehousesRepositoryPostgres(db *gorm.DB) *WarehousesRepositoryPostgres {
	return &WarehousesRepositoryPostgres{db: db}
}

func (s *WarehousesRepositoryPostgres) Create(warehouse models.Warehouse) (*models.Warehouse, error) {
	err := s.db.Create(&warehouse).Error

	return &warehouse, err
}
