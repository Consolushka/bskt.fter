-- +goose Up
-- +goose StatementBegin
CREATE TABLE games
(
    id         serial PRIMARY KEY,
    scheduled_at timestamp    NOT NULL,
    tournament_id int          NOT NULL,
    title      varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX games_scheduled_at_tournament_id_title_unique_idx ON games (scheduled_at, tournament_id, title);

CREATE TRIGGER update_games_updated_at
    BEFORE UPDATE
    ON games
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS games;
-- +goose StatementEnd
