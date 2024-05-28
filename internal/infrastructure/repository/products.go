package repository

import (
	"warehouse/internal/config"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
)

type ProductsRepositoryPostgres struct {
	db     *gorm.DB
	config config.Config
}

func NewProductsRepositoryPostgres(db *gorm.DB, config config.Config) *ProductsRepositoryPostgres {
	return &ProductsRepositoryPostgres{db: db, config: config}
}

func (s *ProductsRepositoryPostgres) ProductById(productId int) (*models.Product, error) {
	var product models.Product

	return &product, s.db.Where("id = ?", productId).First(&product).Error
}

func (s *ProductsRepositoryPostgres) Create(products []*models.Product) (models.ProductCollection, error) {
	err := s.db.CreateInBatches(&products, repository.CreateProductsBatchSize).Error

	return models.ProductCollection(products), err
}

func (s *ProductsRepositoryPostgres) FindBySKUCodes(codes []string) (models.ProductCollection, error) {
	var products []*models.Product

	err := s.db.Table("products").
		Where("sku IN ?", codes).
		Find(&products).
		Error

	return models.ProductCollection(products), err
}
