package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"marketplace/catalog/internal/domain"
)

const indexName = "products"

// ProductDoc for OpenSearch. Flat fields for full-text and facets.
type ProductDoc struct {
	ID          string            `json:"id"`
	CategoryID  string            `json:"category_id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Brand       string            `json:"brand"`
	SellerID    string            `json:"seller_id"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// Client for OpenSearch indexing. Index products on create/update; delete on delete.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates OpenSearch client. baseURL e.g. http://localhost:9200.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IndexProduct indexes or updates product document. Idempotent.
func (c *Client) IndexProduct(ctx context.Context, p *domain.Product, attrs map[string]string) error {
	doc := ProductDoc{
		ID:          p.ID.String(),
		CategoryID:  p.CategoryID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Brand:       p.Brand,
		SellerID:    p.SellerID,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Attributes:  attrs,
	}
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	url := c.baseURL + "/" + indexName + "/_doc/" + p.ID.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch index: %s %s", resp.Status, string(b))
	}
	return nil
}

// DeleteProduct removes product document from index.
func (c *Client) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	url := c.baseURL + "/" + indexName + "/_doc/" + id.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch delete: %s %s", resp.Status, string(b))
	}
	return nil
}

// EnsureIndex creates index with mapping if not exists. Optional on startup.
func (c *Client) EnsureIndex(ctx context.Context) error {
	url := c.baseURL + "/" + indexName
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	mapping := `{"mappings":{"properties":{"id":{"type":"keyword"},"category_id":{"type":"keyword"},"name":{"type":"text"},"slug":{"type":"keyword"},"description":{"type":"text"},"price":{"type":"float"},"brand":{"type":"keyword"},"seller_id":{"type":"keyword"},"status":{"type":"keyword"},"created_at":{"type":"date"},"updated_at":{"type":"date"},"attributes":{"type":"object","enabled":false}}}}`
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader([]byte(mapping)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opensearch create index: %s %s", resp.Status, string(b))
	}
	return nil
}
