package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"marketplace/catalog/internal/domain"
	"marketplace/catalog/internal/opensearch"
	"marketplace/catalog/internal/repository"
)

// Register mounts REST CRUD + health. OpenSearch indexing on product create/update/delete.
func Register(
	e *echo.Echo,
	pool *pgxpool.Pool,
	catRepo *repository.CategoryRepository,
	attrRepo *repository.AttributeRepository,
	prodRepo *repository.ProductRepository,
	osClient *opensearch.Client,
) {
	e.GET("/health", health(pool, osClient))

	e.GET("/categories", listCategories(catRepo))
	e.POST("/categories", createCategory(catRepo))
	e.GET("/categories/:id", getCategory(catRepo))
	e.PUT("/categories/:id", updateCategory(catRepo))
	e.DELETE("/categories/:id", deleteCategory(catRepo))

	e.GET("/attributes", listAttributes(attrRepo))
	e.POST("/attributes", createAttribute(attrRepo))
	e.GET("/attributes/:id", getAttribute(attrRepo))
	e.PUT("/attributes/:id", updateAttribute(attrRepo))
	e.DELETE("/attributes/:id", deleteAttribute(attrRepo))

	e.GET("/products", listProducts(prodRepo))
	e.POST("/products", createProduct(prodRepo, osClient))
	e.GET("/products/:id", getProduct(prodRepo))
	e.PUT("/products/:id", updateProduct(prodRepo, osClient))
	e.DELETE("/products/:id", deleteProduct(prodRepo, osClient))
}

func health(pool *pgxpool.Pool, osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

// --- Categories ---

func listCategories(repo *repository.CategoryRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var parentID *uuid.UUID
		if s := c.QueryParam("parent_id"); s != "" {
			id, err := uuid.Parse(s)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid parent_id"})
			}
			parentID = &id
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		list, err := repo.List(c.Request().Context(), parentID, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func createCategory(repo *repository.CategoryRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Name     string  `json:"name"`
			Slug     string  `json:"slug"`
			ParentID *string `json:"parent_id"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.Name == "" || req.Slug == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and slug required"})
		}
		var parentID *uuid.UUID
		if req.ParentID != nil && *req.ParentID != "" {
			id, err := uuid.Parse(*req.ParentID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid parent_id"})
			}
			parentID = &id
		}
		now := time.Now().UTC()
		cat := &domain.Category{
			ID:        uuid.New(),
			Name:      req.Name,
			Slug:      req.Slug,
			ParentID:  parentID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := repo.Create(c.Request().Context(), cat); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusCreated, cat)
	}
}

func getCategory(repo *repository.CategoryRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		cat, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if cat == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusOK, cat)
	}
}

func updateCategory(repo *repository.CategoryRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		cat, err := repo.GetByID(c.Request().Context(), id)
		if err != nil || cat == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		var req struct {
			Name     string  `json:"name"`
			Slug     string  `json:"slug"`
			ParentID *string `json:"parent_id"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.Name != "" {
			cat.Name = req.Name
		}
		if req.Slug != "" {
			cat.Slug = req.Slug
		}
		if req.ParentID != nil {
			if *req.ParentID == "" {
				cat.ParentID = nil
			} else {
				pid, err := uuid.Parse(*req.ParentID)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid parent_id"})
				}
				cat.ParentID = &pid
			}
		}
		if err := repo.Update(c.Request().Context(), cat); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, cat)
	}
}

func deleteCategory(repo *repository.CategoryRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		if err := repo.Delete(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// --- Attributes ---

func listAttributes(repo *repository.AttributeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		list, err := repo.List(c.Request().Context(), limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func createAttribute(repo *repository.AttributeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Name      string `json:"name"`
			Slug      string `json:"slug"`
			ValueType string `json:"value_type"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.Name == "" || req.Slug == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and slug required"})
		}
		if req.ValueType == "" {
			req.ValueType = "text"
		}
		a := &domain.Attribute{
			ID:        uuid.New(),
			Name:      req.Name,
			Slug:      req.Slug,
			ValueType: req.ValueType,
			CreatedAt: time.Now().UTC(),
		}
		if err := repo.Create(c.Request().Context(), a); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusCreated, a)
	}
}

func getAttribute(repo *repository.AttributeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		a, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if a == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusOK, a)
	}
}

func updateAttribute(repo *repository.AttributeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		a, err := repo.GetByID(c.Request().Context(), id)
		if err != nil || a == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		var req struct {
			Name      string `json:"name"`
			Slug      string `json:"slug"`
			ValueType string `json:"value_type"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.Name != "" {
			a.Name = req.Name
		}
		if req.Slug != "" {
			a.Slug = req.Slug
		}
		if req.ValueType != "" {
			a.ValueType = req.ValueType
		}
		if err := repo.Update(c.Request().Context(), a); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, a)
	}
}

func deleteAttribute(repo *repository.AttributeRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		if err := repo.Delete(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// --- Products ---

func listProducts(repo *repository.ProductRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var categoryID *uuid.UUID
		if s := c.QueryParam("category_id"); s != "" {
			id, err := uuid.Parse(s)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category_id"})
			}
			categoryID = &id
		}
		status := c.QueryParam("status")
		if status == "" {
			status = "active"
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 || limit > 100 {
			limit = 20
		}
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		if offset < 0 {
			offset = 0
		}
		sellerID := c.QueryParam("seller_id")
		list, err := repo.List(c.Request().Context(), categoryID, sellerID, status, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return c.JSON(http.StatusOK, list)
	}
}

func createProduct(repo *repository.ProductRepository, osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			CategoryID  string   `json:"category_id"`
			Name        string   `json:"name"`
			Slug        string   `json:"slug"`
			Description string   `json:"description"`
			Price       float64  `json:"price"`
			Brand       string   `json:"brand"`
			SellerID    string   `json:"seller_id"`
			Status      string   `json:"status"`
			Attributes  []struct {
				AttributeID string `json:"attribute_id"`
				Value      string `json:"value"`
			} `json:"attributes"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.CategoryID == "" || req.Name == "" || req.Slug == "" || req.SellerID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "category_id, name, slug, seller_id required"})
		}
		if req.Price < 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "price must be >= 0"})
		}
		if req.Status == "" {
			req.Status = "draft"
		}
		cid, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category_id"})
		}
		p := &domain.Product{
			ID:          uuid.New(),
			CategoryID:  cid,
			Name:        req.Name,
			Slug:        req.Slug,
			Description: req.Description,
			Price:       req.Price,
			Brand:       req.Brand,
			SellerID:    req.SellerID,
			Status:      req.Status,
		}
		if err := repo.Create(c.Request().Context(), p); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		attrs := make([]domain.ProductAttributeValue, 0, len(req.Attributes))
		attrMap := make(map[string]string)
		for _, a := range req.Attributes {
			aid, err := uuid.Parse(a.AttributeID)
			if err != nil {
				continue
			}
			attrs = append(attrs, domain.ProductAttributeValue{ProductID: p.ID, AttributeID: aid, Value: a.Value})
			attrMap[a.AttributeID] = a.Value
		}
		if len(attrs) > 0 {
			_ = repo.ReplaceAttributeValues(c.Request().Context(), p.ID, attrs)
		}
		if p.Status == "active" && osClient != nil {
			_ = osClient.IndexProduct(c.Request().Context(), p, attrMap)
		}
		return c.JSON(http.StatusCreated, p)
	}
}

func getProduct(repo *repository.ProductRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		p, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if p == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		attrs, _ := repo.GetAttributeValues(c.Request().Context(), id)
		res := map[string]interface{}{
			"id": p.ID, "category_id": p.CategoryID, "name": p.Name, "slug": p.Slug,
			"description": p.Description, "price": p.Price, "brand": p.Brand,
			"seller_id": p.SellerID, "status": p.Status, "created_at": p.CreatedAt, "updated_at": p.UpdatedAt,
			"attributes": attrs,
		}
		return c.JSON(http.StatusOK, res)
	}
}

func updateProduct(repo *repository.ProductRepository, osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		p, err := repo.GetByID(c.Request().Context(), id)
		if err != nil || p == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		var req struct {
			CategoryID  *string  `json:"category_id"`
			Name        *string  `json:"name"`
			Slug        *string  `json:"slug"`
			Description *string  `json:"description"`
			Price       *float64 `json:"price"`
			Brand       *string  `json:"brand"`
			Status      *string  `json:"status"`
			Attributes  *[]struct {
				AttributeID string `json:"attribute_id"`
				Value       string `json:"value"`
			} `json:"attributes"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
		}
		if req.CategoryID != nil {
			cid, err := uuid.Parse(*req.CategoryID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category_id"})
			}
			p.CategoryID = cid
		}
		if req.Name != nil {
			p.Name = *req.Name
		}
		if req.Slug != nil {
			p.Slug = *req.Slug
		}
		if req.Description != nil {
			p.Description = *req.Description
		}
		if req.Price != nil {
			if *req.Price < 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "price must be >= 0"})
			}
			p.Price = *req.Price
		}
		if req.Brand != nil {
			p.Brand = *req.Brand
		}
		if req.Status != nil {
			p.Status = *req.Status
		}
		if err := repo.Update(c.Request().Context(), p); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		var attrMap map[string]string
		if req.Attributes != nil {
			attrs := make([]domain.ProductAttributeValue, 0, len(*req.Attributes))
			attrMap = make(map[string]string)
			for _, a := range *req.Attributes {
				aid, err := uuid.Parse(a.AttributeID)
				if err != nil {
					continue
				}
				attrs = append(attrs, domain.ProductAttributeValue{ProductID: p.ID, AttributeID: aid, Value: a.Value})
				attrMap[a.AttributeID] = a.Value
			}
			_ = repo.ReplaceAttributeValues(c.Request().Context(), p.ID, attrs)
		} else {
			vals, _ := repo.GetAttributeValues(c.Request().Context(), p.ID)
			attrMap = make(map[string]string)
			for _, v := range vals {
				attrMap[v.AttributeID.String()] = v.Value
			}
		}
		if p.Status == "active" && osClient != nil {
			_ = osClient.IndexProduct(c.Request().Context(), p, attrMap)
		} else if p.Status != "active" && osClient != nil {
			_ = osClient.DeleteProduct(c.Request().Context(), p.ID)
		}
		return c.JSON(http.StatusOK, p)
	}
}

func deleteProduct(repo *repository.ProductRepository, osClient *opensearch.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		if err := repo.Delete(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		if osClient != nil {
			_ = osClient.DeleteProduct(c.Request().Context(), id)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
