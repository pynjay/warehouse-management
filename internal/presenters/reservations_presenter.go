package presenters

import "warehouse/internal/repository/models"

type ReservationApiModel struct {
	Id       uint   `json:"reservation_id"`
	OrderId  uint   `json:"order_id"`
	Quantity uint   `json:"quantity"`
	Status   string `json:"status"`
}

type ReservationsApiPresenter interface {
	PresentReservationsForApi(models.ReservationCollection) []ReservationApiModel
	PresentReservationForApi(*models.Reservation) ReservationApiModel
}

type ReservationsApiPresenterImpl struct {
}

func NewReservationsApiPresenterImpl() *ReservationsApiPresenterImpl {
	return &ReservationsApiPresenterImpl{}
}

func (s *ReservationsApiPresenterImpl) PresentReservationsForApi(reservations models.ReservationCollection) []ReservationApiModel {
	reservationsApi := make([]ReservationApiModel, 0, len(reservations))

	for _, reservation := range reservations {
		reservationsApi = append(reservationsApi, s.PresentReservationForApi(reservation))
	}

	return reservationsApi
}

func (s *ReservationsApiPresenterImpl) PresentReservationForApi(reservation *models.Reservation) ReservationApiModel {
	return ReservationApiModel{
		Id:       reservation.ID,
		OrderId:  reservation.OrderId,
		Quantity: reservation.Quantity,
		Status:   reservation.Status,
	}
}
