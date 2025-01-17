-- +goose Up
-- +goose StatementBegin
CREATE TABLE games
(
    id             serial PRIMARY KEY,
    league_id      int       NOT NULL REFERENCES leagues (id),
    home_team_id   int       NOT NULL REFERENCES teams (id),
    away_team_id   int       NOT NULL REFERENCES teams (id),
    played_minutes int       NOT NULL,
    scheduled_at   timestamp NOT NULL,
    created_at     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_games_updated_at
    BEFORE UPDATE
    ON games
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


CREATE UNIQUE INDEX games_unique_index ON games (league_id, home_team_id, away_team_id, scheduled_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS games;
-- +goose StatementEnd
