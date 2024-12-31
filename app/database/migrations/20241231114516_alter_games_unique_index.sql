-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX games_unique_index ON games (league_id, home_team_id, away_team_id, scheduled_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX games_unique_index;
-- +goose StatementEnd
