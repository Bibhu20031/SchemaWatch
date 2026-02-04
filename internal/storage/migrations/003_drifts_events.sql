-- +goose Up
CREATE TABLE drift_events (
    id BIGSERIAL PRIMARY KEY,
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON DELETE CASCADE,
    version_from INT NOT NULL,
    version_to INT NOT NULL,
    change_type TEXT NOT NULL,
    column_name TEXT NOT NULL,
    impact TEXT NOT NULL,
    before_value JSONB,
    after_value JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_drift_schema_id ON drift_events(schema_id);
CREATE INDEX idx_drift_created_at ON drift_events(created_at);

-- +goose Down
DROP TABLE IF EXISTS drift_events;
