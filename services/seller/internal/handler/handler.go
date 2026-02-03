package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"marketplace/seller/internal/catalog"
	"marketplace/seller/internal/middleware"
)

func Register(e *echo.Echo, catalogClient *catalog.Client, jwtSecret string) {
	e.GET("/health", health(catalogClient))
	seller := e.Group("/seller", middleware.JWTValidate([]byte(jwtSecret)))
	seller.GET("/products", listProducts(catalogClient))
	seller.POST("/products", createProduct(catalogClient))
	seller.GET("/products/:id", getProduct(catalogClient))
	seller.PUT("/products/:id", updateProduct(catalogClient))
	seller.DELETE("/products/:id", deleteProduct(catalogClient))
}

func health(catalogClient *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func listProducts(client *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sellerID := middleware.GetSellerID(c)
		if sellerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		status := c.QueryParam("status")
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		list, err := client.ListProducts(c.Request().Context(), sellerID, status, limit, offset)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func createProduct(client *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sellerID := middleware.GetSellerID(c)
		if sellerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		body, err := readJSON(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		out, err := client.CreateProduct(c.Request().Context(), sellerID, body)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		return c.JSON(http.StatusCreated, out)
	}
}

func getProduct(client *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sellerID := middleware.GetSellerID(c)
		if sellerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		id := c.Param("id")
		product, err := client.GetProduct(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		if product == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		if sid, _ := product["seller_id"].(string); sid != sellerID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return c.JSON(http.StatusOK, product)
	}
}

func updateProduct(client *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sellerID := middleware.GetSellerID(c)
		if sellerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		id := c.Param("id")
		product, err := client.GetProduct(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		if product == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		if sid, _ := product["seller_id"].(string); sid != sellerID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		body, err := readJSON(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		delete(body, "id")
		delete(body, "seller_id")
		delete(body, "created_at")
		if err := client.UpdateProduct(c.Request().Context(), id, body); err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		updated, _ := client.GetProduct(c.Request().Context(), id)
		return c.JSON(http.StatusOK, updated)
	}
}

func deleteProduct(client *catalog.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sellerID := middleware.GetSellerID(c)
		if sellerID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		id := c.Param("id")
		product, err := client.GetProduct(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		if product == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		if sid, _ := product["seller_id"].(string); sid != sellerID {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		if err := client.DeleteProduct(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "catalog unavailable"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func readJSON(r io.Reader) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
