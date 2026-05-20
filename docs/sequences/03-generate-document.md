# Sequence — Generate Document (PROCESSING → GENERATED)

Executed inside the `OnToGenerated` transition when PATCH sets `status: GENERATED`.

## Diagram

```mermaid
sequenceDiagram
    autonumber
    participant Tr as transitions.toGenerated
    participant TplRepo as TemplatesRepository
    participant VerRepo as VersionsRepository
    participant Sel as GeneratorSelector<br/>(factory.go)
    participant Gen as Generator<br/>(PDF|HTML|CSV)
    participant Store as storage.SaveDocument
    participant Doc as Document entity

    Tr->>TplRepo: GetByID(template_id)
    TplRepo-->>Tr: engine (HANDLEBARS|HTML|...)
    Tr->>VerRepo: GetByID(template_version_id)
    VerRepo-->>Tr: content (template string)

    Tr->>Sel: Select(output_format, engine)
    Sel-->>Tr: Generator instance

    Tr->>Gen: Generate(ctx, content, payload)
    Note over Gen: PDF: html/template → wkhtmltopdf<br/>HTML: html/template<br/>CSV: text/template + csv helpers
    Gen-->>Tr: bytes, content_type

    Tr->>Store: SaveDocument(id, request_id, ext, bytes)
    Store-->>Tr: file_path, file_name

    Tr->>Doc: Set GENERATED, checksum, file_size,<br/>processed_at, storage_provider=LOCAL
    Tr-->>Tr: return Document

    alt Generate error
        Tr->>Doc: status=FAILED, error_message
        Tr-->>Tr: return error
    end
```

## Generator Selection

```mermaid
flowchart TD
    A[output_format] --> B{format?}
    B -->|PDF| C[pdf.WKHTMLToPDFGenerator]
    B -->|HTML| D[html.Generator]
    B -->|DOCX| D
    B -->|other| E{engine?}
    E -->|HTML| D
    E -->|default| F[csv.TmplCSVGenerator]
```

## Output File

Path: `./storage/documents/{document_id}/{request_id}.{ext}`

Download: `GET /documents/:id/download` → redirect to `file_path`.
