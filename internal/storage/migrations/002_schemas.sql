-- +goose Up
CREATE TABLE schemas (
    id BIGSERIAL PRIMARY KEY,
    db_host TEXT NOT NULL,
    db_port INT NOT NULL,
    db_name TEXT NOT NULL,
    schema_name TEXT NOT NULL,
    table_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE schema_versions (
    id BIGSERIAL PRIMARY KEY,
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON DELETE CASCADE,
    version INT NOT NULL,
    snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(schema_id, version)
);

-- +goose Down
DROP TABLE IF EXISTS schema_versions;
DROP TABLE IF EXISTS schemas;
