-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournament_external_ids
(
    tournament_id int          NOT NULL REFERENCES tournaments (id) ON DELETE CASCADE,
    provider_name varchar(100) NOT NULL,
    external_id   varchar(255) NOT NULL,
    created_at    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (tournament_id, provider_name)
);

CREATE UNIQUE INDEX tournament_external_ids_tournament_provider_unique_idx
    ON tournament_external_ids (tournament_id, provider_name);

CREATE INDEX tournament_external_ids_provider_external_id_idx
    ON tournament_external_ids (provider_name, external_id);

CREATE TRIGGER update_tournament_external_ids_updated_at
    BEFORE UPDATE
    ON tournament_external_ids
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tournament_external_ids;
-- +goose StatementEnd