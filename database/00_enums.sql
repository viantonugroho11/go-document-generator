-- Document Generator — PostgreSQL enum types (run first)

CREATE TYPE template_engine AS ENUM ('HANDLEBARS', 'MUSTACHE', 'HTML');

CREATE TYPE output_format AS ENUM ('PDF', 'HTML', 'DOCX');

CREATE TYPE document_status AS ENUM (
    'PENDING',
    'QUEUED',
    'PROCESSING',
    'GENERATED',
    'FAILED',
    'CANCELLED'
);

CREATE TYPE dms_status AS ENUM ('NOT_SENT', 'QUEUED', 'SENT', 'FAILED');

CREATE TYPE callback_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED', 'RETRYING');

CREATE TYPE storage_provider AS ENUM ('LOCAL', 'S3', 'MINIO', 'GCS', 'AZURE');
