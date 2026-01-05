-- +goose Up
-- +goose StatementBegin
ALTER TABLE game_team_player_stats
    ADD COLUMN points                 smallint NOT NULL DEFAULT 0,
    ADD COLUMN assists                smallint NOT NULL DEFAULT 0,
    ADD COLUMN rebounds               smallint NOT NULL DEFAULT 0,
    ADD COLUMN steals                 smallint NOT NULL DEFAULT 0,
    ADD COLUMN blocks                 smallint NOT NULL DEFAULT 0,
    ADD COLUMN field_goals_percentage float4   NOT NULL DEFAULT 0,
    ADD COLUMN turnovers              smallint NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE game_team_player_stats
    DROP COLUMN IF EXISTS points,
    DROP COLUMN IF EXISTS assists,
    DROP COLUMN IF EXISTS rebounds,
    DROP COLUMN IF EXISTS steals,
    DROP COLUMN IF EXISTS blocks,
    DROP COLUMN IF EXISTS field_goals_percentage,
    DROP COLUMN IF EXISTS turnovers;
-- +goose StatementEnd