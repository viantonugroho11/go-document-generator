# go-document-generator — Document Generation Service (PDF/CSV)

This repository provides a template-driven document generation service following Clean Architecture. Core capabilities:

- Generate documents from templates and template versions
- Validate payloads using JSON Schema (compiled schema cached with LRU)
- Outputs supported: PDF (via wkhtmltopdf) and CSV (via text/template)
- Store document metadata in PostgreSQL
- Publish event to Kafka after a document is created
- Redis as a cache/utility client
- HTTP server (Echo) for health check and example `users` CRUD

Note: HTTP endpoints for documents are not yet exposed; document usecases and infrastructure are implemented and ready to be wired.

## Key Features
- Template & Versioning:
  - Entities: `document_templates`, `document_template_versions`, `documents`
  - Payload validation against the version’s `schema` using JSON Schema
  - `sample_payload` for example data
- Generator:
  - CSV: `text/template` with CSV helper functions
  - PDF: `wkhtmltopdf` via `github.com/SebastiaanKlippert/go-wkhtmltopdf`
  - Generator selection based on `output_format` and `engine`
- Integrations:
  - Kafka (producer/consumer)
  - Redis client
  - GORM + PostgreSQL
- Clean Architecture:
  - `usecase` separated from `repository` and `infrastructure`
  - HTTP transport separated from domain/usecase

## Folder Structure
```
cmd/
  app/                 # main HTTP application (Echo)
  consumer/            # Kafka worker/consumer
internal/
  config/              # configuration (loader, DSN helpers)
  entity/              # domain entities (documents, templates, users, etc.)
  infrastructure/
    database/postgres/ # GORM connection + migrations (AutoMigrate user)
    broker/kafka/      # producer & consumer wrappers
    cache/redis/       # Redis client
    documents/         # selector + generators (csv, pdf)
  repository/          # repository interfaces & implementations (gorm/postgres)
  shared/              # utilities (csv helpers, validators)
  transport/
    apis/              # Echo router + example users handler
    event/kafka/       # runner consumer & example handler
  usecase/
    documents/         # services to generate & create documents
    users/             # example users usecase
database/              # SQL definitions for documents & templates
configs/
  config.yaml          # example configuration (optional)
Dockerfile
docker-compose.yml
```

## Prerequisites
- Go (as specified in `go.mod`)
- PostgreSQL
- Kafka broker (+ Zookeeper as per compose)
- Redis
- wkhtmltopdf (required for PDF output)
  - macOS (Homebrew): `brew install wkhtmltopdf`
  - Linux: install the `wkhtmltopdf` package provided by your distro

## Configuration
Configuration is loaded using `github.com/viantonugroho11/go-config-library`. Sources:
- Consul (optional) via `CONSUL_URL` with app name `go-document-generator`
- Local file (default search path in code: `./config`)
- Environment variables

Sample local file included (`configs/config.yaml`):
```yaml
port: "8080"

# Database
databaseurl: ""
dbhost: "postgres"
dbport: "5432"
dbuser: "postgres"
dbpassword: "postgres"
dbname: "appdb"
dbsslmode: "disable"

# Kafka
kafkabrokers: "kafka:9092"
kafkaclientid: "go-document-generator"
kafkagroupid: "go-document-generator-group"
kafkatopic: "user-events"

# Redis
redisaddr: "redis:6379"
redispassword: ""
redisdb: "0"
```

Common environment variables (see `docker-compose.yml`):
- `PORT` (default 8080)
- `DATABASE_URL` OR `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `KAFKA_BROKERS`, `KAFKA_CLIENT_ID`, `KAFKA_GROUP_ID`, `KAFKA_TOPIC`
- `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`

Helpers in code:
- `Configuration.PGDSN()` for Postgres DSN
- `Configuration.KafkaBrokersList()` to get the brokers list

## Database Schema
SQL definitions are under `database/`:
- `document_templates.sql`
- `document-template-versions.sql`
- `documents.sql`

The app currently performs `AutoMigrate` for the `users` table (see `internal/infrastructure/database/postgres/connection.go`). Run the SQL above manually for the document schema, or add automatic migrations per your preference.

## Document Generators
Selector: `internal/infrastructure/documents/factory.go`
- `output_format = "pdf"` → `WKHTMLToPDFGenerator`
- `output_format = "csv"` → `TmplCSVGenerator`

CSV (`text/template`):
```go
// internal/infrastructure/documents/csv/generator.go
bytes, contentType, err := csvGen.Generate(ctx, templateSource, payload)
// contentType: "text/csv"
```

PDF (wkhtmltopdf):
```go
// internal/infrastructure/documents/pdf/generator.go
bytes, contentType, err := pdfGen.Generate(ctx, htmlTemplateSource, payload)
// contentType: "application/pdf"
```

Simple CSV template example:
```txt
name,amount
{{ .Name }},{{ .Amount }}
```

Simple HTML template example for PDF:
```html
<!doctype html>
<html><body>
  <h1>Invoice #{{ .InvoiceNo }}</h1>
  <p>Total: {{ .Total }}</p>
</body></html>
```

## Document Usecases
`internal/usecase/documents/service.go`:
- `GenerateByVersionID(ctx, versionID, payload)` → returns document bytes + `contentType` and the associated `DocumentTemplate`
- `Create(ctx, doc)` → validates payload with the template version’s `schema`, saves the document, then publishes to Kafka

JSON Schema validation:
- Implementation: `internal/shared/validators/schema.go` using `github.com/santhosh-tekuri/jsonschema/v6`
- Compiled schema is cached using LRU (`github.com/hashicorp/golang-lru/v2`)

## Run Locally
Install dependencies (optional; `go mod tidy` will fetch automatically):
```bash
go mod tidy
```

Run HTTP server:
```bash
go run ./cmd/app
```

Run Kafka consumer (optional):
```bash
go run ./cmd/consumer
```

On server startup:
- Initialize GORM connection to PostgreSQL
- AutoMigrate the `users` table (example)
- Initialize Redis client
- Initialize Kafka producer & consumer (consumer runs in background)
- Serve HTTP (Echo)

## Docker & Compose
Build image:
```bash
docker build -t go-document-generator:latest .
```

Run the full stack:
```bash
docker compose up -d --build
```

Ports:
- App: `8080`
- Postgres: `5432`
- Kafka: `9092`
- Redis: `6379`

Ensure `wkhtmltopdf` is available on the host/container when PDF output is used.

## HTTP Endpoints (Current)
Base URL: `http://localhost:${PORT}`

- `GET /healthz` → health check
- User CRUD (example):
  - `POST /users`
  - `GET /users`
  - `GET /users/:id`
  - `PUT /users/:id`
  - `DELETE /users/:id`

Document-related endpoints are not yet exposed. Use `usecase/documents` within services or add routes/handlers as needed.

## Troubleshooting
- Dependencies: run `go mod tidy`
- Postgres connection: verify DSN/env
- Kafka unavailable: check `KAFKA_BROKERS` and broker status
- Redis connection failed: check `REDIS_ADDR`
- PDF generation failed: ensure `wkhtmltopdf` is installed and executable

## License
MIT — feel free to use and modify.

# Go Boilerplate – Clean Architecture

A structured Go boilerplate (clean-architecture inspired) using:
- Echo (HTTP server)
- GORM + PostgreSQL (ORM & database)
- Kafka (producer & consumer) via `github.com/IBM/sarama`
- Redis client via `github.com/redis/go-redis/v9`
- Viper for configuration (file + environment override)

Goals: clear folder layout, separation of concerns (usecase, repository, transport, infrastructure), and explicit dependency wiring.

## Prerequisites
- Go (per `go.mod`, Go 1.23+)
- PostgreSQL
- Kafka broker (+ Zookeeper for the provided compose)
- Redis

## Folder Structure
```
cmd/
  app/                 # application entrypoint (main)
internal/
  config/              # configuration loader (Viper) & DSN helpers
  entity/              # domain entities (e.g., User)
  infrastructure/
    database/
      postgres/        # GORM connection & migration
        connection.go
        migrate/       # optional SQL examples
    broker/
      kafka/           # Kafka producer & consumer wrappers
    cache/
      redis/           # Redis client initialization
  repository/
    user/
      model/           # GORM model mappings
      postgres/        # repository implementation with GORM
      user_repository.go  # repository interface
  transport/
    apis/              # HTTP router
    http/
      dto/             # request/response DTOs
      handler/         # Echo handlers
    event/
      kafka/           # consumer runner & example handlers
  usecase/             # application/service layer
configs/
  config.yaml          # default configuration (overridable by env)
Dockerfile
docker-compose.yml
```

## Configuration (Viper)
Configuration is read from `configs/config.yaml` and can be overridden by environment variables (Viper’s `AutomaticEnv`).

HTTP:
- `PORT` (default: `8080`)

PostgreSQL (choose one approach):
- `DATABASE_URL` (e.g., `postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable`)
- or separate fields:
  - `DB_HOST` (default: `127.0.0.1`)
  - `DB_PORT` (default: `5432`)
  - `DB_USER` (default: `postgres`)
  - `DB_PASSWORD` (default: empty)
  - `DB_NAME` (default: `appdb`)
  - `DB_SSLMODE` (default: `disable`)

Kafka:
- `KAFKA_BROKERS` (default: `127.0.0.1:9092`, comma-separated)
- `KAFKA_CLIENT_ID` (default: `go-document-generator`)
- `KAFKA_GROUP_ID` (default: `go-document-generator-group`)
- `KAFKA_TOPIC` (default: `user-events`)

Redis:
- `REDIS_ADDR` (default: `127.0.0.1:6379`)
- `REDIS_PASSWORD` (default: empty)
- `REDIS_DB` (default: `0`)

See `configs/config.yaml` for a docker-friendly baseline.

## Run Locally
1) Install dependencies (optional; `go mod tidy` will fetch them):
```bash
go get gorm.io/gorm gorm.io/driver/postgres
go get github.com/labstack/echo/v4
go get github.com/IBM/sarama
go get github.com/redis/go-redis/v9
go get github.com/google/uuid
go get github.com/spf13/viper
go mod tidy
```

2) Set environment variables (example):
```bash
export PORT=8080
export DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"
export KAFKA_BROKERS="127.0.0.1:9092"
export KAFKA_CLIENT_ID="go-document-generator"
export KAFKA_GROUP_ID="go-document-generator-group"
export KAFKA_TOPIC="user-events"
export REDIS_ADDR="127.0.0.1:6379"
export REDIS_DB="0"
```

3) Run the server:
```bash
go run ./cmd/app
```

At startup the app will:
- Initialize a GORM connection to PostgreSQL
- AutoMigrate the `users` table
- Start Echo HTTP server
- Initialize Redis client
- Initialize Kafka producer & consumer (consumer runs in background)

## Docker & Compose
Build the image:
```bash
docker build -t go-document-generator:latest .
```

Run with compose (app + Postgres + Zookeeper + Kafka + Redis):
```bash
docker compose up -d --build
```

Exposed ports:
- App: `8080`
- Postgres: `5432`
- Kafka: `9092`
- Redis: `6379`

## HTTP Endpoints
Base URL: `http://localhost:${PORT}`

Healthcheck:
```bash
GET /healthz
```

User CRUD:
- `POST /users`
  - Body:
    ```json
    { "name": "Jane", "email": "jane@example.com" }
    ```
- `GET /users`
- `GET /users/:id`
- `PUT /users/:id`
  - Body:
    ```json
    { "name": "Jane Updated", "email": "jane.updated@example.com" }
    ```
- `DELETE /users/:id`

Example cURL:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane","email":"jane@example.com"}'
```

## Kafka
Producer:
- Wrapper at `internal/infrastructure/broker/kafka/producer.go`, use `Publish(ctx, topic, key, value)`.

Consumer:
- Wrapper at `internal/infrastructure/broker/kafka/consumer.go` using a consumer group.
- Registration & start in `internal/transport/event/kafka/consumer_runner.go` (example handler logs messages).
- Wired in `cmd/app/main.go` via group and topic.

Notes:
- Ensure `KAFKA_BROKERS` points to a running broker.
- Replace `ExampleHandler` to call real usecases as needed.

## Redis
Client:
- Initialization in `internal/infrastructure/cache/redis/client.go`.
- Wired in `cmd/app/main.go`. You can inject it into layers for caching, rate limiting, etc.

## Database
ORM:
- GORM model for `users` in `internal/repository/user/model/user.go`.
- AutoMigrate runs on startup.

Optional SQL:
- Example file at `internal/infrastructure/database/postgres/migrate/init_users.sql`.

## Repository & Usecase
- Repository interface: `internal/repository/user/user_repository.go`
- Postgres (GORM) implementation: `internal/repository/user/postgres/repository.go`
- Usecase: `internal/usecase/user_usecase.go`
- HTTP handlers (Echo): `internal/transport/http/handler/user_handler.go`
- HTTP DTOs: `internal/transport/http/dto/`

## Architecture Notes
- `usecase` depends only on the repository interface
- `repository/*/postgres` adapts the interface with GORM
- `transport/http` (Echo) calls usecases
- `infrastructure` holds I/O details (database, kafka, redis)

## Troubleshooting
- Build dependency issues: run `go mod tidy`
- Postgres connection issues: verify `DATABASE_URL` or `DB_*` variables
- Kafka broker unavailable: verify `KAFKA_BROKERS` and broker status
- Redis connection refused: verify `REDIS_ADDR` and Redis status

## License
MIT. Feel free to use and modify.

# go-document-generator (English)