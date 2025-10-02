-- +goose Up
-- +goose StatementBegin
ALTER TABLE tournaments
    add column start_at timestamp default now(),
    add column end_at timestamp default now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table tournaments
    drop column if exists start_at,
    drop column if exists end_at;
-- +goose StatementEnd