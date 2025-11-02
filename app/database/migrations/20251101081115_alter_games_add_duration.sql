-- +goose Up
-- +goose StatementBegin
alter table games
    add duration int;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table games
    drop column if exists duration;
-- +goose StatementEnd
