-- +goose Up
-- +goose StatementBegin
ALTER TABLE players
    ADD COLUMN draft_year int,
    ALTER COLUMN birth_date DROP NOT NULL;

CREATE UNIQUE INDEX players_unique_name_draft_birthdate_index ON players (full_name, draft_year, birth_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS players_unique_name_draft_birthdate_index;

ALTER TABLE players
    DROP COLUMN draft_year,
    ALTER COLUMN birth_date SET NOT NULL;
-- +goose StatementEnd
