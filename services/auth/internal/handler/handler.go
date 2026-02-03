// Package handler — HTTP handlers. Health + auth (register, login, refresh).
package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"marketplace/auth/internal/auth"
)

// Register mounts routes on e.
func Register(e *echo.Echo, pool *pgxpool.Pool, authSvc *auth.Service) {
	e.GET("/health", health(pool))
	e.POST("/auth/register", RegisterHandler(authSvc))
	e.POST("/auth/login", LoginHandler(authSvc))
	e.POST("/auth/refresh", RefreshHandler(authSvc))
}

// health returns 200 if DB is reachable (readiness).
func health(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}
