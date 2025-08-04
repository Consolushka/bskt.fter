-- +goose Up
-- +goose StatementBegin
alter table game_team_player_stats add played_seconds int;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table game_team_player_stats drop column played_seconds;
-- +goose StatementEnd
