# Document Status Flow (State Machine)

The state machine follows the pattern in `internal/usecase/sample/states` and the PostgreSQL enum `document_status`.

## State Diagram

```mermaid
stateDiagram-v2
    [*] --> QUEUED: POST /documents\n(initial status QUEUED)

    PENDING --> QUEUED: PATCH status=QUEUED
    QUEUED --> PROCESSING: PATCH status=PROCESSING
    PROCESSING --> GENERATED: PATCH status=GENERATED\n(render file)
    PROCESSING --> FAILED: render error / PATCH status=FAILED

    PENDING --> CANCELLED: PATCH / cancel
    QUEUED --> CANCELLED: PATCH / cancel
    PROCESSING --> CANCELLED: PATCH / cancel

    FAILED --> QUEUED: POST retry / PATCH

    GENERATED --> [*]
    CANCELLED --> [*]
    FAILED --> [*]: terminal unless retry
```

## Transition Rules

| From | To | How | Handler |
|------|-----|------|---------|
| PENDING | QUEUED | PATCH | `OnToQueued` |
| QUEUED | PROCESSING | PATCH | `OnToProcessing` |
| PROCESSING | GENERATED | PATCH | `OnToGenerated` (+ generate) |
| PROCESSING | FAILED | generate error / PATCH | `OnToFailed` |
| *active* | CANCELLED | PATCH / POST cancel | `OnToCancelled` |
| FAILED | QUEUED | POST retry | `OnRetry` |
| PENDING/QUEUED | (fields) | PATCH without status change | `OnFieldUpdate` |

**Field patches** (`payload`, `metadata`, `callback_url`, …) are only allowed when status is **PENDING** or **QUEUED**.

## Typical Operational Flow

```mermaid
flowchart TD
    A[Client POST /documents] --> B[Status QUEUED + Kafka document-queued]
    B --> C[Worker/Client PATCH → PROCESSING]
    C --> D[PATCH → GENERATED]
    D --> E{Generate OK?}
    E -->|Yes| F[Save file + metadata]
    E -->|No| G[Status FAILED]
    F --> H[GET /documents/:id/download]
    F --> I[Optional webhook callback]
```

## Terminal Statuses

| Status | PATCH behavior |
|--------|----------------|
| GENERATED | No-op only; no outbound transitions |
| CANCELLED | Terminal |
| FAILED | Retry only → QUEUED |
