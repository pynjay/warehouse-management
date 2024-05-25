package repository

import (
	"fmt"
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductsRepositoryPostgres struct {
	db *gorm.DB
}

func NewProductsRepositoryPostgres(db *gorm.DB) *ProductsRepositoryPostgres {
	return &ProductsRepositoryPostgres{db: db}
}

func (s *ProductsRepositoryPostgres) Product(accountId, productId int) (*models.Product, error) {
	var product models.Product

	return &product, s.db.Where("crm_account_id = ? AND product_id = ?", accountId, productId).First(&product).Error
}

func (s *ProductsRepositoryPostgres) ProductById(productId int) (*models.Product, error) {
	var product models.Product

	return &product, s.db.Where("product_id = ?", productId).First(&product).Error
}

// func (s *ProductsRepositoryPostgres) ProductsByUserId(userId int) ([]*models.Product, error) {
// 	var products []*models.Product
//
// 	return products, s.db.Find(&products, userId).Error
// }

func (s *ProductsRepositoryPostgres) Create(product *models.Product) (*models.Product, error) {
	err := s.db.
		Clauses(
			clause.OnConflict{
				DoUpdates: clause.AssignmentColumns([]string{"product_id"}),
			},
		).Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *ProductsRepositoryPostgres) Update(product *models.Product) error {
	err := s.db.Where("product_id = ?", product.ID).Updates(product).Error

	if err != nil {
		return fmt.Errorf("error save product to mysql. %w", err)
	}

	return nil
}
