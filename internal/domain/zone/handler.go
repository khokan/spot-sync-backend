package zone

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/httpresponse"
	"spotsync/internal/domain/zone/dto"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid request payload", Errors: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "validation failed", Errors: err.Error()})
	}

	resp, err := h.service.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to create zone", Errors: err.Error()})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessEnvelope{Success: true, Message: "Parking zone created successfully", Data: resp})
}

func (h *handler) UpdateZone(c *echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid zone id", Errors: err.Error()})
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid request payload", Errors: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "validation failed", Errors: err.Error()})
	}

	resp, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorEnvelope{Success: false, Message: "zone not found"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to update zone", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "zone updated successfully", Data: resp})
}

func (h *handler) DeleteZone(c *echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid zone id", Errors: err.Error()})
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorEnvelope{Success: false, Message: "zone not found"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to delete zone", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "zone deleted successfully"})
}

func (h *handler) ListZones(c *echo.Context) error {
	resp, err := h.service.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to list zones", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "zones fetched successfully", Data: resp})
}

func (h *handler) GetZone(c *echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{Success: false, Message: "invalid zone id", Errors: err.Error()})
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.ErrorEnvelope{Success: false, Message: "zone not found"})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{Success: false, Message: "failed to get zone", Errors: err.Error()})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{Success: true, Message: "zone fetched successfully", Data: resp})
}

func parseID(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
