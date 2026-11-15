-- +goose Up
-- +goose StatementBegin
CREATE TABLE team_external_ids
(
    id          serial PRIMARY KEY,
    scope       varchar(255) NOT NULL,
    external_id integer      NOT NULL,
    team_id     int          NOT NULL REFERENCES teams (id) ON DELETE CASCADE,
    created_at  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX team_external_ids_scope_external_id_unique_idx ON team_external_ids (scope, external_id);

CREATE TRIGGER update_team_external_ids_updated_at
    BEFORE UPDATE
    ON team_external_ids
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_external_ids;
-- +goose StatementEnd
