package presenters

import "warehouse/internal/repository/models"

type ProductApiModel struct {
	ProductId uint   `json:"product_id"`
	Name      string `json:"name"`
	Size      string `json:"size"`
	SKU       string `json:"sku"`
}

type ProductsApiPresenter interface {
	PresentProductsForApi(models.ProductCollection) []ProductApiModel
}

type ProductsApiPresenterImpl struct {
}

func NewProductsApiPresenterImpl() *ProductsApiPresenterImpl {
	return &ProductsApiPresenterImpl{}
}

func (s *ProductsApiPresenterImpl) PresentProductsForApi(products models.ProductCollection) []ProductApiModel {
	productsApi := make([]ProductApiModel, 0, len(products))

	for _, product := range products {
		productsApi = append(productsApi, ProductApiModel{
			ProductId: product.ID,
			Name:      product.Name,
			Size:      product.Size,
			SKU:       product.SKU,
		})
	}

	return productsApi
}
