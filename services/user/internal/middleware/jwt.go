package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var ErrMissingOrInvalidToken = errors.New("missing or invalid token")

// JWTValidate validates Bearer token and sets user_id in context. Same secret as Auth service.
func JWTValidate(secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header"})
			}
			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}
			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}
			claims, ok := token.Claims.(*jwt.RegisteredClaims)
			if !ok || claims.Subject == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}
			c.Set(string(UserIDKey), claims.Subject)
			return next(c)
		}
	}
}

// GetUserID returns user_id from context (set by JWTValidate). Empty if not set.
func GetUserID(c echo.Context) string {
	v, _ := c.Get(string(UserIDKey)).(string)
	return v
}
