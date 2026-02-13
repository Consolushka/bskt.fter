-- +goose Up
-- +goose StatementBegin
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

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS poll_watermarks;
-- +goose StatementEnd
