-- +goose Up
CREATE TABLE webhooks (
    id BIGSERIAL PRIMARY KEY,
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE webhooks;