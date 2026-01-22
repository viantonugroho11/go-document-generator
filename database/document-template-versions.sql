CREATE TABLE document_template_versions (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL REFERENCES document_templates(id),

    version         INTEGER NOT NULL,
    content         TEXT NOT NULL,

    schema          JSONB,       -- input validation
    sample_payload  JSONB,

    is_published    BOOLEAN DEFAULT FALSE,
    published_at    TIMESTAMP,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (template_id, version)
);

CREATE INDEX idx_template_versions_template_id ON document_template_versions(template_id);
CREATE INDEX idx_template_versions_version ON document_template_versions(version);
CREATE INDEX idx_template_versions_published ON document_template_versions(is_published);
CREATE INDEX idx_template_versions_created_at ON document_template_versions(created_at);