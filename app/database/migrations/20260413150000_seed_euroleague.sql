-- +goose Up
-- +goose StatementBegin
-- Синхронизируем последовательности на случай ручных вставок ID
SELECT setval('leagues_id_seq', COALESCE((SELECT MAX(id) FROM leagues), 1), (SELECT MAX(id) FROM leagues) IS NOT NULL);
SELECT setval('tournaments_id_seq', COALESCE((SELECT MAX(id) FROM tournaments), 1), (SELECT MAX(id) FROM tournaments) IS NOT NULL);

INSERT INTO leagues (name, alias)
VALUES ('Euroleague', 'euroleague')
ON CONFLICT (name, alias) DO NOTHING;

INSERT INTO tournaments (league_id, name, start_at, end_at, regulation_duration)
SELECT id, 'Euroleague 2025-2026', '2025-10-01 00:00:00', '2026-05-31 23:59:59', 40
FROM leagues
WHERE alias = 'euroleague'
LIMIT 1
ON CONFLICT (league_id, name) DO NOTHING;

INSERT INTO tournament_providers (tournament_id, provider_name, external_id, params)
SELECT id, 'API_BASKETBALL', '120', '{"season": "2025"}'::jsonb
FROM tournaments
WHERE name = 'Euroleague 2025-2026'
LIMIT 1
ON CONFLICT (tournament_id, provider_name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tournament_providers WHERE provider_name = 'API_BASKETBALL' AND external_id = '120';
DELETE FROM tournaments WHERE name = 'Euroleague 2025-2026';
-- +goose StatementEnd
