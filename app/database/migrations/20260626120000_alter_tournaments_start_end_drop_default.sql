-- +goose Up
-- +goose StatementBegin
-- Неизвестные даты турнира храним как NULL, а не как now():
-- значение по умолчанию now() делало только что добавленный турнир "уже завершённым"
-- (его end_at оказывался в прошлом сразу после вставки).
ALTER TABLE tournaments ALTER COLUMN start_at DROP DEFAULT;
ALTER TABLE tournaments ALTER COLUMN end_at DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournaments ALTER COLUMN start_at SET DEFAULT now();
ALTER TABLE tournaments ALTER COLUMN end_at SET DEFAULT now();
-- +goose StatementEnd
