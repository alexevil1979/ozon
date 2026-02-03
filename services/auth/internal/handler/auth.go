// Package handler — HTTP handlers for auth (register, login, refresh).
package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"marketplace/auth/internal/auth"
)

// RegisterRequest body. Validate: email format, password policy.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest body.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshRequest body.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RegisterHandler creates user and returns tokens.
func RegisterHandler(svc *auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		req.Email = strings.TrimSpace(req.Email)
		req.Password = strings.TrimSpace(req.Password)
		if req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password required"})
		}
		u, tp, err := svc.Register(c.Request().Context(), req.Email, req.Password)
		if err != nil {
			if err == auth.ErrEmailTaken {
				return c.JSON(http.StatusConflict, map[string]string{"error": "email already registered"})
			}
			if err == auth.ErrInvalidCredentials || err == auth.ErrPasswordTooShort || err == auth.ErrPasswordTooLong || err == auth.ErrPasswordWeak {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"user_id":        u.ID.String(),
			"email":          u.Email,
			"access_token":   tp.AccessToken,
			"refresh_token":  tp.RefreshToken,
			"expires_in":     tp.ExpiresIn,
		})
	}
}

// LoginHandler returns tokens.
func LoginHandler(svc *auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password required"})
		}
		tp, err := svc.Login(c.Request().Context(), req.Email, req.Password)
		if err != nil {
			if err == auth.ErrInvalidCredentials {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"access_token":  tp.AccessToken,
			"refresh_token":  tp.RefreshToken,
			"expires_in":     tp.ExpiresIn,
		})
	}
}

// RefreshHandler returns new token pair (rotation).
func RefreshHandler(svc *auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RefreshRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "refresh_token required"})
		}
		tp, err := svc.Refresh(c.Request().Context(), req.RefreshToken)
		if err != nil {
			if err == auth.ErrInvalidRefresh {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"access_token":  tp.AccessToken,
			"refresh_token":  tp.RefreshToken,
			"expires_in":     tp.ExpiresIn,
		})
	}
}

