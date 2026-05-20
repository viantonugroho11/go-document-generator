CREATE TABLE document_template_versions (
    id              BIGSERIAL PRIMARY KEY,

    tenant_id       UUID,

    template_id     BIGINT NOT NULL REFERENCES document_templates (id) ON DELETE CASCADE,

    version         INTEGER NOT NULL,
    content         TEXT NOT NULL,

    schema          JSONB,
    variables       JSONB,
    sample_payload  JSONB,

    output_format   output_format NOT NULL,

    checksum        VARCHAR(64),

    is_published    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at    TIMESTAMP,

    created_by      VARCHAR(100),

    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT document_template_versions_template_version_uq
        UNIQUE (template_id, version)
);

CREATE INDEX idx_template_versions_template_id
    ON document_template_versions (template_id);

CREATE INDEX idx_template_versions_version
    ON document_template_versions (version);

CREATE INDEX idx_template_versions_published
    ON document_template_versions (is_published);

CREATE INDEX idx_template_versions_created_at
    ON document_template_versions (created_at);
