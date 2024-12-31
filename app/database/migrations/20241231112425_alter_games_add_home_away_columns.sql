-- +goose Up
-- +goose StatementBegin
ALTER TABLE games
    ADD COLUMN home_team_id int NOT NULL REFERENCES teams (id),
    ADD COLUMN away_team_id int NOT NULL REFERENCES teams (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE games
    DROP COLUMN home_team_id,
    DROP COLUMN away_team_id;
-- +goose StatementEnd
