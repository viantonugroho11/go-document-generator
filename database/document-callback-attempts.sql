CREATE TABLE document_callback_attempts (
    id                      BIGSERIAL PRIMARY KEY,

    document_id             BIGINT NOT NULL REFERENCES documents (id) ON DELETE CASCADE,

    callback_url            TEXT NOT NULL,

    request_payload         JSONB,
    response_payload        JSONB,

    response_status_code    INTEGER,

    is_success              BOOLEAN NOT NULL DEFAULT FALSE,
    error_message           TEXT,

    attempted_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_callback_attempts_document_id
    ON document_callback_attempts (document_id);

CREATE INDEX idx_callback_attempts_attempted_at
    ON document_callback_attempts (attempted_at);
