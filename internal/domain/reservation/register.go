package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)
	jwtService := auth.NewJWTService(cfg.JwtSecret, cfg.JWTExpiryHours)

	reservations := e.Group("/api/v1/reservations", middlewares.AuthMiddleware(jwtService))
	reservations.POST("", handler.CreateReservation)
	reservations.GET("/my-reservations", handler.MyReservations)
	reservations.DELETE("/:id", handler.CancelReservation)
	reservations.GET("", handler.ListAllReservations, middlewares.RequireRole("admin"))
}
