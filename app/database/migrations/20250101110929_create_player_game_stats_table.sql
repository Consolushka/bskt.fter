-- +goose Up
-- +goose StatementBegin
CREATE TABLE player_game_stats
(
    player_id      int REFERENCES players (id),
    game_id        int REFERENCES games (id),
    team_id        int REFERENCES teams,
    pls_min        int,
    played_seconds int,
    is_bench       bool,
    created_at     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (player_id, game_id)
);

CREATE TRIGGER update_player_game_stats_updated_at
    BEFORE UPDATE
    ON player_game_stats
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS player_game_stats;
-- +goose StatementEnd
