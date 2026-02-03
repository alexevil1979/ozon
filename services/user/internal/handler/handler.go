package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"marketplace/user/internal/domain"
	"marketplace/user/internal/redis"
	"marketplace/user/internal/middleware"
	"marketplace/user/internal/repository"
)

func Register(e *echo.Echo, pool *pgxpool.Pool, rdb *redis.Client, profileRepo *repository.ProfileRepository, jwtSecret string) {
	e.GET("/health", health(pool, rdb))
	me := e.Group("/users/me", middleware.JWTValidate([]byte(jwtSecret)))
	me.GET("", getMe(profileRepo))
	me.PATCH("", patchMe(profileRepo))
}

func health(pool *pgxpool.Pool, rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		if rdb != nil {
			if err := rdb.Ping(c.Request().Context()).Err(); err != nil {
				return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": "redis: " + err.Error()})
			}
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func getMe(profileRepo *repository.ProfileRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.GetUserID(c)
		if userID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		p, err := profileRepo.GetByUserID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if p == nil {
			if err := profileRepo.Create(c.Request().Context(), userID); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
			}
			p = &domain.Profile{UserID: userID}
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id":      p.UserID,
			"display_name": p.DisplayName,
			"avatar_url":   p.AvatarURL,
			"created_at":   p.CreatedAt,
			"updated_at":   p.UpdatedAt,
		})
	}
}

func patchMe(profileRepo *repository.ProfileRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.GetUserID(c)
		if userID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		var req struct {
			DisplayName *string `json:"display_name"`
			AvatarURL   *string `json:"avatar_url"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if err := profileRepo.Update(c.Request().Context(), userID, req.DisplayName, req.AvatarURL); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		p, _ := profileRepo.GetByUserID(c.Request().Context(), userID)
		if p == nil {
			p = &domain.Profile{UserID: userID}
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id":      p.UserID,
			"display_name": p.DisplayName,
			"avatar_url":   p.AvatarURL,
			"updated_at":   p.UpdatedAt,
		})
	}
}
