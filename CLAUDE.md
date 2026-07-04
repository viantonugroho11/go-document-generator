# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

**go-document-generator** — document generation API

Stack: Echo (HTTP), GORM + PostgreSQL, Kafka via `github.com/viantonugroho11/go-lib/kafka`, Redis, Viper via `go-config-library`.

## Commands

```bash
go mod tidy
go build ./...
go vet ./...
go run ./cmd/app
go run ./cmd/consumer -consumer=<name>   # if consumers exist
docker compose up -d --build               # local deps (Postgres, Kafka, Redis)
```

Config defaults live in `configs/config.yaml`; env vars override via automatic binding (`DATABASE_URL`, `KAFKA_BROKERS`, `REDIS_ADDR`, etc.).

## Architecture

Clean-architecture layering — one vertical slice per aggregate:

```
cmd/app/                     HTTP entrypoint
cmd/consumer/                Kafka consumer entrypoint (-consumer flag)
internal/
  entity/<aggregate>/        Domain structs (no framework deps)
  repository/<aggregate>/    Interface + model/ + postgres/
  repository/begin/          Transaction manager (Begin/Commit/Rollback)
  usecase/<aggregate>/       Application services + events.go (publisher ports)
  transport/apis/            Echo router, handlers, DTOs
  transport/event/           Kafka payloads and consumers
  infrastructure/            kafka, redis, postgres adapters
  bootstrap/                 Config load, wire.go, server setup
database/                    DDL (*.sql), openapi.yaml when present
```

**Dependency rule:** `transport` → `usecase` → `repository` interface. Wire concrete adapters only in `bootstrap/wire.go`.

**Transactions:** usecases call `txManager.Begin(ctx)`, pass `*gorm.DB` to repos, `Commit` on success, deferred `Rollback` on error. Read-only calls pass `nil` tx.

**New domain checklist:** entity → repository → usecase (+ events.go if Kafka) → handler/dto → `router.go` → `wire.go` → GORM model in `Migrate()`.

## Conventions

- Match existing aggregate naming (entity package may be plural; repository folder often singular).
- HTTP contract: prefer `database/openapi.yaml` when present.
- Do not put GORM tags or DB logic in `internal/entity/*`.
- Do not commit `.env` or secrets.
- Respond in the user's language when specified.
