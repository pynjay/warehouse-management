package warehouses

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"warehouse/internal/usecases"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleCreateWarehouse(c echo.Context) error {
	var req usecases.CreateWarehouseParams
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("failed to read create warehouse request body: %w", err))
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("failed to decode create warehouse request: %w", err))
	}

	if err := h.validate.Struct(req); err != nil {
		return errors.ErrDataValidation(err)
	}

	warehouse, err := h.createWarehouseUseCase.Invoke(&req, h.log)

	if err != nil {
		return fmt.Errorf("error adding a new warehouse: %w", err)
	}

	return c.JSON(
		http.StatusOK,
		warehouse,
	)
}
