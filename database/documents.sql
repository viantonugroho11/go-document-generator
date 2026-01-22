CREATE TABLE documents (
    id              BIGSERIAL PRIMARY KEY,
    request_id      VARCHAR(100) NOT NULL UNIQUE, -- idempotency key

    -- Template
    template_code   VARCHAR(100) NOT NULL,
    template_version INTEGER,

    -- Input
    payload         JSONB NOT NULL,
    metadata        JSONB,

    -- Rendering lifecycle
    status          document_status NOT NULL DEFAULT 'PENDING',
    error_message   TEXT,

    -- Output file
    file_name       VARCHAR(255),
    file_path       TEXT,
    file_size       BIGINT,
    checksum        VARCHAR(64),
    content_type    VARCHAR(50),

    -- DMS integration
    store_to_dms    BOOLEAN DEFAULT FALSE,
    dms_document_id VARCHAR(100),
    dms_status      dms_status DEFAULT 'NOT_SENT',

    -- Callback
    has_callback    BOOLEAN DEFAULT FALSE,
    callback_url    TEXT,
    callback_status callback_status DEFAULT 'PENDING',

    -- Audit
    created_by      VARCHAR(100),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_documents_request_id ON documents(request_id);
CREATE INDEX idx_documents_template_code ON documents(template_code);
CREATE INDEX idx_documents_status ON documents(status);
CREATE INDEX idx_documents_created_at ON documents(created_at);
CREATE INDEX idx_documents_processed_at ON documents(processed_at);
CREATE INDEX idx_documents_updated_at ON documents(updated_at);