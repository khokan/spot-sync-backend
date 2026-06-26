package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || role == "" {
				return c.JSON(http.StatusForbidden, map[string]any{
					"success": false,
					"message": "forbidden",
				})
			}

			if _, exists := allowed[role]; !exists {
				return c.JSON(http.StatusForbidden, map[string]any{
					"success": false,
					"message": "forbidden",
				})
			}

			return next(c)
		}
	}
}
