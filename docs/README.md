# Architecture Documentation — Document Generator Service

This documentation describes the **go-document-generator** project using the **C4 model** (Levels 1–4) and **sequence diagrams** for the main business flows.

## Table of Contents

| Document | Description |
|----------|-------------|
| [**C4 Overview (C1–C4 summary)**](./architecture/c4-overview.md) | Single-page summary of all levels |
| [C1 — System Context](./architecture/c1-system-context.md) | System in context of users & external systems |
| [C2 — Containers](./architecture/c2-containers.md) | Deployable processes: API, Consumer, DB, Kafka, Redis |
| [C3 — Components (API)](./architecture/c3-components-api.md) | Internal HTTP application components |
| [C3 — Components (Domain)](./architecture/c3-components-domain.md) | Use cases, repositories, state machine |
| [C4 — Code (Documents)](./architecture/c4-code-documents.md) | Package & class detail for the document domain |
| [Document Status Flow](./flows/document-status-flow.md) | State machine PENDING → GENERATED |
| [Sequence: Create Document](./sequences/01-create-document.md) | POST `/documents` |
| [Sequence: Status Transition (PATCH)](./sequences/02-document-patch-state.md) | PATCH `/documents/:id` + state machine |
| [Sequence: Generate (PROCESSING→GENERATED)](./sequences/03-generate-document.md) | Render PDF/HTML/CSV |
| [Sequence: Template Management](./sequences/04-template-management.md) | Templates & versions |
| [Sequence: Cancel & Retry](./sequences/05-cancel-retry.md) | Cancel / retry job |
| [Sequence: Kafka Events](./sequences/06-kafka-events.md) | Producers & consumers |

## Conventions

- API aligns with [`database/openapi.yaml`](../database/openapi.yaml).
- Document statuses align with [`database/00_enums.sql`](../database/00_enums.sql).
- Diagrams use **Mermaid** (renders on GitHub, VS Code, Cursor).

## Technical Summary

```
Client → HTTP API (Echo) → Usecase → Repository (GORM) → PostgreSQL
                              ↓
                         Kafka Producer → document-events / template-events
                              ↓
                         Generator (PDF/HTML/CSV) → Local Storage
```
