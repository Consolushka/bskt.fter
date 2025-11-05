-- +goose Up
-- +goose StatementBegin
-- Переименование таблицы
ALTER TABLE tournament_external_ids RENAME TO tournament_providers;
-- Изменение колонки external_id на nullable
ALTER TABLE tournament_providers ALTER COLUMN external_id DROP NOT NULL;
-- Добавление колонки params
ALTER TABLE tournament_providers ADD COLUMN params JSONB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournament_providers DROP COLUMN params;
ALTER TABLE tournament_providers ALTER COLUMN external_id SET NOT NULL;
ALTER TABLE tournament_providers RENAME TO tournament_external_ids;
-- +goose StatementEnd
