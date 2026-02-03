# Marketplace (Ozon-like)

Monorepo: Go + Next.js, DDD, Clean Architecture. Turborepo + pnpm.

## Stack

- **Backend**: Go 1.23+ (auth, user, catalog, order), Echo, pgx, JWT
- **Frontend**: Next.js 15 App Router, React Server Components, Tailwind, shadcn/ui
- **DB**: PostgreSQL 16, Redis 7, OpenSearch, MinIO (infra/docker-compose)

## Quick start

```bash
pnpm install
docker compose -f infra/docker-compose.yml up -d
pnpm dev                    # apps/web
cd services/auth && go run ./cmd/server    # auth :8080
cd services/user && go run ./cmd/server   # user :8081
cd services/catalog && go run ./cmd/server # catalog :8082
cd services/search && go run ./cmd/server  # search :8083
cd services/seller && go run ./cmd/server  # seller :8084
cd services/order && go run ./cmd/server   # order :8085
cd services/admin && go run ./cmd/server   # admin :8086
```

## Structure

- `apps/web` — Next.js (buyer + seller + admin)
- `apps/mobile` — React Native stub
- `packages/ui`, `types`, `utils`, `config`, `eslint-config` — shared
- `services/auth`, `user`, `catalog`, `search`, `order` — Go/TS services
- `infra/` — docker-compose, k8s stubs

## Phase 2: Auth + User (v0.2)

- **Auth**: Register, Login, Refresh (JWT + refresh token rotation), bcrypt, migrations (users, refresh_tokens)
- **User**: GET/PATCH /users/me (profile), PostgreSQL + Redis, JWT validation (same secret as Auth), user_profiles migration
- Both services: Echo, pgx, health checks, rate limit

## Phase 4: Catalog (v0.4)

- **Catalog**: Products, categories, attributes — full CRUD (REST), PostgreSQL migrations, OpenSearch indexing on product create/update/delete (status=active)
- Endpoints: GET/POST /categories, GET/PUT/DELETE /categories/:id; GET/POST /attributes, GET/PUT/DELETE /attributes/:id; GET/POST /products, GET/PUT/DELETE /products/:id
- Product attributes: product_attribute_values table; optional attributes in create/update; index in OpenSearch for search

## Phase 5: Search API (v0.5)

- **Search**: GET /search — full-text (multi_match on name^2, description, brand, fuzziness AUTO), filters: category_id, brand, min_price, max_price; aggregations: brands (terms), price_ranges (0-500, 500-1000, 1000-5000, 5000+); sort by updated_at desc
- Service reads from OpenSearch index `products` (same as catalog writes)

## Phase 6: Buyer frontend (v0.6)

- **apps/web**: Next.js 15 App Router, buyer layout (Header + Footer), route group `(buyer)`
- **Home** `/`: hero + ссылка «Перейти в каталог»
- **Каталог** `/catalog`: Server Component, fetch Search API; фильтры (q, min/max price, бренд из facets); ProductCard grid; пагинация
- **Товар** `/product/[id]`: Server Component, fetch Catalog API GET /products/:id; кнопка «В корзину» (server action addToCart)
- **Корзина** `/cart`: Server Component, чтение cookie; список товаров, удаление (server action removeFromCart); итого; ссылка «Оформить заказ» → /checkout (stub)
- **Cart**: cookie `marketplace_cart`, server actions addToCart / removeFromCart / getCartItems
- **Env**: API_CATALOG_URL, API_SEARCH_URL (см. apps/web/.env.local.example)

## Phase 7: Seller service + cabinet (v0.7)

- **Catalog**: GET /products — добавлен фильтр `seller_id` (для личного кабинета продавца).
- **Seller service** (Go, :8084): BFF с JWT; проксирует в Catalog. Endpoints: GET/POST /seller/products, GET/PUT/DELETE /seller/products/:id; seller_id из JWT sub; проверка владельца при get/update/delete.
- **Web**: страница входа `/login` (email + пароль → Auth API → cookie `marketplace_token`), редирект `?next=/seller`.
- **API routes**: GET/POST /api/seller/products, GET/PUT/DELETE /api/seller/products/[id] — читают cookie и вызывают Seller service с Bearer.
- **Кабинет продавца** `/seller`: layout с проверкой cookie (редирект на /login при отсутствии), нав: Мои товары, Добавить товар. Страницы: /seller (дашборд), /seller/products (таблица товаров), /seller/products/new (форма создания), /seller/products/[id]/edit (форма редактирования). Статус товара: draft / active (модерация — заглушка).

## Phase 8: Orders + checkout + payment stub (v0.8)

- **Order service** (Go, :8085): PostgreSQL (orders, order_items), JWT (buyer_id из sub). Endpoints: POST /orders (body: items[]), GET /orders (мои заказы), GET /orders/:id, POST /orders/:id/pay (заглушка: status=paid, payment_id=stub_<id>).
- **Web**: /checkout — корзина + «Оформить заказ» (server action placeOrder: cart → Order API → clearCart → redirect /orders/:id). /orders — список заказов (требует логин). /orders/[id] — детали заказа + кнопка «Оплатить (заглушка)». API routes: GET/POST /api/orders, GET /api/orders/[id], POST /api/orders/[id]/pay.
- **Cart**: clearCart() в actions; placeOrder() в lib/actions/order.ts.

## Phase 9: Admin basic panel (v0.9)

- **Admin service** (Go, :8086): один PostgreSQL (marketplace), защита по заголовку X-Admin-Key. Endpoints: GET /admin/users (id, email, created_at), GET /admin/products (query: status), PATCH /admin/products/:id/status (body: { status: "draft"|"active" }).
- **Web**: /admin — вход по секретному ключу (форма → POST /api/admin/login → cookie admin_secret). Layout /admin проверяет cookie === ADMIN_SECRET, иначе редирект на /admin/login. Страницы: /admin (дашборд), /admin/users (таблица пользователей), /admin/products (таблица товаров, фильтр по статусу, кнопки «Одобрить» / «В черновик»). API: POST /api/admin/login, GET /api/admin/users, GET /api/admin/products, PATCH /api/admin/products/[id]/status.
- **Env**: ADMIN_API_URL, ADMIN_SECRET (то же значение, что ADMIN_API_KEY у Admin service).

## Дополнительно

- **CORS**: В сервисах Catalog и Search включено middleware CORS (AllowOrigins: `*`) для публичных API.
- **Тесты (apps/web)**: Vitest, `pnpm test` — тесты в `**/*.test.{ts,tsx}`; пример: `lib/utils.test.ts` (проверка `cn` из `@marketplace/utils`).
- **Observability**: В `infra/docker-compose.yml` добавлены сервисы Prometheus (порт 9090) и Grafana (порт 3001) под профилем `observability`. Запуск: `docker compose -f infra/docker-compose.yml --profile observability up -d`. Конфиг Prometheus: `infra/prometheus.yml` (при необходимости добавьте scrape_configs для Go-сервисов с эндпоинтом `/metrics`).
- **Логи**: Go-сервисы пишут структурированные JSON-логи через `log/slog` (уровень Info по умолчанию). Для сбора в Loki/ELK можно направить stdout в агрегатор.
