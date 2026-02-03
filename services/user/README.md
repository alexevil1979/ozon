# User service (Go). Profile (display_name, avatar), PostgreSQL + Redis, JWT from Auth.

## Run

```bash
docker compose -f infra/docker-compose.yml up -d postgres redis
cd services/user
go mod tidy
go run ./cmd/server
```

## Endpoints

- `GET /health` — readiness (DB + Redis ping)
- `GET /users/me` — Bearer token required; returns profile (lazy-create if missing)
- `PATCH /users/me` — body: `{ "display_name?", "avatar_url?" }`; partial update

## Env

- `PORT` (default 8081)
- `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET` (must match Auth service)
