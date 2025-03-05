-- +goose Up
-- +goose StatementBegin
CREATE TABLE players_metrics
(
    player_id int NOT NULL primary key REFERENCES players(id),
    avg_clear_imp numeric(12,11),
    avg_played_seconds numeric(6,2),
    played_games_count int,
    from_bench_count int,
    from_start_count int,
    updated_at timestamp
);

CREATE TRIGGER update_players_metrics_updated_at
    BEFORE UPDATE
    ON players_metrics
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players_metrics;
-- +goose StatementEnd
