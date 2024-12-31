-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams
(
    id         serial PRIMARY KEY,
    alias      varchar(15),
    league_id  int          NOT NULL REFERENCES leagues (id),
    name       varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX team_alias_league_unique_index ON teams (alias, league_id);

CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE
    ON teams
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;

DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
