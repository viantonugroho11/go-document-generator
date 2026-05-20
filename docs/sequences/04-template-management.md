# Sequence — Template Management & Versions

Master template and versioning flow before documents can be generated.

## 4.1 Create Template

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as TemplateHandler
    participant UC as documenttemplates.Service
    participant Repo as DocumentTemplatesRepository
    participant Kafka as TemplateEventPublisher

    Client->>API: POST /templates
    API->>UC: Create(template)
    UC->>Repo: Create
    Repo-->>UC: created
    UC->>Kafka: PublishTemplateCreated
    UC-->>API: template
    API-->>Client: 201 Created
```

## 4.2 Create & Publish Version

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant API as TemplateVersionHandler
    participant UC as documenttemplateversions.Service
    participant TplRepo as TemplatesRepository
    participant VerRepo as VersionsRepository
    participant Tx as BeginRepository

    Client->>API: POST /templates/:id/versions<br/>{content, schema, output_format}
    API->>UC: Create(template_id, version)
    UC->>TplRepo: GetByID (exists)
    UC->>VerRepo: NextVersionNumber
    UC->>VerRepo: Create (SHA256 checksum)
    UC-->>API: 201 version

    Client->>API: POST /templates/:id/versions/:vid/publish
    API->>UC: Publish
    UC->>Tx: Begin
    UC->>VerRepo: UnpublishOthers
    UC->>VerRepo: Publish(is_published=true)
    UC->>Tx: Commit
    UC->>Kafka: PublishVersionPublished
    API-->>Client: 200 version
```

## Prerequisites for Document Generation

```mermaid
flowchart LR
    T[Active template] --> V[Published version]
    V --> D[POST /documents]
    D --> OK[Resolve template_code + latest published]
```

`Create` uses the latest **published** version when `template_version` is omitted.
