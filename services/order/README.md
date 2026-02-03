# Order service (Go). Orders + checkout + payment stub.

## Run

```bash
docker compose -f infra/docker-compose.yml up -d postgres
cd services/order
go mod tidy
go run ./cmd/server
```

## Endpoints (all require Authorization: Bearer)

- `GET /health` — DB ping
- `POST /orders` — create order (body: { items: [{ product_id, name, price, quantity }] })
- `GET /orders` — list my orders (query: limit, offset)
- `GET /orders/:id` — get one (must be owner)
- `POST /orders/:id/pay` — payment stub: set status=paid, payment_id=stub_<id>

## Env

- `PORT` (default 8085)
- `DATABASE_URL`, `JWT_SECRET` (must match Auth)
