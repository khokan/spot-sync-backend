package user

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/httpresponse"
	"spotsync/internal/domain/user/dto"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.CreateUser(req)
	if err != nil {

		if errors.Is(err, ErrorAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.ErrorEnvelope{
				Success: false,
				Message: "Failed to create user",
				Errors:  err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Failed to create user",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.SuccessEnvelope{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})

}
func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.LoginUser(req)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.ErrorEnvelope{
				Success: false,
				Message: "Cannot login user",
				Errors:  err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Failed to login user",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})

}

func (h *handler) GetMe(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.ErrorEnvelope{
			Success: false,
			Message: "Cannot get user information",
			Errors:  "missing user id in context",
		})
	}

	email, _ := c.Get("user_email").(string)
	name, _ := c.Get("user_name").(string)
	role, _ := c.Get("user_role").(string)

	return c.JSON(http.StatusOK, httpresponse.SuccessEnvelope{
		Success: true,
		Message: "Profile fetched successfully",
		Data: dto.UserResponse{
			ID:    userID,
			Name:  name,
			Email: email,
			Role:  role,
		},
	})
}
