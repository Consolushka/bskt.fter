-- +goose Up
-- +goose StatementBegin
alter table scheduled_tasks
    alter column next_execution_at type timestamptz using next_execution_at::timestamptz;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table scheduled_tasks
    alter column next_execution_at type timestamp using next_execution_at::timestamp;
-- +goose StatementEnd
