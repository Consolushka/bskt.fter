-- +goose Up
-- +goose StatementBegin
alter table game_team_stats
    add column final_differential int default 0;

-- Обновляем колонку final_differential для всех существующих записей
UPDATE game_team_stats
SET final_differential = (
    SELECT gts1.score - gts2.score
    FROM game_team_stats gts1, game_team_stats gts2
    WHERE gts1.game_id = game_team_stats.game_id
      AND gts2.game_id = game_team_stats.game_id
      AND gts1.team_id = game_team_stats.team_id
      AND gts2.team_id != game_team_stats.team_id
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game_team_stats
    drop column if exists final_differential;
-- +goose StatementEnd