# go-document-generator — Layanan Pembangkitan Dokumen (PDF/CSV)

Repository ini adalah layanan generator dokumen berbasis template dengan pendekatan Clean Architecture. Fungsionalitas inti:

- Pembangkitan dokumen dari template dan versi template
- Validasi payload menggunakan JSON Schema (cache kompilasi schema dengan LRU)
- Output saat ini: PDF (via wkhtmltopdf) dan CSV (via text/template)
- Penyimpanan metadata dokumen di PostgreSQL
- Publikasi event ke Kafka setelah dokumen dibuat
- Redis sebagai klien cache/utility
- HTTP server (Echo) untuk health check dan contoh CRUD `users`

Catatan: endpoint HTTP untuk dokumen belum diekspos; usecase dan infrastruktur dokumen sudah tersedia dan siap di-wire.

## Fitur Utama
- Template & Versioning:
  - Entitas `document_templates`, `document_template_versions`, `documents`
  - Validasi payload terhadap `schema` pada versi template menggunakan JSON Schema
  - `sample_payload` untuk contoh data
- Generator:
  - CSV: `text/template` dengan helper fungsi CSV
  - PDF: `wkhtmltopdf` melalui `github.com/SebastiaanKlippert/go-wkhtmltopdf`
  - Selector pemilihan generator berdasarkan `output_format` dan `engine`
- Integrasi:
  - Kafka (producer/consumer)
  - Redis client
  - GORM + PostgreSQL
- Clean Architecture:
  - `usecase` terpisah dari `repository` dan `infrastructure`
  - transport HTTP terpisah dari domain/usecase

## Struktur Folder
```
cmd/
  app/                 # aplikasi HTTP utama (Echo)
  consumer/            # worker/consumer Kafka
internal/
  config/              # konfigurasi (loader, DSN helpers)
  entity/              # domain entities (documents, templates, users, dst.)
  infrastructure/
    database/postgres/ # koneksi GORM + migrasi (AutoMigrate user)
    broker/kafka/      # producer & consumer wrappers
    cache/redis/       # client Redis
    documents/         # selector + generators (csv, pdf)
  repository/          # interface & implementasi repos (gorm/postgres)
  shared/              # util (csv helpers, validators)
  transport/
    apis/              # router Echo + handler contoh users
    event/kafka/       # runner consumer & example handler
  usecase/
    documents/         # service untuk generate & create dokumen
    users/             # contoh usecase users
database/              # SQL definisi tabel dokumen & template
configs/
  config.yaml          # contoh konfigurasi (opsional)
Dockerfile
docker-compose.yml
```

## Persyaratan
- Go (sesuai `go.mod`)
- PostgreSQL
- Kafka broker (+ Zookeeper sesuai compose)
- Redis
- wkhtmltopdf (dibutuhkan untuk output PDF)
  - Instalasi macOS (Homebrew): `brew install wkhtmltopdf`
  - Linux: instal paket `wkhtmltopdf` sesuai distro

## Konfigurasi
Konfigurasi dimuat menggunakan `github.com/viantonugroho11/go-config-library`. Sumber konfigurasi:
- Consul (opsional) via env `CONSUL_URL` dengan app name `go-document-generator`
- File lokal (path pencarian default di kode: `./config`)
- Environment variables

Contoh file lokal yang disertakan (`configs/config.yaml`):
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

Variabel lingkungan yang umum (lihat `docker-compose.yml`):
- `PORT` (default 8080)
- `DATABASE_URL` ATAU `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `KAFKA_BROKERS`, `KAFKA_CLIENT_ID`, `KAFKA_GROUP_ID`, `KAFKA_TOPIC`
- `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`

Helper di kode:
- `Configuration.PGDSN()` untuk DSN Postgres
- `Configuration.KafkaBrokersList()` untuk daftar broker

## Skema Database
Definisi SQL ada di direktori `database/`:
- `document_templates.sql`
- `document-template-versions.sql`
- `documents.sql`

Aplikasi saat ini melakukan `AutoMigrate` untuk tabel `users` (lihat `internal/infrastructure/database/postgres/connection.go`). Jalankan SQL di atas secara manual untuk skema dokumen, atau tambahkan migrasi otomatis sesuai preferensi Anda.

## Generator Dokumen
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

Contoh template CSV sederhana:
```txt
name,amount
{{ .Name }},{{ .Amount }}
```

Contoh template HTML sederhana untuk PDF:
```html
<!doctype html>
<html><body>
  <h1>Invoice #{{ .InvoiceNo }}</h1>
  <p>Total: {{ .Total }}</p>
</body></html>
```

## Usecase Dokumen
`internal/usecase/documents/service.go`:
- `GenerateByVersionID(ctx, versionID, payload)` → menghasilkan bytes dokumen + `contentType` dan mengembalikan `DocumentTemplate` terkait
- `Create(ctx, doc)` → validasi payload dengan `schema` versi template, simpan dokumen, lalu publish ke Kafka

Validasi JSON Schema:
- Implementasi: `internal/shared/validators/schema.go` menggunakan `github.com/santhosh-tekuri/jsonschema/v6`
- Kompilasi schema dicache menggunakan LRU (`github.com/hashicorp/golang-lru/v2`)

## Menjalankan Secara Lokal
Instal dependensi (opsional, `go mod tidy` akan menarik otomatis):
```bash
go mod tidy
```

Jalankan server HTTP:
```bash
go run ./cmd/app
```

Jalankan consumer Kafka (opsional):
```bash
go run ./cmd/consumer
```

Pada startup server:
- Inisialisasi koneksi GORM ke PostgreSQL
- AutoMigrate tabel `users` (contoh)
- Inisialisasi Redis client
- Inisialisasi Kafka producer & consumer (consumer berjalan background)
- Menyajikan HTTP server (Echo)

## Docker & Compose
Build image:
```bash
docker build -t go-document-generator:latest .
```

Jalankan seluruh stack:
```bash
docker compose up -d --build
```

Port:
- App: `8080`
- Postgres: `5432`
- Kafka: `9092`
- Redis: `6379`

Pastikan `wkhtmltopdf` tersedia pada host/container bila output PDF digunakan.

## HTTP Endpoints (Saat Ini)
Base URL: `http://localhost:${PORT}`

- `GET /healthz` → cek kesehatan
- User CRUD (contoh):
  - `POST /users`
  - `GET /users`
  - `GET /users/:id`
  - `PUT /users/:id`
  - `DELETE /users/:id`

Endpoint terkait dokumen belum diekspos. Gunakan `usecase/documents` di dalam service atau tambahkan route/handler sesuai kebutuhan.

## Troubleshooting
- Dependensi: jalankan `go mod tidy`
- Koneksi Postgres: verifikasi DSN/env
- Kafka tidak tersedia: cek `KAFKA_BROKERS` dan status broker
- Redis gagal konek: cek `REDIS_ADDR`
- PDF gagal dibuat: pastikan `wkhtmltopdf` terinstal dan dapat dieksekusi

## Lisensi
MIT — silakan gunakan dan modifikasi sesuai kebutuhan.

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

# go-document-generator