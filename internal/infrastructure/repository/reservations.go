package repository

import (
	"warehouse/internal/repository/models"

	"gorm.io/gorm"
)

type ReservationsRepositoryPostgres struct {
	db *gorm.DB
}

func NewReservationsRepositoryPostgres(db *gorm.DB) *ReservationsRepositoryPostgres {
	return &ReservationsRepositoryPostgres{db: db}
}

func (s *ReservationsRepositoryPostgres) Create(
	reservations []*models.Reservation,
	tx ...*gorm.DB,
) (models.ReservationCollection, error) {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}
	err := connection.CreateInBatches(&reservations, 50).Error

	return reservations, err
}

func (s *ReservationsRepositoryPostgres) Reservation(reservationId uint) (*models.Reservation, error) {
	var reservation models.Reservation

	return &reservation, s.db.Where("id = ?", reservationId).First(&reservation).Error
}

func (s *ReservationsRepositoryPostgres) UpdateStatus(reservationId uint, status string, tx ...*gorm.DB) error {
	var connection *gorm.DB = s.db

	if len(tx) > 0 {
		connection = tx[0]
	}

	return connection.
		Table("reservations").
		Where("id = ?", reservationId).
		Update("status", status).
		Error
}
