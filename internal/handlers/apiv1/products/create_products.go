package products

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"warehouse/internal/usecases"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateProductsTypes(c echo.Context) error {
	req, err := h.parseCreateProductsRequest(c.Request())

	if err != nil {
		return err
	}

	products, err := h.createProductsUsecase.Invoke(&req, h.log)

	if err != nil {
		return fmt.Errorf("error creating products: %w", err)
	}

	return c.JSON(
		http.StatusOK,
		h.productsPresenter.PresentProductsForApi(products),
	)
}

func (h *Handler) parseCreateProductsRequest(request *http.Request) (usecases.CreateProductsParams, error) {
	var req usecases.CreateProductsParams
	body, err := io.ReadAll(request.Body)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to read create products request body: %w", err))
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to decode create products request: %w", err))
	}

	if err := h.validate.Struct(req); err != nil {
		return req, errors.ErrDataValidation(err)
	}

	return req, nil
}
