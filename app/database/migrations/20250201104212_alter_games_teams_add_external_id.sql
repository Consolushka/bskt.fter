-- +goose Up
-- +goose StatementBegin
ALTER TABLE games
    ADD COLUMN official_id VARCHAR(255) NULL;

COMMENT ON COLUMN games.official_id IS 'ID from official API';

ALTER TABLE teams
    ADD COLUMN official_id VARCHAR(255) NULL;
COMMENT ON COLUMN teams.official_id IS 'ID from official API';

ALTER TABLE players
 RENAME COLUMN league_player_id TO official_id;

COMMENT ON COLUMN players.official_id IS 'ID from official API';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
