-- +goose Up
-- +goose StatementBegin
ALTER TABLE player_game_stats
    ADD COLUMN imp_clear numeric(12, 11);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player_game_stats
    DROP COLUMN imp_clear;
-- +goose StatementEnd
