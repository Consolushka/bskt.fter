-- +goose Up
-- +goose StatementBegin
CREATE TABLE game_team_player_stats
(
    id         serial PRIMARY KEY,
    game_id    int       NOT NULL REFERENCES games (id),
    team_id    int       NOT NULL REFERENCES teams (id),
    player_id  int       NOT NULL REFERENCES players (id),
    plus_minus int       NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX game_team_player_stats_game_id_team_id_player_id_unique_idx ON game_team_player_stats (game_id, team_id, player_id);

CREATE TRIGGER update_game_team_player_stats_updated_at
    BEFORE UPDATE
    ON game_team_player_stats
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS game_team_player_stats;
-- +goose StatementEnd
