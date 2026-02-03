# Auth service (Go). Register, Login, Refresh, PostgreSQL (pgx), Echo, JWT + refresh tokens.

## Run

```bash
# From repo root: start infra first
docker compose -f infra/docker-compose.yml up -d postgres

cd services/auth
go mod tidy
go run ./cmd/server
```

## Endpoints

- `GET /health` — readiness (DB ping)
- `POST /auth/register` — body: `{ "email", "password" }` → tokens + user_id
- `POST /auth/login` — body: `{ "email", "password" }` → tokens
- `POST /auth/refresh` — body: `{ "refresh_token" }` → new token pair (rotation)

## Env

- `PORT` (default 8080)
- `DATABASE_URL` (default local postgres)
- `JWT_SECRET`, `JWT_EXPIRY` (default 24h), `REFRESH_EXPIRY` (default 168h)
