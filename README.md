# Gin User Service

A REST API for user, role, and permission management built with Gin, GORM, and PostgreSQL.

## Requirements

- Go 1.25 or newer (use a currently supported, patched release)
- PostgreSQL 16 or newer
- Docker with Compose (recommended for local development)

## Local development

Copy the example configuration and replace the development secret:

```sh
cp .env.example .env
docker compose up -d postgres
go run ./cmd --action=migrate
go run ./cmd --action=rollback # rolls back the latest migration; destructive
go run ./cmd --action=run-server
```

Alternatively, run the complete stack in containers:

```sh
SECRET_KEY="a-random-secret-containing-at-least-32-characters" docker compose up --build -d
docker compose run --rm app --action=migrate
```

The service listens on `http://localhost:8000` by default.

## Configuration

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `DB_URL` | Yes | — | PostgreSQL connection URL |
| `SECRET_KEY` | Yes | — | JWT HMAC secret containing at least 32 characters |
| `METRICS_TOKEN` | No | — | Bearer token required by `/metrics` when configured |
| `PORT` | No | `8000` | HTTP listening port |
| `ENV` | No | `development` | Set to `production` for Gin release mode |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Comma-separated browser origins |
| `ACCESS_TOKEN_TTL` | No | `15m` | Go duration for access-token lifetime |

Secrets must be injected at runtime. Do not commit `.env` files or package them into deployment artifacts.

## Commands

```sh
go run ./cmd --action=migrate
go run ./cmd --action=create-super-user
go run ./cmd --action=run-server
go run ./cmd --action=drop-tables # destructive
```

## API

- `GET /health/live` — process liveness
- `GET /health/ready` — database readiness
- `POST /api/v1/auth/token` — obtain an access token
- `/api/v1/users` — administrator-only user CRUD
- `/api/v1/roles` — administrator-only role CRUD
- `/api/v1/permissions` — administrator-only permission CRUD

Use the returned token as `Authorization: Bearer <token>`.

## Verification

```sh
gofmt -w .
go vet ./...
go test -race -coverpkg=./... -cover ./...
INTEGRATION_DB_URL="postgres://postgres:postgres@localhost:5432/user_service_test?sslmode=disable" go test ./tests -run Postgres
go build ./cmd
docker build -t gin-user-service .
```

CI runs formatting, vet, race-enabled tests, coverage, a binary build, dependency vulnerability scanning, and a container build on pull requests and pushes to `master`.

## Deployment notes

- Run `go run ./cmd --action=migrate` as a separate, controlled release step before starting new application instances. Migrations are versioned and tracked in `schema_migrations`.
- Store production configuration in a secret manager and expose it to the service at runtime.
- Deploy behind TLS and an ingress/load balancer with explicit trusted proxy configuration.
- Use `/health/live` for liveness and `/health/ready` for readiness checks.
- Retain the prior release artifact so a failed deployment can be rolled back atomically.
