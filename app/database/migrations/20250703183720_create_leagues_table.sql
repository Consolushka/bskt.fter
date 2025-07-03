-- +goose Up
-- +goose StatementBegin
CREATE TABLE leagues
(
    id         serial PRIMARY KEY,
    name       varchar(255) NOT NULL,
    alias      varchar(10)  NOT NULL,
    created_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX leagues_name_alias_unique_idx ON leagues (name, alias);

CREATE TRIGGER update_leagues_updated_at
    BEFORE UPDATE
    ON leagues
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS leagues;
-- +goose StatementEnd
