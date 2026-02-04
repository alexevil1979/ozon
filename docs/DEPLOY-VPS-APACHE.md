# Установка маркетплейса (Ozon-like) на VPS с Apache

Вариант развёртывания, когда на сервере используется **только Apache** (без Nginx). Каталог приложения: `/ssd/www/ozon`, домен: `ozon.1tlt.ru`.

## Требования

- **ОС**: Ubuntu 22.04 LTS (или Debian 12)
- **Память**: минимум 2 GB RAM (рекомендуется 4 GB)
- **Домен**: `ozon.1tlt.ru` указывает на IP вашего VPS (A-запись)
- **Веб-сервер**: Apache 2.4 (вместо Nginx)

Устанавливаются:

- Docker и Docker Compose
- Go 1.23+
- Node.js 20+ и pnpm
- **Apache2**
- Certbot с плагином для Apache

---

## 1. Подготовка системы

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y docker.io docker-compose-v2 apache2 certbot python3-certbot-apache
sudo systemctl enable docker
sudo usermod -aG docker $USER
# выйти и зайти снова, чтобы применилась группа docker
```

Включение модулей Apache для проксирования:

```bash
sudo a2enmod proxy proxy_http headers ssl rewrite
sudo systemctl restart apache2
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

Либо загрузите архив и распакуйте в `/ssd/www/ozon`.

---

## 3. Инфраструктура (Docker)

```bash
cd /ssd/www/ozon/infra
```

Для продакшена смените пароль PostgreSQL в `docker-compose.yml` (значение `POSTGRES_PASSWORD` / пароль в `DATABASE_URL`).

Запуск сервисов:

```bash
docker compose -f docker-compose.yml up -d
```

Проверка:

```bash
docker compose -f docker-compose.yml ps
redis-cli -h 127.0.0.1 ping
curl -s http://localhost:9200
```

Строка подключения к БД (тот же пароль, что в compose):

```text
postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable
```

Используйте её как `DATABASE_URL` для всех Go-сервисов.

---

## 4. Переменные окружения для бэкенда

Общий JWT-секрет для Auth, User, Order, Seller:

```bash
export JWT_SECRET="ваш_длинный_случайный_jwt_secret"
export ADMIN_API_KEY="ваш_секретный_ключ_админки"
```

Сводка по сервисам:

| Сервис   | PORT | DATABASE_URL | JWT_SECRET | OPENSEARCH_URL | REDIS_URL | CATALOG_URL | ADMIN_API_KEY |
|----------|------|--------------|------------|----------------|-----------|-------------|---------------|
| auth     | 8080 | ✓            | ✓          | —              | —         | —           | —             |
| user     | 8081 | ✓            | ✓          | —              | ✓         | —           | —             |
| catalog  | 8082 | ✓            | —          | ✓              | —         | —           | —             |
| search   | 8083 | —            | —          | ✓              | —         | —           | —             |
| seller   | 8084 | —            | ✓          | —              | —         | http://127.0.0.1:8082 | — |
| order    | 8085 | ✓            | ✓          | —              | —         | —           | —             |
| admin    | 8086 | ✓            | —          | —              | —         | —           | ✓             |

`OPENSEARCH_URL` = `http://127.0.0.1:9200`, `REDIS_URL` = `redis://127.0.0.1:6379`.

---

## 5. Запуск Go-сервисов

Из корня репозитория:

```bash
export DATABASE_URL="postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable"
export JWT_SECRET="ваш_jwt_secret"
export ADMIN_API_KEY="ваш_admin_secret"
export OPENSEARCH_URL="http://127.0.0.1:9200"
export REDIS_URL="redis://127.0.0.1:6379"
export CATALOG_URL="http://127.0.0.1:8082"

cd /ssd/www/ozon/services/auth    && go run ./cmd/server &
cd /ssd/www/ozon/services/user     && go run ./cmd/server &
cd /ssd/www/ozon/services/catalog  && go run ./cmd/server &
cd /ssd/www/ozon/services/search    && go run ./cmd/server &
cd /ssd/www/ozon/services/seller    && go run ./cmd/server &
cd /ssd/www/ozon/services/order    && go run ./cmd/server &
cd /ssd/www/ozon/services/admin   && go run ./cmd/server &
```

Проверка:

```bash
curl -s http://127.0.0.1:8080/health
curl -s http://127.0.0.1:8082/health
curl -s http://127.0.0.1:8086/health
```

Для продакшена используйте systemd (раздел 8).

---

## 6. Next.js (фронт и API-прокси)

```bash
cd /ssd/www/ozon/apps/web
pnpm install
cp .env.local.example .env.local
```

Отредактируйте `.env.local`:

```env
API_CATALOG_URL=http://127.0.0.1:8082
API_SEARCH_URL=http://127.0.0.1:8083
API_SELLER_URL=http://127.0.0.1:8084
ORDER_API_URL=http://127.0.0.1:8085
ADMIN_API_URL=http://127.0.0.1:8086
ADMIN_SECRET=ваш_admin_secret_как_ADMIN_API_KEY

NEXT_PUBLIC_AUTH_URL=https://ozon.1tlt.ru
NEXT_PUBLIC_API_CATALOG_URL=https://ozon.1tlt.ru
```

Сборка и запуск:

```bash
pnpm build
pnpm start
```

Next.js слушает порт 3000 на localhost.

---

## 7. Apache и SSL (ozon.1tlt.ru)

Сайт отдаётся через Apache; Next.js доступен только как backend на `127.0.0.1:3000`.

### 7.1. Виртуальный хост (до SSL)

Создайте конфиг сайта:

```bash
sudo nano /etc/apache2/sites-available/ozon.1tlt.ru.conf
```

Содержимое:

```apache
<VirtualHost *:80>
    ServerName ozon.1tlt.ru

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:3000/
    ProxyPassReverse / http://127.0.0.1:3000/

    RequestHeader set X-Forwarded-Proto "http"
</VirtualHost>
```

Включите сайт и проверьте конфигурацию:

```bash
sudo a2ensite ozon.1tlt.ru.conf
sudo apache2ctl configtest
sudo systemctl reload apache2
```

Проверка в браузере: `http://ozon.1tlt.ru` (должен открываться фронт Next.js).

### 7.2. SSL через Certbot (Apache)

Выдача и автоматическая подстановка сертификата Let's Encrypt:

```bash
sudo certbot --apache -d ozon.1tlt.ru
```

Certbot создаст или изменит виртуальный хост для HTTPS (порт 443) и при необходимости настроит редирект с HTTP на HTTPS. После этого используйте `https://ozon.1tlt.ru`.

### 7.3. Ручная настройка HTTPS (если не используете Certbot)

Если сертификат у вас уже есть, можно добавить второй виртуальный хост:

```apache
<VirtualHost *:443>
    ServerName ozon.1tlt.ru

    SSLEngine on
    SSLCertificateFile      /path/to/fullchain.pem
    SSLCertificateKeyFile   /path/to/privkey.pem

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:3000/
    ProxyPassReverse / http://127.0.0.1:3000/

    RequestHeader set X-Forwarded-Proto "https"
</VirtualHost>
```

Включите SSL и перезагрузите Apache:

```bash
sudo a2enmod ssl
sudo systemctl reload apache2
```

---

## 8. Автозапуск (systemd) — опционально

Общий env-файл для Go-сервисов:

```bash
sudo mkdir -p /etc/ozon
sudo nano /etc/ozon/env
```

Содержимое `/etc/ozon/env` (без `export`):

```text
DATABASE_URL=postgres://marketplace:ВАШ_ПАРОЛЬ@127.0.0.1:5432/marketplace?sslmode=disable
JWT_SECRET=ваш_jwt_secret
ADMIN_API_KEY=ваш_admin_secret
OPENSEARCH_URL=http://127.0.0.1:9200
REDIS_URL=redis://127.0.0.1:6379
CATALOG_URL=http://127.0.0.1:8082
```

Пример юнита для Auth (`/etc/systemd/system/ozon-auth.service`):

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

Аналогично создайте юниты для user, catalog, search, seller, order, admin (свои `WorkingDirectory`). Для Next.js:

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

## 9. Краткий чеклист (вариант с Apache)

1. Установить Docker, Go, Node.js, pnpm, **Apache2**, certbot (python3-certbot-apache).
2. Включить модули Apache: `proxy`, `proxy_http`, `headers`, `ssl`, `rewrite`.
3. Клонировать репозиторий в `/ssd/www/ozon`.
4. В `infra/` запустить `docker compose up -d`.
5. Задать переменные окружения и запустить все Go-сервисы.
6. В `apps/web` создать `.env.local`, выполнить `pnpm install`, `pnpm build`, `pnpm start`.
7. Создать виртуальный хост Apache для `ozon.1tlt.ru` с проксированием на `http://127.0.0.1:3000`.
8. Выдать SSL: `sudo certbot --apache -d ozon.1tlt.ru`.
9. При необходимости настроить systemd для автозапуска сервисов.

Домен: **ozon.1tlt.ru**, каталог на VPS: **/ssd/www/ozon**, веб-сервер: **Apache**.
