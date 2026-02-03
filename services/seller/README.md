# Seller service (Go). BFF: JWT → Catalog API with seller_id from token.

## Run

```bash
# Catalog must be running (port 8082)
cd services/seller
go mod tidy
go run ./cmd/server
```

## Endpoints (all require Authorization: Bearer <access_token>)

- `GET /health` — ok
- `GET /seller/products` — list my products (query: status, limit, offset)
- `POST /seller/products` — create product (body: category_id, name, slug, description?, price, brand?, status?)
- `GET /seller/products/:id` — get one (must be owner)
- `PUT /seller/products/:id` — update (must be owner)
- `DELETE /seller/products/:id` — delete (must be owner)

## Env

- `PORT` (default 8084)
- `CATALOG_URL` (default http://localhost:8082)
- `JWT_SECRET` (must match Auth service)
