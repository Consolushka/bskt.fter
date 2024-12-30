-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE leagues
(
    id                serial PRIMARY KEY,
    name_local        varchar(255) NOT NULL,
    alias_local       varchar(10)  NOT NULL,
    name_en           varchar(255) NOT NULL,
    alias_en          varchar(10)  NOT NULL,
    periods_number    int          NOT NULL,
    period_duration   int          NOT NULL,
    overtime_duration int          NOT NULL,
    created_at        timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX leagues_alias_local_idx ON leagues (alias_local);
CREATE UNIQUE INDEX leagues_alias_en_idx ON leagues (alias_en);

ALTER TABLE leagues
    ADD CONSTRAINT check_periods_number CHECK (periods_number > 0);
ALTER TABLE leagues
    ADD CONSTRAINT check_period_duration CHECK (period_duration > 0);

CREATE TRIGGER update_leagues_updated_at
    BEFORE UPDATE
    ON leagues
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS leagues;
DROP FUNCTION IF EXISTS trigger_set_timestamp;
-- +goose StatementEnd
