package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"marketplace/order/internal/domain"
	"marketplace/order/internal/middleware"
	"marketplace/order/internal/repository"
)

func Register(e *echo.Echo, pool *pgxpool.Pool, orderRepo *repository.OrderRepository, jwtSecret string) {
	e.GET("/health", health(pool))
	orders := e.Group("/orders", middleware.JWTValidate([]byte(jwtSecret)))
	orders.POST("", createOrder(orderRepo))
	orders.GET("", listOrders(orderRepo))
	orders.GET("/:id", getOrder(orderRepo))
	orders.POST("/:id/pay", payOrder(orderRepo))
}

func health(pool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func createOrder(repo *repository.OrderRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		buyerID := middleware.GetBuyerID(c)
		if buyerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		var req struct {
			Items []struct {
				ProductID string  `json:"product_id"`
				Name      string  `json:"name"`
				Price     float64 `json:"price"`
				Quantity  int     `json:"quantity"`
			} `json:"items"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if len(req.Items) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "items required"})
		}
		var total float64
		items := make([]domain.OrderItem, 0, len(req.Items))
		for _, it := range req.Items {
			if it.ProductID == "" || it.Name == "" || it.Price < 0 || it.Quantity < 1 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid item"})
			}
			total += it.Price * float64(it.Quantity)
			items = append(items, domain.OrderItem{
				ProductID: it.ProductID,
				Name:      it.Name,
				Price:     it.Price,
				Quantity:  it.Quantity,
			})
		}
		now := time.Now().UTC()
		o := &domain.Order{
			ID:        uuid.New(),
			BuyerID:   buyerID,
			Status:    "pending",
			Total:     total,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := repo.Create(c.Request().Context(), o, items); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"id":         o.ID.String(),
			"buyer_id":   o.BuyerID,
			"status":     o.Status,
			"total":      o.Total,
			"created_at": o.CreatedAt,
		})
	}
}

func listOrders(repo *repository.OrderRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		buyerID := middleware.GetBuyerID(c)
		if buyerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		list, err := repo.ListByBuyerID(c.Request().Context(), buyerID, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func getOrder(repo *repository.OrderRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		buyerID := middleware.GetBuyerID(c)
		if buyerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		o, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if o == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		if o.BuyerID != buyerID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		items, _ := repo.GetItems(c.Request().Context(), o.ID)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"id":         o.ID.String(),
			"buyer_id":   o.BuyerID,
			"status":     o.Status,
			"total":      o.Total,
			"payment_id": o.PaymentID,
			"created_at": o.CreatedAt,
			"updated_at": o.UpdatedAt,
			"items":      items,
		})
	}
}

func payOrder(repo *repository.OrderRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		buyerID := middleware.GetBuyerID(c)
		if buyerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		o, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if o == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		if o.BuyerID != buyerID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if o.Status != "pending" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "order already paid or cancelled"})
		}
		paymentID := "stub_" + id.String()
		if err := repo.UpdateStatus(c.Request().Context(), id, "paid", paymentID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status":     "paid",
			"payment_id": paymentID,
		})
	}
}
