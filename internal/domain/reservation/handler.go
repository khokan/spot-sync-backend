package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/httpresponse"
	"spotsync/internal/domain/reservation/dto"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorEnvelope{Success: false, Message: "unauthorized"})
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid request payload", Errors: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "validation failed", Errors: err.Error()})
	}

	response, err := h.service.Create(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorEnvelope{Success: false, Message: "zone is full"})
		}
		if errors.Is(err, ErrDuplicateLicensePlate) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorEnvelope{Success: false, Message: "duplicate license plate"})
		}
		if errors.Is(err, ErrReservationNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorEnvelope{Success: false, Message: "zone not found"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to create reservation", Errors: err.Error()})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessEnvelope{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    response,
	})
}

func (h *handler) MyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorEnvelope{Success: false, Message: "unauthorized"})
	}

	rows, err := h.service.ListMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to fetch reservations", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "My reservations retrieved successfully", Data: rows})
}

func (h *handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorEnvelope{Success: false, Message: "unauthorized"})
	}
	role, _ := c.Get("user_role").(string)

	reservationID, err := parseReservationID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid reservation id", Errors: err.Error()})
	}

	err = h.service.Cancel(userID, role, reservationID)
	if err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorEnvelope{Success: false, Message: "reservation not found"})
		}
		if errors.Is(err, ErrReservationForbidden) {
			return c.JSON(http.StatusForbidden, httpresponse.ErrorEnvelope{Success: false, Message: "forbidden"})
		}
		if errors.Is(err, ErrReservationConflict) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorEnvelope{Success: false, Message: "reservation is not active"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to cancel reservation", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "Reservation cancelled successfully"})
}

func (h *handler) ListAllReservations(c *echo.Context) error {
	rows, err := h.service.ListAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to fetch all reservations", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "All reservations retrieved successfully", Data: rows})
}

func parseReservationID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
