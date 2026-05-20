CREATE TABLE documents (
    id                    BIGSERIAL PRIMARY KEY,

    tenant_id             UUID,

    request_id            VARCHAR(100) NOT NULL,

    template_id           BIGINT REFERENCES document_templates (id),
    template_version_id   BIGINT REFERENCES document_template_versions (id),

    template_code         VARCHAR(100) NOT NULL,
    template_version      INTEGER NOT NULL,

    payload               JSONB NOT NULL,
    metadata              JSONB,

    status                document_status NOT NULL DEFAULT 'PENDING',

    error_message         TEXT,

    output_format         output_format NOT NULL,

    -- Generated file
    file_name             VARCHAR(255),
    file_path             TEXT,

    storage_provider      storage_provider,

    file_size             BIGINT,
    checksum              VARCHAR(64),
    content_type          VARCHAR(100),

    -- Digital signature
    is_signed             BOOLEAN NOT NULL DEFAULT FALSE,
    signature_provider    VARCHAR(100),
    signed_at             TIMESTAMP,

    -- DMS
    store_to_dms          BOOLEAN NOT NULL DEFAULT FALSE,
    dms_document_id       VARCHAR(100),
    dms_status            dms_status NOT NULL DEFAULT 'NOT_SENT',

    -- Callback / webhook
    has_callback          BOOLEAN NOT NULL DEFAULT FALSE,
    callback_url          TEXT,
    callback_status       callback_status NOT NULL DEFAULT 'PENDING',
    callback_last_at      TIMESTAMP,

    -- Retry
    retry_count           INTEGER NOT NULL DEFAULT 0,
    next_retry_at         TIMESTAMP,

    expired_at            TIMESTAMP,

    created_by            VARCHAR(100),
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at          TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP
);

CREATE UNIQUE INDEX uq_documents_tenant_request_id
    ON documents (tenant_id, request_id)
    WHERE tenant_id IS NOT NULL AND deleted_at IS NULL;

CREATE UNIQUE INDEX uq_documents_request_id_global
    ON documents (request_id)
    WHERE tenant_id IS NULL AND deleted_at IS NULL;

CREATE INDEX idx_documents_request_id ON documents (request_id);
CREATE INDEX idx_documents_template_code ON documents (template_code);
CREATE INDEX idx_documents_status ON documents (status);
CREATE INDEX idx_documents_dms_status ON documents (dms_status);
CREATE INDEX idx_documents_callback_status ON documents (callback_status);
CREATE INDEX idx_documents_created_at ON documents (created_at);
CREATE INDEX idx_documents_processed_at ON documents (processed_at);
CREATE INDEX idx_documents_next_retry ON documents (next_retry_at);
CREATE INDEX idx_documents_expired_at ON documents (expired_at);
CREATE INDEX idx_documents_deleted_at ON documents (deleted_at)
    WHERE deleted_at IS NOT NULL;
