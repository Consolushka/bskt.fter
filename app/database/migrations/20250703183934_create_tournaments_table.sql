-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournaments
(
    id         serial PRIMARY KEY,
    league_id  int          NOT NULL REFERENCES leagues (id),
    name       varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX tournaments_league_id_name_unique_idx ON tournaments (league_id, name);

CREATE TRIGGER update_tournaments_updated_at
    BEFORE UPDATE
    ON tournaments
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tournaments;
-- +goose StatementEnd
