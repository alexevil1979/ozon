package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const indexName = "products"

// Client for OpenSearch search only (read). Catalog service writes to same index.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SearchRequest query params.
type SearchRequest struct {
	Q          string   // full-text (name, description)
	CategoryID string   // filter
	Brand      string   // filter
	MinPrice   *float64 // filter
	MaxPrice   *float64 // filter
	Limit      int
	Offset     int
}

// SearchResult hit + facets.
type SearchResult struct {
	Hits   []SearchHit   `json:"hits"`
	Total  int64         `json:"total"`
	Facets SearchFacets  `json:"facets"`
}

type SearchHit struct {
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

type SearchFacets struct {
	Brands []BrandBucket  `json:"brands"`
	Prices []PriceBucket  `json:"prices"`
}

type BrandBucket struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type PriceBucket struct {
	From  *float64 `json:"from"`
	To    *float64 `json:"to"`
	Count int64   `json:"count"`
}

// Search runs full-text + filters + aggregations (brand terms, price range).
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Build OpenSearch request body
	body := c.buildSearchBody(req)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := c.baseURL + "/" + indexName + "/_search"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opensearch search: %s %s", resp.Status, string(b))
	}

	var osResp osSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&osResp); err != nil {
		return nil, err
	}
	return c.mapResponse(&osResp, req.Offset), nil
}

func (c *Client) buildSearchBody(req SearchRequest) map[string]interface{} {
	// bool: must status=active, filter by category/brand/price; should multi_match if q
	var must []map[string]interface{}
	var filter []map[string]interface{}

	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{"status": "active"},
	})

	if req.Q != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Q,
				"fields": []string{"name^2", "description", "brand"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	if req.CategoryID != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"category_id": req.CategoryID},
		})
	}
	if req.Brand != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"brand": req.Brand},
		})
	}
	if req.MinPrice != nil || req.MaxPrice != nil {
		rangeQ := map[string]interface{}{}
		if req.MinPrice != nil {
			rangeQ["gte"] = *req.MinPrice
		}
		if req.MaxPrice != nil {
			rangeQ["lte"] = *req.MaxPrice
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{"price": rangeQ},
		})
	}

	boolQ := map[string]interface{}{"must": must}
	if len(filter) > 0 {
		boolQ["filter"] = filter
	}

	// Aggregations: brands (terms), price_ranges (range buckets)
	aggs := map[string]interface{}{
		"brands": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "brand",
				"size":  50,
				"missing": "—",
			},
		},
		"price_ranges": map[string]interface{}{
			"range": map[string]interface{}{
				"field": "price",
				"ranges": []map[string]interface{}{
					{"key": "0-500", "from": 0, "to": 500},
					{"key": "500-1000", "from": 500, "to": 1000},
					{"key": "1000-5000", "from": 1000, "to": 5000},
					{"key": "5000+", "from": 5000},
				},
			},
		},
	}

	return map[string]interface{}{
		"query": map[string]interface{}{"bool": boolQ},
		"aggs":  aggs,
		"from":  req.Offset,
		"size":  req.Limit,
		"sort":  []map[string]interface{}{{"updated_at": map[string]string{"order": "desc"}}},
	}
}

type osSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		Brands struct {
			Buckets []struct {
				Key   string `json:"key"`
				Count int64  `json:"doc_count"`
			} `json:"buckets"`
		} `json:"brands"`
		PriceRanges struct {
			Buckets []struct {
				Key     string  `json:"key"`
				From    *float64 `json:"from"`
				To      *float64 `json:"to"`
				DocCount int64  `json:"doc_count"`
			} `json:"buckets"`
		} `json:"price_ranges"`
	} `json:"aggregations"`
}

func (c *Client) mapResponse(r *osSearchResponse, offset int) *SearchResult {
	out := &SearchResult{
		Total: r.Hits.Total.Value,
		Hits:  make([]SearchHit, 0, len(r.Hits.Hits)),
		Facets: SearchFacets{
			Brands: make([]BrandBucket, 0),
			Prices: make([]PriceBucket, 0),
		},
	}
	for _, h := range r.Hits.Hits {
		hit := mapSourceToHit(h.Source)
		if hit != nil {
			out.Hits = append(out.Hits, *hit)
		}
	}
	for _, b := range r.Aggregations.Brands.Buckets {
		out.Facets.Brands = append(out.Facets.Brands, BrandBucket{Key: b.Key, Count: b.Count})
	}
	for _, b := range r.Aggregations.PriceRanges.Buckets {
		out.Facets.Prices = append(out.Facets.Prices, PriceBucket{From: b.From, To: b.To, Count: b.DocCount})
	}
	return out
}

func mapSourceToHit(m map[string]interface{}) *SearchHit {
	if m == nil {
		return nil
	}
	hit := &SearchHit{}
	if v, ok := m["id"].(string); ok {
		hit.ID = v
	}
	if v, ok := m["category_id"].(string); ok {
		hit.CategoryID = v
	}
	if v, ok := m["name"].(string); ok {
		hit.Name = v
	}
	if v, ok := m["slug"].(string); ok {
		hit.Slug = v
	}
	if v, ok := m["description"].(string); ok {
		hit.Description = v
	}
	if v, ok := m["price"].(float64); ok {
		hit.Price = v
	}
	if v, ok := m["brand"].(string); ok {
		hit.Brand = v
	}
	if v, ok := m["seller_id"].(string); ok {
		hit.SellerID = v
	}
	if v, ok := m["status"].(string); ok {
		hit.Status = v
	}
	if v, ok := m["created_at"].(string); ok {
		t, _ := time.Parse(time.RFC3339, v)
		hit.CreatedAt = t
	}
	if v, ok := m["updated_at"].(string); ok {
		t, _ := time.Parse(time.RFC3339, v)
		hit.UpdatedAt = t
	}
	if v, ok := m["attributes"].(map[string]interface{}); ok {
		hit.Attributes = make(map[string]string)
		for k, val := range v {
			if s, ok := val.(string); ok {
				hit.Attributes[k] = s
			}
		}
	}
	return hit
}

// Ping checks OpenSearch availability.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/_cluster/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("opensearch ping: %s", resp.Status)
	}
	return nil
}
