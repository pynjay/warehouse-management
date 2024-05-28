package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"
	"warehouse/pkg/log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type ErrMakeReservationCode int

const (
	DuplicateSKU ErrMakeReservationCode = iota
	ProductsDoNotExist
	NotEnoughQuantityInStock
)

type ErrMakeReservation struct {
	Message string
	ErrCode ErrMakeReservationCode
}

func (e ErrMakeReservation) Error() string {
	return e.Message
}

func (e ErrMakeReservation) Code() ErrMakeReservationCode {
	return e.ErrCode
}

type MakeReservationsUseCase interface {
	Invoke(params *MakeReservationsParams, logger log.Logger) (models.ReservationCollection, error)
}

type MakeReservationsParams struct {
	OrderId int `json:"order_id" validate:"required,min=1"`
	Items   []struct {
		SKU      string `json:"product_sku" validate:"required,min=1,max=255"`
		Quantity uint   `json:"quantity" validate:"required,min=1"`
	} `json:"items" validate:"required,min=1,max=50,dive"`
}

type MakeReservationsUseCaseImpl struct {
	inventoriesRepo  repository.InventoriesRepository
	reservationsRepo repository.ReservationsRepository
	productsRepo     repository.ProductsRepository
	db               *gorm.DB
}

func NewMakeReservationsUseCase(
	inventoriesRepo repository.InventoriesRepository,
	reservationsRepo repository.ReservationsRepository,
	productsRepo repository.ProductsRepository,
	db *gorm.DB,
) *MakeReservationsUseCaseImpl {
	return &MakeReservationsUseCaseImpl{
		inventoriesRepo:  inventoriesRepo,
		reservationsRepo: reservationsRepo,
		productsRepo:     productsRepo,
		db:               db,
	}
}

func (u *MakeReservationsUseCaseImpl) Invoke(
	params *MakeReservationsParams,
	logger log.Logger,
) (models.ReservationCollection, error) {
	logger.Info(
		fmt.Sprintf(
			"Creating a reservation: %v",
			*params,
		),
	)

	skuCodes := make([]string, 0, len(params.Items))
	quantityByProductId := make(map[uint]uint, len(params.Items))
	quantityBySKU := make(map[string]uint, len(params.Items))

	for _, paramItem := range params.Items {
		if _, ok := quantityBySKU[paramItem.SKU]; ok {
			return nil, ErrMakeReservation{
				fmt.Sprintf("Invalid params: duplicate SKUs not allowed: %s", paramItem.SKU),
				DuplicateSKU,
			}
		}

		skuCodes = append(skuCodes, paramItem.SKU)
		quantityBySKU[paramItem.SKU] = paramItem.Quantity
	}

	products, err := u.productsRepo.FindBySKUCodes(skuCodes)

	if err != nil {
		return nil, fmt.Errorf("unable to get products by codes: %w", err)
	}

	for _, product := range products {
		quantityByProductId[product.ID] = quantityBySKU[product.SKU]
		delete(quantityBySKU, product.SKU)
	}

	if len(quantityBySKU) > 0 {
		keys := make([]string, 0, len(quantityBySKU))
		for k := range quantityBySKU {
			keys = append(keys, k)
		}

		return nil, ErrMakeReservation{
			fmt.Sprintf("Products with SKUs do not exist: %v", keys),
			ProductsDoNotExist,
		}
	}

	var createdReservations models.ReservationCollection
	retryCount := 0

	for retryCount < 3 {
		// TODO идемпотентность
		createdReservations, err = u.performMakeReservationsTransaction(quantityByProductId, params, logger)

		if err == nil {
			break
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch err.(*pgconn.PgError).Code {
			case pgerrcode.SerializationFailure:
				logger.Warn(
					fmt.Sprintf(
						"Get error serializing access due to concurrent update. Retry count: %d",
						retryCount,
					),
				)
				retryCount++
			default:
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return createdReservations, nil
}

func (u *MakeReservationsUseCaseImpl) performMakeReservationsTransaction(
	quantityByProductId map[uint]uint,
	params *MakeReservationsParams,
	logger log.Logger,
) (models.ReservationCollection, error) {
	tx := u.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	var err error

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	reserveProductsParams := make([]repository.ReserveProductsParams, 0, len(quantityByProductId))
	reservations := make([]*models.Reservation, 0, len(quantityByProductId))

	for productId, productQuantity := range quantityByProductId {
		quantitiesByInventoryId, err := u.inventoriesRepo.GetAvailableQuantities(productId, productQuantity, tx)

		if err != nil {
			return nil, fmt.Errorf("Error getting available quantities of a product: %w", err)
		}

		var sumQuantities uint

		for inventoryId, quantityParam := range quantitiesByInventoryId {
			if sumQuantities >= productQuantity {
				break
			}

			quantityToReserve := min(quantityParam.Quantity, productQuantity-sumQuantities)
			if quantityToReserve == 0 {
				continue
			}

			reserveProductsParams = append(reserveProductsParams, repository.ReserveProductsParams{
				InventoryId:       inventoryId,
				QuantityToReserve: quantityToReserve,
			})
			reservations = append(reservations, &models.Reservation{
				WarehouseId: quantityParam.WarehouseId,
				ProductId:   productId,
				OrderId:     uint(params.OrderId),
				Quantity:    quantityToReserve,
				Status:      models.ReservationStatusPending,
			})
			sumQuantities += quantityParam.Quantity
		}

		if sumQuantities < productQuantity {
			logger.Err(fmt.Sprintf(
				"Not enough quantity in stock for product id %d - Wanted %d and got %d",
				productId,
				productQuantity,
				sumQuantities,
			))

			return nil, ErrMakeReservation{
				fmt.Sprintf("Not enough quantity in stock for product id %d", productId),
				NotEnoughQuantityInStock,
			}
		}
	}

	var createdReservations models.ReservationCollection

	if err = u.inventoriesRepo.ReserveProducts(reserveProductsParams, tx); err != nil {
		return nil, err
	}

	if createdReservations, err = u.reservationsRepo.Create(reservations, tx); err != nil {
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return createdReservations, nil
}
