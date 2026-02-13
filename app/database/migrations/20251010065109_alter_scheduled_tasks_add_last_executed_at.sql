-- +goose Up
-- +goose StatementBegin
alter table scheduled_tasks
    add column last_executed_at timestamptz default (now() - interval '1 day');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table scheduled_tasks
    drop column last_executed_at;
-- +goose StatementEnd
