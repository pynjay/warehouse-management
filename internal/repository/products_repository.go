package repository

import (
	"warehouse/internal/repository/models"
)

const (
	CreateProductsBatchSize = 50
)

type ProductsRepository interface {
	ProductById(id int) (*models.Product, error)
	FindBySKUCodes(codes []string) (models.ProductCollection, error)
	Create(products []*models.Product) (models.ProductCollection, error)
}
