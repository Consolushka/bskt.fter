-- +goose Up
-- +goose StatementBegin
-- Seed scheduled tasks with dynamic dates (Moscow timezone)
INSERT INTO scheduled_tasks (type, next_execution_at) VALUES
      -- American tournaments: today at 10:59 MSK (or tomorrow if already passed)
      ('process_american_tournaments_task',
       CASE
           WHEN CURRENT_TIME > TIME '10:59:00'
               THEN (CURRENT_DATE + INTERVAL '1 day' + TIME '10:59:00')
           ELSE (CURRENT_DATE + TIME '10:59:00')
           END
      ),
      -- European not urgent: today at 04:59 MSK (or tomorrow if already passed)
      ('process_not_urgent_european_tournaments_task',
       CASE
           WHEN CURRENT_TIME > TIME '04:59:00'
               THEN (CURRENT_DATE + INTERVAL '1 day' + TIME '04:59:00')
           ELSE (CURRENT_DATE + TIME '04:59:00')
           END
      ),
      -- European urgent: today at 23:59 MSK (or tomorrow if already passed)
      ('process_urgent_european_tournaments_task',
       CASE
           WHEN CURRENT_TIME > TIME '23:59:00'
               THEN (CURRENT_DATE + INTERVAL '1 day' + TIME '23:59:00')
           ELSE (CURRENT_DATE + TIME '23:59:00')
           END
      );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM scheduled_tasks WHERE type IN (
       'process_american_tournaments_task',
       'process_not_urgent_european_tournaments_task',
       'process_urgent_european_tournaments_task'
    );
-- +goose StatementEnd