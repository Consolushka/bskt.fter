-- +goose Up
-- +goose StatementBegin
CREATE TABLE games (
    id serial PRIMARY KEY,
    league_id int NOT NULL REFERENCES leagues(id),
    played_minutes int NOT NULL,
    scheduled_at timestamp NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_games_updated_at
    BEFORE UPDATE
    ON games
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_games_updated_at ON games;

DROP TABLE IF EXISTS games;
-- +goose StatementEnd
