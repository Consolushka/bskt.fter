-- +goose Up
-- +goose StatementBegin
alter table game_team_stats
    add column is_home boolean not null default false;

-- До этой миграции хозяином считалась строка, вставленная первой в игре (минимальный id) —
-- так её и сейчас показывает API. Переносим это правило в явный признак.
update game_team_stats g
set is_home = (g.id = (select min(id) from game_team_stats where game_id = g.game_id));

-- Без значения по умолчанию запись, забывшая указать сторону, упадёт сразу,
-- а не сохранит хозяина как гостя.
alter table game_team_stats
    alter column is_home drop default;

create unique index game_team_stats_game_id_home_unique_idx on game_team_stats (game_id) where is_home;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists game_team_stats_game_id_home_unique_idx;

alter table game_team_stats
    drop column if exists is_home;
-- +goose StatementEnd
