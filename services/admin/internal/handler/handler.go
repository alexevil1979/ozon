package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, pool *pgxpool.Pool, adminKey string) {
	e.GET("/health", health(pool))
	admin := e.Group("/admin", requireAdminKey(adminKey))
	admin.GET("/users", listUsers(pool))
	admin.GET("/products", listProducts(pool))
	admin.PATCH("/products/:id/status", updateProductStatus(pool))
}

func requireAdminKey(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-Admin-Key")
			if key == "" || key != apiKey {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}

func health(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func listUsers(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		rows, err := pool.Query(c.Request().Context(),
			`SELECT id, email, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
			limit, offset,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		defer rows.Close()
		var list []map[string]interface{}
		for rows.Next() {
			var id uuid.UUID
			var email string
			var createdAt interface{}
			if err := rows.Scan(&id, &email, &createdAt); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
			}
			list = append(list, map[string]interface{}{
				"id":         id.String(),
				"email":      email,
				"created_at": createdAt,
			})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func listProducts(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := c.QueryParam("status")
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		query := `SELECT p.id, p.category_id, p.name, p.slug, p.price, p.brand, p.seller_id, p.status, p.created_at, p.updated_at FROM products p WHERE 1=1`
		args := []interface{}{limit, offset}
		n := 1
		if status != "" {
			query += ` AND p.status = $` + strconv.Itoa(n)
			args = append([]interface{}{status}, args...)
			n++
		}
		query += ` ORDER BY p.updated_at DESC LIMIT $` + strconv.Itoa(n) + ` OFFSET $` + strconv.Itoa(n+1)
		rows, err := pool.Query(c.Request().Context(), query, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		defer rows.Close()
		var list []map[string]interface{}
		for rows.Next() {
			var id, categoryID uuid.UUID
			var name, slug, sellerID, pstatus string
			var price float64
			var brand *string
			var createdAt, updatedAt interface{}
			if err := rows.Scan(&id, &categoryID, &name, &slug, &price, &brand, &sellerID, &pstatus, &createdAt, &updatedAt); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
			}
			m := map[string]interface{}{
				"id": id.String(), "category_id": categoryID.String(), "name": name, "slug": slug,
				"price": price, "seller_id": sellerID, "status": pstatus, "created_at": createdAt, "updated_at": updatedAt,
			}
			if brand != nil {
				m["brand"] = *brand
			}
			list = append(list, m)
		}
		return c.JSON(http.StatusOK, list)
	}
}

func updateProductStatus(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		var req struct {
			Status string `json:"status"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.Status != "draft" && req.Status != "active" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "status must be draft or active"})
		}
		_, err = pool.Exec(c.Request().Context(),
			`UPDATE products SET status = $2, updated_at = now() WHERE id = $1`,
			id, req.Status,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": req.Status})
	}
}
