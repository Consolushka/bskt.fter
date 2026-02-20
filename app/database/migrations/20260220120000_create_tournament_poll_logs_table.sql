-- +goose Up
-- +goose StatementBegin
CREATE TABLE tournament_poll_logs (
    id BIGSERIAL PRIMARY KEY,
    tournament_id BIGINT NOT NULL,
    poll_start_at TIMESTAMPTZ NOT NULL,
    poll_end_at TIMESTAMPTZ,
    interval_start TIMESTAMPTZ NOT NULL,
    interval_end TIMESTAMPTZ NOT NULL,
    saved_games_count INT DEFAULT 0,
    status VARCHAR(50) NOT NULL, -- 'success', 'error'
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_poll_logs_tournament FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE
);

-- Индекс для быстрого поиска последней успешной отметки для турнира
CREATE INDEX idx_poll_logs_tournament_status_end ON tournament_poll_logs (tournament_id, status, interval_end DESC);

-- Удаляем старую таблицу, так как она больше не нужна
DROP TABLE IF EXISTS poll_watermarks;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tournament_poll_logs;

CREATE TABLE poll_watermarks (
    tournament_id BIGINT PRIMARY KEY,
    last_successful_poll_at timestamptz NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_poll_watermarks_tournament FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE
);
-- +goose StatementEnd
