package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const SellerIDKey = "seller_id"

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
			c.Set(SellerIDKey, claims.Subject)
			return next(c)
		}
	}
}

func GetSellerID(c echo.Context) string {
	v, _ := c.Get(SellerIDKey).(string)
	return v
}
