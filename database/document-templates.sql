CREATE TABLE document_templates (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(100) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,

    engine          VARCHAR(50) NOT NULL,   -- handlebars | mustache | html
    output_format   VARCHAR(20) NOT NULL,   -- PDF | HTML | DOCX

    is_active       BOOLEAN DEFAULT TRUE,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_document_templates_code ON document_templates(code);
CREATE INDEX idx_document_templates_is_active ON document_templates(is_active);
CREATE INDEX idx_document_templates_created_at ON document_templates(created_at);
CREATE INDEX idx_document_templates_updated_at ON document_templates(updated_at);