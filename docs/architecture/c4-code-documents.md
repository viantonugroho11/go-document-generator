# C4 — Code Diagram (Documents Usecase)

Code level: main interfaces/classes in the `documents` package and its dependencies.

## Class Diagram (Documents)

```mermaid
classDiagram
    class Service {
        <<interface>>
        +Create(ctx, CreateInput) Document, bool, error
        +Patch(ctx, Document) Document, error
        +Cancel(ctx, id, tenantID) Document, error
        +Retry(ctx, id, tenantID) Document, error
    }

    class service {
        -docs DocumentsRepository
        -templates DocumentTemplatesRepository
        -versions DocumentTemplateVersionsRepository
        -stateMachine IDocumentStateMachineFactory
        -selector GeneratorSelector
        -publisher DocumentEventPublisher
        +applyStateMachine()
        +transitionDocument()
    }

    class IDocumentStateMachineFactory {
        <<interface>>
        +NewStateMachine(current) IDocumentStateMachine
    }

    class IDocumentState {
        <<interface>>
        +Do(ctx, update) Document, error
    }

    class IOnStateTransition {
        <<interface>>
        +OnStateTransition(ctx, update) Document, error
    }

    class GeneratorSelector {
        <<interface>>
        +Select(outputFormat, engine) Generator
    }

    class Generator {
        <<interface>>
        +Generate(ctx, template, data) bytes, contentType, error
    }

    class toGenerated {
        -deps Deps
        +OnStateTransition() render + save file
    }

    Service <|.. service
    service --> IDocumentStateMachineFactory
    service --> GeneratorSelector
    IDocumentStateMachineFactory ..> IDocumentState
    IDocumentState ..> IOnStateTransition
    IOnStateTransition <|.. toGenerated
    toGenerated --> GeneratorSelector
```

## File Map

```
internal/usecase/documents/
├── service.go              # Service interface + Create, List, Cancel, Retry
├── patch.go                # applyStateMachine, transitionDocument
├── routing.go              # Generator, GeneratorSelector interfaces
├── events.go               # DocumentEventPublisher
├── statemachine_wire.go    # BuildStateHandlers (wiring)
├── selector_adapter.go     # Adapter GeneratorSelector → transitions
├── states/
│   ├── state.go            # Factory, Handlers, interfaces
│   ├── pending.go          # → QUEUED | CANCELLED | field update
│   ├── queued.go           # → PROCESSING | CANCELLED | field update
│   ├── processing.go       # → GENERATED | FAILED | CANCELLED
│   ├── generated.go        # terminal
│   ├── failed.go           # → QUEUED (retry)
│   └── cancelled.go        # terminal
└── transitions/
    ├── deps.go
    ├── field_update.go     # JSON Schema validation
    ├── to_generated.go     # generateAndFinalize()
    ├── to_queued.go
    ├── to_processing.go
    ├── to_cancelled.go
    ├── to_failed.go
    └── retry.go
```

## Internal Sequence: Patch + State Machine

```mermaid
sequenceDiagram
    participant S as service.Patch
    participant R as DocumentsRepository
    participant F as StateMachineFactory
    participant St as current State
    participant T as Transition Handler

    S->>R: GetByID
    S->>S: mergeDocumentPatch
    S->>F: NewStateMachine(existing)
    F->>St: Do(ctx, update)
    St->>T: OnStateTransition (by target status)
    T-->>St: updated Document
    St-->>S: result
    S->>R: Update(result)
```

## Generator Infrastructure

| Output | Engine | Implementation |
|--------|--------|----------------|
| PDF | HTML template | `infrastructure/documents/pdf` (wkhtmltopdf) |
| HTML | HTML | `infrastructure/documents/html` |
| CSV / default | HANDLEBARS/MUSTACHE | `infrastructure/documents/csv` (text/template) |

Selector: `infrastructure/documents/factory.go` → `usecase/documents.GeneratorSelector`.
