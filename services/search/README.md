# Search service (Go). Full-text + price/brand facets via OpenSearch.

## Run

```bash
docker compose -f infra/docker-compose.yml up -d opensearch
# Ensure catalog has indexed products (status=active)
cd services/search
go mod tidy
go run ./cmd/server
```

## Endpoints

- `GET /health` — readiness (OpenSearch ping)
- `GET /search` — query params:
  - `q` — full-text (name, description, brand; fuzziness AUTO)
  - `category_id` — filter by category (UUID)
  - `brand` — filter by brand
  - `min_price`, `max_price` — price range
  - `limit` (default 20), `offset` (default 0)
- Response: `{ hits, total, facets: { brands: [{ key, count }], prices: [{ from, to, count }] } }`

## Env

- `PORT` (default 8083)
- `OPENSEARCH_URL` (default http://localhost:9200)
