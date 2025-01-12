-- +goose Up
-- +goose StatementBegin
ALTER TABLE players
    ADD COLUMN league_player_id int NOT NULL;

CREATE UNIQUE INDEX players_league_player_id_uindex ON players (league_player_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players
    DROP COLUMN league_player_id;

DROP INDEX IF EXISTS players_league_player_id_uindex;
-- +goose StatementEnd
