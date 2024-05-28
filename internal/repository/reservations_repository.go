package repository

import (
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
)

type ReservationsRepository interface {
	Create(reservations []*models.Reservation, tx ...*gorm.DB) (models.ReservationCollection, error)
	Reservation(reservationId uint) (*models.Reservation, error)
	UpdateStatus(reservationId uint, status string, tx ...*gorm.DB) error
}
