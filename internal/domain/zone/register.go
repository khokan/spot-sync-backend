package zone

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

	zones := e.Group("/api/v1/zones")
	zones.GET("", handler.ListZones)
	zones.GET("/:id", handler.GetZone)

	admin := e.Group("/api/v1/zones", middlewares.AuthMiddleware(jwtService), middlewares.RequireRole("admin"))
	admin.POST("", handler.CreateZone)
	admin.PUT("/:id", handler.UpdateZone)
	admin.DELETE("/:id", handler.DeleteZone)
}
