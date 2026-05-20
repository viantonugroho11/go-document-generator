# C3 — Component Diagram (Domain & State Machine)

Focus on the **documents** domain and the state machine pattern (similar to `usecase/sample`).

## State Machine Diagram

```mermaid
flowchart LR
    subgraph states["documents/states"]
        P[PENDING]
        Q[QUEUED]
        PR[PROCESSING]
        G[GENERATED]
        F[FAILED]
        C[CANCELLED]
    end

    P -->|OnToQueued| Q
    Q -->|OnToProcessing| PR
    PR -->|OnToGenerated| G
    PR -->|OnToFailed| F
    P & Q & PR -->|OnToCancelled| C
    F -->|OnRetry| Q

    subgraph transitions["documents/transitions"]
        T1[field_update]
        T2[to_generated + Generator]
        T3[to_cancelled]
    end

    PR -.-> T2
    Q -.-> T1
```

## Domain Component Diagram

```mermaid
flowchart TB
    ucDoc[documents.Service]

    subgraph sm["State Machine"]
        smFactory[StateMachineFactory]
        smStates[State Impl]
        trans[Transitions]
    end

    gen[GeneratorSelector]
    valid[JSON Schema Validator]
    store[Local Storage]

    repoDoc[DocumentsRepository]
    repoTpl[TemplatesRepository]
    repoVer[VersionsRepository]

    ucDoc --> smFactory --> smStates --> trans
    trans --> gen & valid & store
    ucDoc --> repoDoc & repoTpl & repoVer
```

## Packages & Responsibilities

| Package | Pattern | Purpose |
|---------|---------|---------|
| `entity/` | Domain model | Pure structs + enums |
| `repository/` | Repository | Interface + `postgres/` + `model/` |
| `usecase/*/events.go` | Publisher interface | Kafka per domain (separate files) |
| `usecase/documents/states/` | State Machine | Factory + state per status |
| `usecase/documents/transitions/` | Command handlers | Side effects per transition |
| `usecase/documents/statemachine_wire.go` | Composition root | Wire handlers (avoids import cycles) |
| `infrastructure/documents/` | Adapter | `Selector` → PDF/HTML/CSV generator |

## Entity & Tables

```mermaid
erDiagram
    document_templates ||--o{ document_template_versions : has
    document_templates ||--o{ documents : references
    document_template_versions ||--o{ documents : uses
    documents ||--o{ document_render_logs : logs
    documents ||--o{ document_callback_attempts : callbacks

    document_templates {
        bigint id PK
        uuid tenant_id
        varchar code
        template_engine engine
    }
    document_template_versions {
        bigint id PK
        bigint template_id FK
        int version
        jsonb schema
        text content
    }
    documents {
        bigint id PK
        varchar request_id
        document_status status
        jsonb payload
    }
```
