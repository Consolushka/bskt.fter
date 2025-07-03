-- +goose Up
-- +goose StatementBegin
CREATE TABLE players
(
    id            serial PRIMARY KEY,
    full_name     varchar(255) NOT NULL,
    birth_date_at timestamp    NOT NULL,
    created_at    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX players_full_name_birth_date_at_unique_idx ON players (full_name, birth_date_at);

CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE
    ON players
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
