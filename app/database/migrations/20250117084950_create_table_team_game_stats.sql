-- +goose Up
-- +goose StatementBegin
CREATE TABLE team_game_stats
(
    id         serial PRIMARY KEY,
    team_id    int REFERENCES teams (id),
    game_id    int REFERENCES games (id),
    points     int       NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_team_game_stats_updated_at
    BEFORE UPDATE
    ON team_game_stats
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_game_stats;
-- +goose StatementEnd
