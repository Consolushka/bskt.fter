-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams
(
    id         serial PRIMARY KEY,
    name       varchar(255) NOT NULL,
    home_town  varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX teams_name_home_town_unique_idx ON teams (name, home_town);

CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE
    ON teams
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
