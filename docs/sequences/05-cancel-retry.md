# Sequence — Cancel & Retry Document

## 5.1 Cancel (POST /documents/:id/cancel)

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as DocumentHandler.Cancel
    participant UC as documents.Service
    participant SM as State Machine
    participant Tr as OnToCancelled
    participant Repo as DocumentsRepository

    Client->>API: POST /documents/:id/cancel
    API->>UC: Cancel(id, tenant)
    UC->>Repo: GetByID
    UC->>SM: transitionDocument(CANCELLED)
    SM->>Tr: OnToCancelled
    Tr-->>SM: status=CANCELLED
    UC->>Repo: Update
    API-->>Client: 200 Document
```

Allowed only from: **PENDING**, **QUEUED**, **PROCESSING**.

## 5.2 Retry (POST /documents/:id/retry)

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as DocumentHandler.Retry
    participant UC as documents.Service
    participant SM as State Machine
    participant Tr as OnRetry
    participant Repo as DocumentsRepository
    participant Kafka as DocumentEventPublisher

    Client->>API: POST /documents/:id/retry
    API->>UC: Retry(id, tenant)
    UC->>Repo: GetByID (status=FAILED)
    UC->>SM: transitionDocument(QUEUED)
    SM->>Tr: OnRetry
    Note over Tr: retry_count++, clear error_message
    UC->>Repo: Update
    UC->>Kafka: PublishDocumentRetried
    API-->>Client: 202 Document
```

## Comparison

| Action | Initial status | Final status | Kafka |
|--------|----------------|--------------|-------|
| Cancel | PENDING/QUEUED/PROCESSING | CANCELLED | — |
| Retry | FAILED | QUEUED | document-retried |
