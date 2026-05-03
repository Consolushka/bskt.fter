-- +goose Up
-- +goose StatementBegin
ALTER TABLE leagues ADD COLUMN tier SMALLINT;
ALTER TABLE leagues ADD COLUMN country VARCHAR(255);
ALTER TABLE tournaments ADD COLUMN tier SMALLINT;
ALTER TABLE teams ADD COLUMN alias VARCHAR(10);
UPDATE teams SET alias = UPPER(SUBSTRING(name, 1, 3));
ALTER TABLE teams ALTER COLUMN alias SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE teams DROP COLUMN alias;
ALTER TABLE tournaments DROP COLUMN tier;
ALTER TABLE leagues DROP COLUMN country;
ALTER TABLE leagues DROP COLUMN tier;
-- +goose StatementEnd
