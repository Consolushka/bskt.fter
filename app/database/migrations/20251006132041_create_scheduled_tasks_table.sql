-- +goose Up
-- +goose StatementBegin
CREATE TABLE scheduled_tasks (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    next_execution_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_scheduled_tasks_next_execution_at ON scheduled_tasks(next_execution_at);
CREATE INDEX idx_scheduled_tasks_deleted_at ON scheduled_tasks(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled_tasks;
-- +goose StatementEnd
