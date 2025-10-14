-- +goose Up
-- +goose StatementBegin
ALTER TABLE tournaments
    add column regulation_duration int default 40;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournaments
    drop column if exists regulation_duration;
-- +goose StatementEnd
