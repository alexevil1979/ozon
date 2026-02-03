# Catalog service (Go). Products, categories, attributes, CRUD, OpenSearch indexing.

## Run

```bash
docker compose -f infra/docker-compose.yml up -d postgres opensearch
cd services/catalog
go mod tidy
go run ./cmd/server
```

## Endpoints

- `GET /health` — readiness (DB ping)
- **Categories**: `GET/POST /categories`, `GET/PUT/DELETE /categories/:id` (query: parent_id, limit, offset)
- **Attributes**: `GET/POST /attributes`, `GET/PUT/DELETE /attributes/:id` (query: limit, offset)
- **Products**: `GET/POST /products`, `GET/PUT/DELETE /products/:id` (query: category_id, status, limit, offset). On create/update with status=active, document is indexed to OpenSearch.

## Env

- `PORT` (default 8082)
- `DATABASE_URL`, `OPENSEARCH_URL` (default http://localhost:9200)
