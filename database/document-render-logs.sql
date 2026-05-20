CREATE TABLE document_render_logs (
    id                  BIGSERIAL PRIMARY KEY,

    document_id         BIGINT NOT NULL REFERENCES documents (id) ON DELETE CASCADE,

    status              document_status NOT NULL,

    message             TEXT,
    execution_time_ms   BIGINT,
    stack_trace         TEXT,

    worker_name         VARCHAR(100),

    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_render_logs_document_id ON document_render_logs (document_id);
CREATE INDEX idx_render_logs_status ON document_render_logs (status);
CREATE INDEX idx_render_logs_created_at ON document_render_logs (created_at);
