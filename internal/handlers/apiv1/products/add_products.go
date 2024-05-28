package products

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"warehouse/internal/repository/models"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

type addProductsRequest struct {
	WarehouseId uint `json:"warehouse_id" validate:"required,min=1"`
	Quantity    uint `json:"quantity" validate:"required,min=1,max=5000"`
}

func (h *Handler) HandleAddProducts(c echo.Context) error {
	req, err := h.parseAddProductsRequest(c.Request())

	if err != nil {
		return err
	}

	productIdParam := c.Param("product_id")
	productId, err := strconv.Atoi(productIdParam)

	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("Got error parsing product id: %w", err))
	}

	if productId < 1 {
		return errors.ErrValidation(fmt.Sprintf("Product Id should be positive: %d", productId))
	}

	inventoryModel := models.Inventory{
		ProductId:   uint(productId),
		WarehouseId: req.WarehouseId,
		Quantity:    req.Quantity,
	}

	inventory, err := h.inventoriesRepo.Create(inventoryModel)

	if err != nil {
		return fmt.Errorf("Error creating a new inventory entry: %w", err)
	}

	return c.JSON(http.StatusOK, struct {
		WarehouseId uint
		ProductId   uint
	}{
		WarehouseId: inventory.WarehouseId,
		ProductId:   inventory.ProductId,
	})
}

func (h *Handler) parseAddProductsRequest(request *http.Request) (addProductsRequest, error) {
	var req addProductsRequest
	body, err := io.ReadAll(request.Body)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to read add products request body: %w", err))
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to decode add products request: %w", err))
	}

	if err := h.validate.Struct(req); err != nil {
		return req, errors.ErrDataValidation(err)
	}

	return req, nil
}
