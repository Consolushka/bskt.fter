-- +goose Up
-- +goose StatementBegin
ALTER TABLE tournaments ADD COLUMN deleted_at TIMESTAMP NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournaments DROP COLUMN deleted_at;
-- +goose StatementEnd
