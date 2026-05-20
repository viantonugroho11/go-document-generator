CREATE TABLE document_templates (
    id              BIGSERIAL PRIMARY KEY,

    tenant_id       UUID,

    code            VARCHAR(100) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,

    engine          template_engine NOT NULL,
    default_format  output_format NOT NULL,

    category        VARCHAR(100),

    is_active       BOOLEAN NOT NULL DEFAULT TRUE,

    created_by      VARCHAR(100),
    updated_by      VARCHAR(100),

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uq_document_templates_tenant_code
    ON document_templates (tenant_id, code)
    WHERE tenant_id IS NOT NULL;

CREATE UNIQUE INDEX uq_document_templates_code_global
    ON document_templates (code)
    WHERE tenant_id IS NULL;

CREATE INDEX idx_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_category ON document_templates (category);
CREATE INDEX idx_document_templates_active ON document_templates (is_active);
CREATE INDEX idx_document_templates_created_at ON document_templates (created_at);
