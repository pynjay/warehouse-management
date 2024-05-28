package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"warehouse/internal/repository"
	"warehouse/internal/repository/models"
	"warehouse/pkg/log"

	"gorm.io/gorm"
)

type ErrUpdateReservationStatusCode int

const (
	ReservationNotFound ErrUpdateReservationStatusCode = iota
	InvalidReservationStatus
	InvalidReservationUpdateRequest
)

type ErrUpdateReservationStatus struct {
	Message string
	ErrCode ErrUpdateReservationStatusCode
}

func (e ErrUpdateReservationStatus) Error() string {
	return e.Message
}

func (e ErrUpdateReservationStatus) Code() ErrUpdateReservationStatusCode {
	return e.ErrCode
}

type UpdateReservationStatusUseCase interface {
	Invoke(reservationId uint, status string, loggel log.Logger) (*models.Reservation, error)
}

type UpdateReservationStatusUseCaseImpl struct {
	db               *gorm.DB
	reservationsRepo repository.ReservationsRepository
	inventoriesRepo  repository.InventoriesRepository
}

func NewUpdateReservationStatusUseCase(
	db *gorm.DB,
	reservationsRepo repository.ReservationsRepository,
	inventoriesRepo repository.InventoriesRepository,
) *UpdateReservationStatusUseCaseImpl {
	return &UpdateReservationStatusUseCaseImpl{
		db:               db,
		reservationsRepo: reservationsRepo,
		inventoriesRepo:  inventoriesRepo,
	}
}

func (u *UpdateReservationStatusUseCaseImpl) Invoke(
	reservationId uint,
	status string,
	logger log.Logger,
) (*models.Reservation, error) {
	logger.Info(
		fmt.Sprintf(
			"Updating reservation %d with status %s",
			reservationId,
			status,
		),
	)

	var err error

	if !slices.Contains([]string{models.ReservationStatusCancelled, models.ReservationStatusFulfilled}, status) {
		return nil, ErrUpdateReservationStatus{
			fmt.Sprintf("The status provided is not appropriate for update: %s", status),
			InvalidReservationStatus,
		}
	}
	reservation, err := u.reservationsRepo.Reservation(reservationId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUpdateReservationStatus{
				fmt.Sprintf("reservation not found by id: %d", reservationId),
				ReservationNotFound,
			}
		}

		return nil, fmt.Errorf("Error getting a reservation by id: %w", err)
	}

	if slices.Contains([]string{models.ReservationStatusCancelled, models.ReservationStatusFulfilled}, reservation.Status) {
		return nil, ErrUpdateReservationStatus{
			fmt.Sprintf("The reservation already has a complete status: %s", reservation.Status),
			InvalidReservationUpdateRequest,
		}
	}

	tx := u.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	if status == models.ReservationStatusCancelled {
		err = u.inventoriesRepo.FreeReservedQuantity(reservation.ProductId, reservation.WarehouseId, reservation.Quantity, tx)
	} else {
		err = u.inventoriesRepo.UpdateQuantityClaimed(reservation.ProductId, reservation.WarehouseId, reservation.Quantity, tx)
	}

	if err != nil {
		return nil, fmt.Errorf("Error update inventories with reservation's quantity: %w", err)
	}

	err = u.reservationsRepo.UpdateStatus(reservationId, status, tx)

	if err != nil {
		return nil, fmt.Errorf("Error update reservation status: %w", err)
	}

	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("Error committing update reservation status transaction: %w", err)
	}

	return reservation, nil
}
