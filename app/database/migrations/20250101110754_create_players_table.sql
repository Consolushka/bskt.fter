-- +goose Up
-- +goose StatementBegin
CREATE TABLE players
(
    id               serial PRIMARY KEY,
    league_player_id int          NOT NULL,
    full_name        varchar(255) NOT NULL,
    birth_date       date         NOT NULL,
    created_at       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE
    ON players
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE UNIQUE INDEX players_league_player_id_uindex ON players (league_player_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
