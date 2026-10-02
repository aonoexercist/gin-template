# Gin Template

A production-ready Go backend template built with [Gin](https://github.com/gin-gonic/gin), using a **feature-based** project layout, PostgreSQL (GORM), structured logging, graceful shutdown, Docker, and CI.

## Features

- Feature-based (modular) structure that scales with your domains
- Handler → Service → Repository layering with interfaces for easy testing
- Typed application errors and a consistent JSON response format
- Middleware: request ID, structured logging, panic recovery, CORS
- Health endpoints: `/healthz` (liveness) and `/readyz` (readiness)
- Config from environment variables (12-factor)
- Graceful shutdown with signal handling
- SQL migrations with golang-migrate
- Multi-stage Dockerfile (distroless) and Docker Compose
- GitHub Actions CI (vet, race tests, golangci-lint)

## Project Structure

```
.
├── cmd/api/                  # Entry point: dependency wiring only
├── internal/
│   ├── user/                 # Example feature
│   │   ├── handler.go        # HTTP layer + RegisterRoutes
│   │   ├── service.go        # Business logic
│   │   ├── repository.go     # Data access (interface + GORM impl)
│   │   ├── model.go          # DB entity
│   │   ├── dto.go            # Request/response structs + validation
│   │   └── service_test.go
│   ├── health/               # Liveness/readiness handlers
│   └── platform/             # Shared infrastructure
│       ├── config/           # Env config loading
│       ├── database/         # Postgres connection + pool
│       ├── logger/           # slog setup
│       ├── server/           # Router + HTTP server + graceful shutdown
│       ├── middleware/       # RequestID, Logger, Recovery, CORS
│       ├── apperror/         # Typed errors mapped to HTTP status
│       └── response/         # Standard JSON response helpers
├── migrations/               # SQL migrations (up/down)
├── api/openapi.yaml          # OpenAPI spec
├── deployments/              # Dockerfile, docker-compose.yml
├── Makefile
├── .env.example
└── .golangci.yml
```

## Prerequisites

| Tool | Needed for | Install |
|---|---|---|
| Go 1.22+ | Building and running | https://go.dev/dl |
| Docker | Local Postgres, `make up` | https://docs.docker.com/get-docker |
| make | Running the Makefile | Preinstalled on macOS/Linux; Windows: WSL or `choco install make` |
| golang-migrate | `make migrate-up/down` | `brew install golang-migrate` or `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| golangci-lint | `make lint` | `brew install golangci-lint` |

## Getting Started

### 1. Use this template

Click **Use this template** on GitHub (or clone), then rename the module path:

```bash
grep -rl "github.com/yourname/gin-template" . \
  | xargs sed -i 's#github.com/yourname/gin-template#github.com/YOU/REPO#g'
```

On macOS, use `sed -i ''` instead of `sed -i`.

### 2. Configure

```bash
cp .env.example .env
go mod tidy
```

### 3. Run locally (Go on your machine, Postgres in Docker)

```bash
docker compose -f deployments/docker-compose.yml up -d db
make migrate-up
make run
```

### 4. Or run everything in Docker

```bash
make up
```

The API is available at `http://localhost:8080`.

## Configuration

All configuration is via environment variables.

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `development` | `development` or `production` (JSON logs, Gin release mode) |
| `PORT` | `8080` | HTTP port |
| `DATABASE_URL` | required | Postgres connection string |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `SHUTDOWN_TIMEOUT` | `15s` | Max time to drain requests on shutdown |

## Makefile Commands

Run from the project root with `make <target>`.

| Command | Description |
|---|---|
| `make run` | Load `.env` and run the API |
| `make build` | Build a static binary to `bin/api` |
| `make test` | Run tests with race detector and coverage |
| `make lint` | Run golangci-lint |
| `make tidy` | Run `go mod tidy` |
| `make up` | Build and start db, migrations, and API via Docker Compose |
| `make down` | Stop the stack and delete the DB volume |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the last migration |

Override the database URL for migrations:

```bash
make migrate-up DB_URL="postgres://user:pass@host:5432/mydb?sslmode=disable"
```

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Liveness probe |
| GET | `/readyz` | Readiness probe (checks DB) |
| POST | `/api/v1/users` | Create a user |
| GET | `/api/v1/users` | List users (`?page=1&page_size=20`) |
| GET | `/api/v1/users/:id` | Get a user |
| PUT | `/api/v1/users/:id` | Update a user |
| DELETE | `/api/v1/users/:id` | Delete a user |

Full spec: [`api/openapi.yaml`](api/openapi.yaml)

### Example

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ana","email":"ana@example.com"}'
```

Success response:

```json
{
  "success": true,
  "data": { "id": 1, "name": "Ana", "email": "ana@example.com", "created_at": "2026-01-01T00:00:00Z" }
}
```

Error response:

```json
{
  "success": false,
  "error": { "code": "conflict", "message": "email already in use" }
}
```

## Adding a New Feature

1. Create `internal/<feature>/` with `model.go`, `dto.go`, `repository.go`, `service.go`, and `handler.go`.
2. Implement `RegisterRoutes(rg *gin.RouterGroup)` on the handler.
3. Wire it in `cmd/api/main.go`:
   ```go
   orderHandler := order.NewHandler(order.NewService(order.NewRepository(db)))
   router := server.NewRouter(cfg, log, healthHandler, userHandler, orderHandler)
   ```
4. Add a migration:
   ```bash
   migrate create -ext sql -dir migrations -seq create_orders
   ```
5. Add tests and update `api/openapi.yaml`.

### Cross-feature dependencies

When one feature needs another, define a small interface in the **consumer** package and let `main.go` wire the concrete implementation. This avoids circular imports:

```go
// internal/order/service.go
type UserFinder interface {
    Get(ctx context.Context, id uint) (*user.User, error)
}
```

## Conventions

- Dependencies flow one way: `handler → service → repository → database`.
- Services return `*apperror.Error`; handlers call `response.Fail`.
- Always pass `c.Request.Context()` down to the database layer.
- Keep `main.go` free of business logic.
- Configuration comes from environment variables only. Never commit `.env`.

## Testing

```bash
make test
```

Service tests use an in-memory fake repository (see `internal/user/service_test.go`), so no database is required.

## Deployment

Build and run the container image:

```bash
docker build -f deployments/Dockerfile -t gin-template .
docker run -p 8080:8080 -e DATABASE_URL=... -e APP_ENV=production gin-template
```

Use `/healthz` for liveness and `/readyz` for readiness probes in Kubernetes. Run migrations as a separate step (job or init container) rather than on app startup.

## Troubleshooting

| Problem | Fix |
|---|---|
| `DATABASE_URL is required` | Create `.env` with `cp .env.example .env` |
| `make: *** missing separator` | Makefile commands must be indented with a tab, not spaces |
| `make migrate-up` can't connect | Start Postgres first: `docker compose -f deployments/docker-compose.yml up -d db` |
| Docker build fails on `go.sum` | Run `go mod tidy` and commit `go.sum` |

## License

MIT. Add a `LICENSE` file before publishing.