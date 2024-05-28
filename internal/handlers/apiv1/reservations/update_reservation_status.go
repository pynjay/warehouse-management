package reservations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"warehouse/internal/repository/models"
	"warehouse/internal/usecases"
	"warehouse/pkg/errors"

	"github.com/labstack/echo/v4"
)

type updateReservationStatusRequest struct {
	Status string `json:"status" validate:"required,min=1,max=20"`
}

func (h *Handler) HandleUpdateReservationStatus(c echo.Context) error {
	req, err := h.parseUpdateReservationRequest(c.Request())

	if err != nil {
		return err
	}

	reservationIdAttribute := c.Param("reservation_id")
	reservationId, err := strconv.Atoi(reservationIdAttribute)

	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("Got error parsing reservation id: %w", err))
	}

	if reservationId < 1 {
		return errors.ErrValidation(fmt.Sprintf("Warehouse Id should be positive: %d", reservationId))
	}

	reservation, err := h.updateReservationStatusUseCase.Invoke(uint(reservationId), req.Status, h.log)

	if err != nil {
		errWrapped := fmt.Errorf("Error updating reservation status: %w", err)

		if updateReservationStatusErr, ok := err.(usecases.ErrUpdateReservationStatus); ok {
			switch updateReservationStatusErr.Code() {
			case usecases.ReservationNotFound:
				return errors.ErrResourceNotFound(errWrapped)
			default:
				return errors.ErrBadRequest(errWrapped)
			}
		} else {
			return errWrapped
		}
	}

	reservationAPiModel := h.reservationsPresenter.PresentReservationForApi(reservation)
	reservationAPiModel.Status = req.Status

	return c.JSON(
		http.StatusOK,
		reservationAPiModel,
	)
}

func (h *Handler) parseUpdateReservationRequest(request *http.Request) (updateReservationStatusRequest, error) {
	var req updateReservationStatusRequest
	body, err := io.ReadAll(request.Body)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("Failed to read update reservation status request body: %w", err))
	}

	err = json.Unmarshal(body, &req)

	if err != nil {
		return req, errors.ErrBadRequest(fmt.Errorf("Failed to decode update reservation status request: %w", err))
	}

	if err := h.validate.Struct(req); err != nil {
		return req, errors.ErrDataValidation(err)
	}

	if !slices.Contains(
		[]string{
			models.ReservationStatusPending,
			models.ReservationStatusFulfilled,
			models.ReservationStatusCancelled,
		},
		req.Status,
	) {
		return req, errors.ErrValidation(fmt.Sprintf("%s is not a valid reservation status", req.Status))
	}

	return req, nil
}
