package reservations

import (
	"warehouse/internal/presenters"
	"warehouse/internal/usecases"
	"warehouse/pkg/log"

	"github.com/go-playground/validator"
)

type Handler struct {
	log                            log.Logger
	validate                       *validator.Validate
	makeReservationsUseCase        usecases.MakeReservationsUseCase
	reservationsPresenter          presenters.ReservationsApiPresenter
	updateReservationStatusUseCase usecases.UpdateReservationStatusUseCase
}

func NewHandler(
	log log.Logger,
	validate *validator.Validate,
	makeReservationsUseCase usecases.MakeReservationsUseCase,
	reservationsPresenter presenters.ReservationsApiPresenter,
	updateReservationStatusUseCase usecases.UpdateReservationStatusUseCase,
) *Handler {
	return &Handler{
		log:                            log,
		validate:                       validate,
		makeReservationsUseCase:        makeReservationsUseCase,
		reservationsPresenter:          reservationsPresenter,
		updateReservationStatusUseCase: updateReservationStatusUseCase,
	}
}
