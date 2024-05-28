package warehouses

import (
	"fmt"
	"net/http"
	"strconv"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

// HandleTotalQuantityCount возвращает количество доступных для резервирования товаров на складе,
// но не проверяет, записан ли такой склад в бд, для несуществуюших по задумке вернется 0
func (h *Handler) HandleTotalQuantityCount(c echo.Context) error {
	warehouseIdParam := c.Param("warehouse_id")
	warehouseId, err := strconv.Atoi(warehouseIdParam)

	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("Got error parsing warehouse id: %w", err))
	}

	if warehouseId < 1 {
		return errors.ErrValidation(fmt.Sprintf("Warehouse Id should be positive: %d", warehouseId))
	}

	count, err := h.inventoriesRepo.SumAvailableQuantityByWarehouseId(uint(warehouseId))

	if err != nil {
		return fmt.Errorf("Error getting products quantity count by warehouse id: %w", err)
	}

	return c.JSON(http.StatusOK, struct {
		Count int64 `json:"count"`
	}{
		Count: count,
	})
}
