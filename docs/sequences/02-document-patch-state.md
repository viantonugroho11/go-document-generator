# Sequence — PATCH Document & State Machine

Partial field updates and/or sequential status transitions.

## Diagram

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as DocumentHandler.Patch
    participant UC as documents.Service
    participant DocRepo as DocumentsRepository
    participant SMF as StateMachineFactory
    participant St as Current State<br/>(pending|queued|processing|...)
    participant Tr as Transition Handler

    Client->>API: PATCH /documents/:id<br/>{status?, payload?, ...}
    API->>DocRepo: GetByID (via GetByID usecase)
    API->>API: ApplyPatchDocument(existing, req)

    API->>UC: Patch(merged Document)
    UC->>DocRepo: GetByID
    UC->>UC: mergeDocumentPatch

    UC->>SMF: NewStateMachine(existing)
    SMF-->>UC: state machine (current = existing.status)

    UC->>St: Do(ctx, update)
    Note over St: Switch by update.Status

    alt status = QUEUED (from PENDING)
        St->>Tr: OnToQueued
    else status = PROCESSING (from QUEUED)
        St->>Tr: OnToProcessing
    else status = GENERATED (from PROCESSING)
        St->>Tr: OnToGenerated
        Note over Tr: See sequence 03-generate
    else status = CANCELLED
        St->>Tr: OnToCancelled
    else same status (fields only)
        St->>Tr: OnFieldUpdate
        Tr->>Tr: ValidateSchema if payload changed
    else invalid transition
        St-->>UC: ErrInvalidState
    end

    Tr-->>St: updated Document
    St-->>UC: result
    UC->>DocRepo: Update(result)
    UC-->>API: Document
    API-->>Client: 200 JSON
```

## Validation

- Status transitions must be **sequential** (see [document-status-flow.md](../flows/document-status-flow.md)).
- Field patches are rejected when status is not PENDING/QUEUED (unless a status transition is also requested).
