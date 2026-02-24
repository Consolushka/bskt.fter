-- +goose Up
-- +goose StatementBegin
-- 1. Normalize existing data for api_nba provider (where values are > 1)
UPDATE game_team_player_stats gtps
SET field_goals_percentage = field_goals_percentage / 100
FROM games g
JOIN tournament_providers tp ON g.tournament_id = tp.tournament_id
WHERE gtps.game_id = g.id
  AND tp.provider_name = 'API_NBA'
  AND gtps.field_goals_percentage > 1.0;

-- 2. Add the check constraint
ALTER TABLE game_team_player_stats
    ADD CONSTRAINT check_field_goals_percentage_range
    CHECK (field_goals_percentage >= 0 AND field_goals_percentage <= 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE game_team_player_stats
    DROP CONSTRAINT IF EXISTS check_field_goals_percentage_range;

-- Note: We don't revert the data normalization (dividing by 100) because
-- it's a correction to match the new system standard.
-- +goose StatementEnd
