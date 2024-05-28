package reservations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"warehouse/internal/usecases"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleMakeReservations(c echo.Context) error {
	req, err := h.parseMakeReservationsRequest(c.Request())

	if err != nil {
		return err
	}

	reservations, err := h.makeReservationsUseCase.Invoke(&req, h.log)

	if err != nil {
		makeReservationsErr := fmt.Errorf("error creating products: %w", err)

		if _, ok := err.(usecases.ErrMakeReservation); ok {
			return errors.ErrValidation(err.Error())
		} else {
			return makeReservationsErr
		}
	}

	return c.JSON(
		http.StatusOK,
		h.reservationsPresenter.PresentReservationsForApi(reservations),
	)
}

func (h *Handler) parseMakeReservationsRequest(request *http.Request) (usecases.MakeReservationsParams, error) {
	var req usecases.MakeReservationsParams
	body, err := io.ReadAll(request.Body)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to read make reservations request body: %w", err))
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("failed to decode make reservations request: %w", err))
	}

	if err := h.validate.Struct(req); err != nil {
		return req, errors.ErrDataValidation(err)
	}

	return req, nil
}
