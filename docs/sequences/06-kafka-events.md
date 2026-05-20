# Sequence — Kafka Events

Event-driven integration for downstream notifications and example consumers.

## Topics (config)

| Topic | Producer | Payload |
|-------|----------|---------|
| `user-events` | User create/update | `UserCreatedEvent` |
| `template-events` | Template create/update | `TemplateCreatedEvent` |
| `template-version-events` | Version create/publish | `TemplateVersionCreatedEvent` |
| `document-events` | Document queued/retried | `DocumentQueuedEvent`, `DocumentRetriedEvent` |

## Publish — Document Queued

```mermaid
sequenceDiagram
    autonumber
    participant UC as documents.Service.Create
    participant Pub as DocumentEventPublisherKafka
    participant K as Kafka Producer
    participant Topic as document-events

    UC->>Pub: PublishDocumentQueued(document)
    Pub->>K: Publish(DocumentQueuedEvent)
    Note over K: key = request_id
    K->>Topic: message
```

## Publish — Template / Version

```mermaid
sequenceDiagram
    participant UC_T as documenttemplates.Service
    participant UC_V as documenttemplateversions.Service
    participant K1 as template-events
    participant K2 as template-version-events

    UC_T->>K1: TemplateCreatedEvent
    UC_V->>K2: TemplateVersionCreatedEvent
```

## Consumer (examples)

```mermaid
sequenceDiagram
    autonumber
    actor Ops
    participant Main as cmd/consumer
    participant Boot as bootstrap.RunConsumer
    participant KCons as kafka.Consumer
    participant H as UserCreatedHandler / OrderCreatedHandler

    Ops->>Main: go run cmd/consumer -consumer=user
    Main->>Boot: RunConsumer("user")
    Boot->>KCons: Subscribe(user-events, group_id)
    loop Poll messages
        KCons->>H: Handle(UserCreatedEvent)
        H-->>KCons: Progress Success/Error
    end
```

> A dedicated document consumer (process QUEUED → PROCESSING → GENERATED) can be added as a separate handler calling `documents.Service` — not wired in the repo yet; the current flow triggers the state machine via the **PATCH API**.

## End-to-End Diagram (target architecture)

```mermaid
sequenceDiagram
    participant C as Client
    participant API as HTTP API
    participant DB as PostgreSQL
    participant K as Kafka
    participant W as Worker Consumer

    C->>API: POST /documents
    API->>DB: INSERT QUEUED
    API->>K: document-queued
    API-->>C: 202

    K->>W: consume event
    W->>API: PATCH PROCESSING (internal)
    W->>API: PATCH GENERATED (internal)
    Note over W: or call usecase directly
    W->>DB: UPDATE + file
```
