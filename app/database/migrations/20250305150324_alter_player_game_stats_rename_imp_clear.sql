-- +goose Up
-- +goose StatementBegin

ALTER TABLE player_game_stats
    RENAME COLUMN imp_clear TO imp_clean;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player_game_stats
    RENAME COLUMN imp_clean TO imp_clear;
-- +goose StatementEnd
