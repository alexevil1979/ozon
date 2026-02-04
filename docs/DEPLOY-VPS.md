# Установка маркетплейса (Ozon-like) на VPS

Документ описывает развёртывание проекта на VPS: каталог `/ssd/www/ozon`, домен `ozon.1tlt.ru`.

## Требования

- **ОС**: Ubuntu 22.04 LTS (или Debian 12)
- **Память**: минимум 2 GB RAM (рекомендуется 4 GB для PostgreSQL + OpenSearch + сервисов)
- **Домен**: `ozon.1tlt.ru` указывает на IP вашего VPS (A-запись)

Установлены:

- Docker и Docker Compose
- Go 1.23+
- Node.js 20+ и pnpm
- Nginx
- Certbot (для SSL)

---

## 1. Подготовка системы

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y docker.io docker-compose-v2 nginx certbot python3-certbot-nginx
sudo systemctl enable docker
sudo usermod -aG docker $USER
# выйти и зайти снова, чтобы применилась группа docker
```

Установка Go (если нет):

```bash
wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Установка Node.js и pnpm:

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
sudo npm install -g pnpm
```

---

## 2. Размещение проекта

```bash
sudo mkdir -p /ssd/www
sudo chown $USER:$USER /ssd/www
cd /ssd/www
git clone https://github.com/alexevil1979/ozon.git
cd ozon
```

Альтернатива: загрузить архив и распаковать в `/ssd/www/ozon`.

---

## 3. Инфраструктура (Docker)

Все сервисы БД и поиска запускаются через Docker Compose. Порты наружу можно не открывать (доступ только с localhost).

```bash
cd /ssd/www/ozon/infra
```

Для продакшена задайте пароль PostgreSQL и сохраните его в `.env` в этой же папке:

```bash
echo "POSTGRES_PASSWORD=ваш_надёжный_пароль_postgres" > .env
echo "POSTGRES_USER=marketplace" >> .env
echo "POSTGRES_DB=marketplace" >> .env
```

Измените в `docker-compose.yml` пароль для контейнера postgres: замените `marketplace_dev_secret` на значение из `POSTGRES_PASSWORD` (или оставьте один пароль в `.env` и подставьте его в compose через `env_file`/переменные). Для простоты можно оставить в compose фиксированный пароль и не использовать `.env`.

Запуск только базовых сервисов (без Prometheus/Grafana):

```bash
docker compose -f docker-compose.yml up -d
```

Проверка:

```bash
docker compose -f docker-compose.yml ps
curl -s http://localhost:5432 || true   # pg не отвечает на HTTP, но порт слушает
redis-cli -h 127.0.0.1 ping
curl -s http://localhost:9200
```

Строка подключения к БД (используйте один и тот же пароль, что в compose):

```text
postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable
```

Сохраните её как `DATABASE_URL` для всех Go-сервисов.

---

## 4. Переменные окружения для бэкенда

Общий секрет JWT для Auth, User, Order, Seller (один и тот же):

```bash
export JWT_SECRET="ваш_длинный_случайный_jwt_secret"
```

Секрет админки (то же значение, что будет в Next.js как `ADMIN_SECRET`):

```bash
export ADMIN_API_KEY="ваш_секретный_ключ_админки"
```

Рекомендуется положить их в единый env-файл или systemd environment.

Пример сводки переменных для всех сервисов:

| Сервис   | PORT | DATABASE_URL | JWT_SECRET | OPENSEARCH_URL | REDIS_URL | CATALOG_URL | ADMIN_API_KEY |
|----------|------|--------------|------------|----------------|-----------|-------------|---------------|
| auth     | 8080 | ✓            | ✓          | —              | —         | —           | —             |
| user     | 8081 | ✓            | ✓          | —              | ✓         | —           | —             |
| catalog  | 8082 | ✓            | —          | ✓              | —         | —           | —             |
| search   | 8083 | —            | —          | ✓              | —         | —           | —             |
| seller   | 8084 | —            | ✓          | —              | —         | http://127.0.0.1:8082 | — |
| order    | 8085 | ✓            | ✓          | —              | —         | —           | —             |
| admin    | 8086 | ✓            | —          | —              | —         | —           | ✓             |

`DATABASE_URL` — одна и та же для auth, user, catalog, order, admin.  
`OPENSEARCH_URL` — `http://127.0.0.1:9200`.  
`REDIS_URL` — `redis://127.0.0.1:6379`.

---

## 5. Запуск Go-сервисов

Из корня репозитория `/ssd/www/ozon`:

```bash
export DATABASE_URL="postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable"
export JWT_SECRET="ваш_jwt_secret"
export ADMIN_API_KEY="ваш_admin_secret"
export OPENSEARCH_URL="http://127.0.0.1:9200"
export REDIS_URL="redis://127.0.0.1:6379"
export CATALOG_URL="http://127.0.0.1:8082"
```

Запуск по одному терминалу на сервис (для теста):

```bash
cd /ssd/www/ozon/services/auth && go run ./cmd/server &
cd /ssd/www/ozon/services/user && go run ./cmd/server &
cd /ssd/www/ozon/services/catalog && go run ./cmd/server &
cd /ssd/www/ozon/services/search && go run ./cmd/server &
cd /ssd/www/ozon/services/seller && go run ./cmd/server &
cd /ssd/www/ozon/services/order && go run ./cmd/server &
cd /ssd/www/ozon/services/admin && go run ./cmd/server &
```

Миграции выполняются при старте (auth, user, catalog, order, admin). После первого запуска проверьте:

```bash
curl -s http://127.0.0.1:8080/health
curl -s http://127.0.0.1:8082/health
curl -s http://127.0.0.1:8086/health
```

Для продакшена лучше собрать бинарники и запускать через systemd (см. раздел «Автозапуск (systemd)» ниже).

---

## 6. Next.js (фронт и API-прокси)

В приложении Next.js все запросы к бэкенду идут с сервера (server-side) по localhost, кроме логина/каталога с клиента — для них в браузере используются `NEXT_PUBLIC_*` URL. На VPS имеет смысл направить публичные URL на тот же домен, чтобы запросы шли на ваш сайт, а не на localhost браузера.

```bash
cd /ssd/www/ozon/apps/web
pnpm install
```

Создайте `.env.local`:

```bash
cp .env.local.example .env.local
```

Отредактируйте `.env.local`:

```env
# Серверные (вызовы с Next.js на VPS)
API_CATALOG_URL=http://127.0.0.1:8082
API_SEARCH_URL=http://127.0.0.1:8083
API_SELLER_URL=http://127.0.0.1:8084
ORDER_API_URL=http://127.0.0.1:8085
ADMIN_API_URL=http://127.0.0.1:8086
ADMIN_SECRET=ваш_admin_secret_как_ADMIN_API_KEY

# Публичные URL (браузер): тот же домен
NEXT_PUBLIC_AUTH_URL=https://ozon.1tlt.ru
NEXT_PUBLIC_API_CATALOG_URL=https://ozon.1tlt.ru
```

Сборка и запуск:

```bash
pnpm build
pnpm start
```

Порт по умолчанию — 3000. Проверка: `curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:3000`.

---

## 7. Nginx и SSL (ozon.1tlt.ru)

Сайт отдаётся через Nginx, Next.js слушает только localhost.

Создайте конфиг:

```bash
sudo nano /etc/nginx/sites-available/ozon.1tlt.ru
```

Содержимое (до выдачи SSL):

```nginx
server {
    listen 80;
    server_name ozon.1tlt.ru;
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Включите сайт и проверьте конфиг:

```bash
sudo ln -s /etc/nginx/sites-available/ozon.1tlt.ru /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

Получение сертификата Let's Encrypt:

```bash
sudo certbot --nginx -d ozon.1tlt.ru
```

Certbot сам добавит редирект на HTTPS и параметры SSL. После этого откройте в браузере: `https://ozon.1tlt.ru`.

---

## 8. Автозапуск (systemd) — опционально

Чтобы после перезагрузки поднимались все Go-сервисы и Next.js, можно завести юниты systemd.

Общая папка и env-файл:

```bash
sudo mkdir -p /etc/ozon
sudo nano /etc/ozon/env
```

В `/etc/ozon/env` (без export):

```text
DATABASE_URL=postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable
JWT_SECRET=ваш_jwt_secret
ADMIN_API_KEY=ваш_admin_secret
OPENSEARCH_URL=http://127.0.0.1:9200
REDIS_URL=redis://127.0.0.1:6379
CATALOG_URL=http://127.0.0.1:8082
```

Пример юнита для Auth:

```bash
sudo nano /etc/systemd/system/ozon-auth.service
```

```ini
[Unit]
Description=Ozon Auth Service
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=/ssd/www/ozon/services/auth
EnvironmentFile=/etc/ozon/env
ExecStart=/usr/local/go/bin/go run ./cmd/server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Аналогично создаются юниты для user, catalog, search, seller, order, admin (свои `WorkingDirectory` и при необходимости доп. переменными). Для продакшена лучше собирать бинарники (`go build -o /ssd/www/ozon/bin/auth ./cmd/server`) и в `ExecStart` указывать `/ssd/www/ozon/bin/auth`.

Юнит для Next.js:

```ini
[Unit]
Description=Ozon Next.js
After=network.target ozon-auth.service ozon-catalog.service

[Service]
Type=simple
WorkingDirectory=/ssd/www/ozon/apps/web
ExecStart=/usr/bin/pnpm start
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Включение и запуск:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ozon-auth ozon-user ozon-catalog ozon-search ozon-seller ozon-order ozon-admin ozon-web
sudo systemctl start ozon-auth ozon-user ozon-catalog ozon-search ozon-seller ozon-order ozon-admin ozon-web
```

---

## 9. Краткий чеклист

1. Установить Docker, Go, Node.js, pnpm, Nginx, certbot.
2. Клонировать репозиторий в `/ssd/www/ozon`.
3. В `infra/` запустить `docker compose up -d` (PostgreSQL, Redis, OpenSearch, MinIO).
4. Задать `DATABASE_URL`, `JWT_SECRET`, `ADMIN_API_KEY` и при необходимости остальные переменные.
5. Запустить все Go-сервисы (auth, user, catalog, search, seller, order, admin).
6. В `apps/web` создать `.env.local`, выполнить `pnpm install`, `pnpm build`, `pnpm start`.
7. Настроить Nginx на `ozon.1tlt.ru` → `http://127.0.0.1:3000`, выдать SSL через certbot.
8. При необходимости оформить systemd-юниты для автозапуска.

Домен для тестов: **ozon.1tlt.ru**, каталог приложения на VPS: **/ssd/www/ozon**.
