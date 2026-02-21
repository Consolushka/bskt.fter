-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS poll_watermarks;

CREATE TABLE poll_watermarks (
    tournament_id BIGINT PRIMARY KEY,
    last_successful_poll_at timestamptz NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_poll_watermarks_tournament FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE
);

CREATE TRIGGER update_poll_watermarks_updated_at
    BEFORE UPDATE
    ON poll_watermarks
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS poll_watermarks;

CREATE TABLE poll_watermarks (
    task_type varchar(255) PRIMARY KEY,
    last_successful_poll_at timestamptz NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_poll_watermarks_updated_at
    BEFORE UPDATE
    ON poll_watermarks
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd
