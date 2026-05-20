# Sequence — Create Document (POST /documents)

Creates a document generation job with `request_id` as the idempotency key.

## Diagram

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as Echo Handler
    participant UC as documents.Service
    participant TplRepo as TemplatesRepository
    participant VerRepo as VersionsRepository
    participant Val as validators.ValidateSchema
    participant DocRepo as DocumentsRepository
    participant Tx as BeginRepository
    participant Kafka as DocumentEventPublisher

    Client->>API: POST /documents<br/>{request_id, template_code, payload}
    API->>API: Bind DTO + X-Tenant-Id
    API->>UC: Create(CreateInput)

    UC->>DocRepo: GetByRequestId(request_id)
    alt Already exists (idempotent)
        DocRepo-->>UC: existing document
        UC-->>API: 200 + existing
    else New job
        UC->>TplRepo: GetByCode(template_code)
        TplRepo-->>UC: template (is_active)
        UC->>VerRepo: GetLatestPublished / GetByVersion
        VerRepo-->>UC: version + schema
        UC->>Val: ValidateSchema(schema, payload)
        Val-->>UC: OK

        UC->>Tx: Begin()
        UC->>DocRepo: Create(status=QUEUED)
        DocRepo-->>UC: created
        UC->>Tx: Commit()

        UC->>Kafka: PublishDocumentQueued
        UC-->>API: 202 Accepted
        API-->>Client: GeneratedDocument JSON
    end
```

## HTTP Responses

| Condition | Status |
|-----------|--------|
| New job | `202 Accepted` |
| `request_id` replay | `200 OK` |
| Template not found / not published | `404` |
| Invalid payload schema | `400` |

## Side Effects

- Row in `documents` table with status **QUEUED**
- Kafka event on `document-events` (key = `request_id`)
