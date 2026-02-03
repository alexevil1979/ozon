# Admin service (Go). Read users, list/update product status (moderation).

## Run

```bash
# Same PostgreSQL as Auth/Catalog (marketplace DB)
cd services/admin
go mod tidy
go run ./cmd/server
```

## Endpoints (all require header X-Admin-Key: <ADMIN_API_KEY>)

- `GET /health` — DB ping
- `GET /admin/users` — list users (id, email, created_at); query: limit, offset
- `GET /admin/products` — list products; query: status, limit, offset
- `PATCH /admin/products/:id/status` — body: { "status": "draft"|"active" }

## Env

- `PORT` (default 8086)
- `DATABASE_URL`, `ADMIN_API_KEY`
