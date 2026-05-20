# C2 — Container Diagram

The container diagram breaks the system into **deployable units** and supporting infrastructure.

## Diagram

```mermaid
flowchart TB
    client([Client])

    subgraph platform["Document Generator Platform"]
        api["HTTP API<br/>Go / Echo"]
        consumer["Kafka Consumer<br/>Go"]
        db[(PostgreSQL)]
        kafka[Kafka]
        redis[Redis]
        storage[File Storage]
    end

    webhook[Webhook]
    consul[Consul]

    client -->|HTTPS JSON REST| api
    api --> db
    api --> kafka
    api --> redis
    api --> storage
    api --> webhook
    api --> consul
    consumer --> kafka
    consumer --> db
```

## Containers

| Container | Entrypoint | Responsibility |
|-----------|------------|----------------|
| **HTTP API** | `cmd/app/main.go` → `bootstrap.RunApp()` | Echo routing, DI wiring, graceful shutdown. |
| **Kafka Consumer** | `cmd/consumer/main.go` → `bootstrap.RunConsumer()` | Processes Kafka topics (e.g. user/order). |
| **PostgreSQL** | `database/*.sql` | Schema: `document_templates`, `document_template_versions`, `documents`, `document_render_logs`, `document_callback_attempts`. |
| **Kafka** | `configs/config.yaml` | Topics: `document-events`, `template-events`, `user-events`. |
| **Redis** | `internal/bootstrap/redis.go` | Redis client (ready for cache/rate-limit). |
| **File Storage** | `internal/shared/storage/local.go` | Stores rendered file bytes (PDF/HTML/etc.). |

## Runtime Dependencies

```mermaid
flowchart LR
    app[HTTP API] --> pg[(PostgreSQL)]
    app --> kf[Kafka]
    app --> rd[Redis]
    app --> fs[Local Storage]
    cons[Consumer] --> kf
    cons --> pg
```

## Ports & Configuration

| Service | Default |
|---------|---------|
| HTTP API | `:8080` |
| PostgreSQL | `:5432` |
| Kafka | `:9092` |
| Redis | `:6379` |

See [`configs/config.yaml`](../../configs/config.yaml) and [`docker-compose.yml`](../../docker-compose.yml).
