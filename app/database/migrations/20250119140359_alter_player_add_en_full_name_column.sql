-- +goose Up
-- +goose StatementBegin
ALTER TABLE players
    RENAME COLUMN full_name TO full_name_local;

ALTER TABLE players
    ADD COLUMN full_name_en VARCHAR(255)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players
    RENAME COLUMN full_name_local TO full_name;
ALTER TABLE players
    DROP COLUMN full_name_en;
-- +goose StatementEnd
