package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"marketplace/search/internal/opensearch"
)

// Register mounts GET /search + health.
func Register(e *echo.Echo, osClient *opensearch.Client) {
	e.GET("/health", health(osClient))
	e.GET("/search", search(osClient))
}

func health(osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := osClient.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func search(osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := opensearch.SearchRequest{
			Q:          c.QueryParam("q"),
			CategoryID: c.QueryParam("category_id"),
			Brand:      c.QueryParam("brand"),
		}
		if s := c.QueryParam("min_price"); s != "" {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				req.MinPrice = &v
			}
		}
		if s := c.QueryParam("max_price"); s != "" {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				req.MaxPrice = &v
			}
		}
		req.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
		req.Offset, _ = strconv.Atoi(c.QueryParam("offset"))
		if req.Limit <= 0 {
			req.Limit = 20
		}
		if req.Offset < 0 {
			req.Offset = 0
		}
		result, err := osClient.Search(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "search failed"})
		}
		return c.JSON(http.StatusOK, result)
	}
}
