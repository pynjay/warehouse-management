package usecases

import (
	"fmt"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"
	"warehouse/pkg/log"

	"github.com/google/uuid"
)

type CreateProductsUseCase interface {
	Invoke(params *CreateProductsParams, logger log.Logger) (models.ProductCollection, error)
}

type CreateProductsParams struct {
	Products []struct {
		Name string `json:"name" validate:"required,min=1,max=255"`
		Size string `json:"size" validate:"required,min=1,max=50"`
		// TODO description, price, category, etc
	} `json:"products" validate:"required,min=1,max=50,dive"`
}

type CreateProductsUseCaseImpl struct {
	productsRepo repository.ProductsRepository
}

func NewCreateProductsUseCase(
	productsRepo repository.ProductsRepository,
) *CreateProductsUseCaseImpl {
	return &CreateProductsUseCaseImpl{
		productsRepo: productsRepo,
	}
}

func (u *CreateProductsUseCaseImpl) Invoke(
	params *CreateProductsParams,
	logger log.Logger,
) (models.ProductCollection, error) {
	logger.Info(
		fmt.Sprintf(
			"Creating %d products: %v",
			len(params.Products),
			*params,
		),
	)

	if len(params.Products) == 0 || len(params.Products) > repository.CreateProductsBatchSize {
		return models.ProductCollection{}, fmt.Errorf("invalid amount of products provided to create: %d", len(params.Products))
	}

	products := make([]*models.Product, 0, len(params.Products))

	// TODO предовращать создание дубликатов продуктов
	for _, param := range params.Products {
		products = append(products, &models.Product{Name: param.Name, Size: param.Size, SKU: generateSKU()})
	}

	productModels, err := u.productsRepo.Create(products)

	if err != nil {
		return models.ProductCollection{}, fmt.Errorf("error inserting products into the db: %w", err)
	}

	return productModels, nil
}

func generateSKU() string {
	return fmt.Sprintf("SKU-%s", uuid.New().String())
}
