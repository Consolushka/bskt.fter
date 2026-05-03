-- +goose Up
-- +goose StatementBegin
ALTER TABLE tournament_poll_logs ADD COLUMN next_poll_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournament_poll_logs DROP COLUMN next_poll_at;
-- +goose StatementEnd
