## Document Generator DB

Enterprise asynchronous document generation — templates, versioning, JSON schema validation, rendering, storage, DMS handoff, webhooks, and audit logs. PostgreSQL DDL + OpenAPI contract.

### Artifacts

| File | Purpose |
|------|---------|
| `dbdiagram.dbml` | ERD (source design) |
| `00_enums.sql` | PostgreSQL enum types (**run first**) |
| `document-templates.sql` | Master templates |
| `document-template-versions.sql` | Versioned content + schema |
| `documents.sql` | Generation jobs / outputs |
| `document-render-logs.sql` | Render attempt diagnostics |
| `document-callback-attempts.sql` | Webhook delivery history |
| `openapi.yaml` | REST API contract |

### Execution order

```bash
./run-all-sql.sh --only document-generator-db/
```

Or manually:

1. `00_enums.sql`
2. `document-templates.sql`
3. `document-template-versions.sql`
4. `documents.sql`
5. `document-render-logs.sql`
6. `document-callback-attempts.sql`

### Entities

- **document_templates** — `code`, `engine`, `default_format`, multi-tenant `tenant_id`
- **document_template_versions** — `content`, `schema`, `variables`, publish flag
- **documents** — async job with `request_id` idempotency, file metadata, DMS, callback, retry, signature fields
- **document_render_logs** — per worker attempt
- **document_callback_attempts** — webhook HTTP audit

### API (`openapi.yaml`)

- **Base path:** `/document-generator/v1`
- **Auth:** `Authorization: Bearer`
- **Tenant:** `X-Tenant-Id` header (multi-tenant)
- **Pagination:** `{ data, meta }` on list endpoints
- **Idempotency:** `request_id` on `POST /documents` (+ optional `Idempotency-Key` header)

Key routes:

| Method | Path | Description |
|--------|------|-------------|
| `GET/POST` | `/templates` | List / create templates |
| `GET/PATCH/DELETE` | `/templates/{template_id}` | Detail / update / deactivate |
| `GET/POST` | `/templates/{template_id}/versions` | List / create versions |
| `POST` | `/templates/.../versions/{version_id}/publish` | Publish version |
| `GET/POST` | `/documents` | List / queue generation (`202`) |
| `GET` | `/documents/by-request/{request_id}` | Lookup by idempotency key |
| `GET/DELETE` | `/documents/{document_id}` | Detail / soft-delete |
| `POST` | `/documents/{document_id}/cancel` | Cancel in-flight job |
| `POST` | `/documents/{document_id}/retry` | Retry failed job |
| `GET` | `/documents/{document_id}/download` | Signed URL redirect |
| `GET` | `/documents/{document_id}/render-logs` | Render diagnostics |
| `GET` | `/documents/{document_id}/callback-attempts` | Webhook attempts |

### Operational notes

- Validate `payload` against `document_template_versions.schema` in the service layer.
- Published version is used when `template_version` is omitted on create.
- Soft-delete documents via `deleted_at`; templates deactivate via `is_active = false`.

Import `dbdiagram.dbml` into [dbdiagram.io](https://dbdiagram.io) for visualization.
