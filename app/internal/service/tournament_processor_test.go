package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"errors"
	"strings"
	"testing"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTournamentProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPersistence := NewMockPersistenceServiceInterface(ctrl)
	mockStatsProvider := stats_provider.NewMockStatsProvider(ctrl)
	mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)
	mockGamesRepo := games_repo.NewMockGamesRepo(ctrl)

	mockAggregator := ports.NewMockAggregator(ctrl)
	processor := NewTournamentProcessor(mockStatsProvider, mockPersistence, mockPlayersRepo, mockGamesRepo, mockAggregator, tournaments.TournamentModel{Id: 12})

	assert.NotNil(t, processor)
	assert.Equal(t, uint(12), processor.tournament.Id)
	assert.Equal(t, mockPersistence, processor.persistenceService)
	assert.Equal(t, mockStatsProvider, processor.statsProvider)
	assert.Equal(t, mockPlayersRepo, processor.playersRepo)
	assert.Equal(t, mockGamesRepo, processor.gamesRepo)
}

func TestTournamentProcessor_ProcessByPeriod(t *testing.T) {
	from := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	t.Run("returns error when stats provider fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(nil, errors.New("stats provider error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 1})
		count, err := processor.ProcessByPeriod(from, to)

		require.ErrorContains(t, err, "stats provider error")
		assert.Equal(t, 0, count)
	})

	t.Run("skips existing games without enrichment and saving", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntities := []games.GameStatEntity{
			{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: from}},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(gameEntities, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(true, nil)

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 99})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("enriches and saves new game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "MIA - NYK", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(gameEntity, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(77)).Return(nil)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 77})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("returns error when game existence check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "DAL - PHX", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, errors.New("db error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 11})
		count, err := processor.ProcessByPeriod(from, to)

		require.ErrorContains(t, err, "db error")
		assert.Equal(t, 0, count)
	})

	t.Run("continues when enrich game stats fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		gomock.InOrder(
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(games.GameStatEntity{}, errors.New("enrich error")),
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil),
		)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).Times(1)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(20)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 20})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when save game fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game1, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil)
		gomock.InOrder(
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(errors.New("save error")),
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil),
		)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(21)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 21})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when ListByFullName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "PLAYERS-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{PlayerModel: players.PlayerModel{FullName: "John Doe"}},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("John Doe").Return(nil, errors.New("players repo error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(22)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 22})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when player not found and GetPlayerBio fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "123",
						PlayerModel:      players.PlayerModel{FullName: ""},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("123").Return(players.PlayerBioEntity{}, errors.New("bio error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(23)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 23})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when player not found and GetPlayerBio succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-OK", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "555",
						PlayerModel:      players.PlayerModel{FullName: "", BirthDate: time.Time{}},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("555").Return(players.PlayerBioEntity{
			FullName:  "Recovered Player",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(24)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockGames, mockAggregator, tournaments.TournamentModel{Id: 24})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

type warningsRecorder struct {
	warnings []map[string]interface{}
}

func (r *warningsRecorder) InitLogger() composite_logger.Logger { return r }
func (r *warningsRecorder) Info(string, map[string]interface{}) {}

// Считаем только предупреждение о датах: чужой warning не должен выдавать себя за него.
func (r *warningsRecorder) Warn(msg string, ctx map[string]interface{}) {
	if strings.Contains(msg, "Imported game is scheduled outside tournament dates") {
		r.warnings = append(r.warnings, ctx)
	}
}
func (r *warningsRecorder) Error(string, map[string]interface{}) {}
func (r *warningsRecorder) Fatal(string, map[string]interface{}) {}

func recordWarnings(t *testing.T) *warningsRecorder {
	recorder := &warningsRecorder{}
	composite_logger.Init(recorder)
	t.Cleanup(func() { composite_logger.Init() })
	return recorder
}

func TestTournamentProcessor_ProcessByPeriod_WarnsAboutGamesOutsideTournamentDates(t *testing.T) {
	startAt := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	insideDates := time.Date(2026, 3, 1, 19, 0, 0, 0, time.UTC)
	afterEnd := endAt.AddDate(0, 1, 15)
	beforeStart := startAt.AddDate(0, 0, -3)

	cases := []struct {
		name             string
		scheduledAt      []time.Time
		exists           bool
		saveErr          error
		expectedWarnings int
	}{
		{name: "game inside dates gives no warning", scheduledAt: []time.Time{insideDates}, expectedWarnings: 0},
		{name: "game after end gives warning", scheduledAt: []time.Time{afterEnd}, expectedWarnings: 1},
		{name: "game before start gives warning", scheduledAt: []time.Time{beforeStart}, expectedWarnings: 1},
		{name: "several games outside dates give one warning per poll", scheduledAt: []time.Time{beforeStart, insideDates, afterEnd, afterEnd}, expectedWarnings: 1},
		{name: "already imported game gives no warning", scheduledAt: []time.Time{afterEnd}, exists: true, expectedWarnings: 0},
		{name: "game that failed to save gives no warning", scheduledAt: []time.Time{afterEnd}, saveErr: errors.New("save error"), expectedWarnings: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gameEntities := make([]games.GameStatEntity, 0, len(tc.scheduledAt))
			for _, scheduledAt := range tc.scheduledAt {
				gameEntities = append(gameEntities, games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: scheduledAt}})
			}

			mockStats := stats_provider.NewMockStatsProvider(ctrl)
			mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return(gameEntities, nil)
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).DoAndReturn(func(game games.GameStatEntity) (games.GameStatEntity, error) {
				return game, nil
			}).AnyTimes()
			mockGames := games_repo.NewMockGamesRepo(ctrl)
			mockGames.EXPECT().Exists(gomock.Any()).Return(tc.exists, nil).AnyTimes()
			mockPersistence := NewMockPersistenceServiceInterface(ctrl)
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(tc.saveErr).AnyTimes()
			mockAggregator := ports.NewMockAggregator(ctrl)
			mockAggregator.EXPECT().NotifyGameImported(uint(5)).Return(nil).AnyTimes()

			recorder := recordWarnings(t)
			tournament := tournaments.TournamentModel{Id: 5, Name: "NBA 2025-2026", StartAt: &startAt, EndAt: &endAt}
			processor := NewTournamentProcessor(mockStats, mockPersistence, players_repo.NewMockPlayersRepo(ctrl), mockGames, mockAggregator, tournament)

			_, err := processor.ProcessByPeriod(startAt, endAt)

			require.NoError(t, err)
			assert.Len(t, recorder.warnings, tc.expectedWarnings)
			assert.Equal(t, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), *tournament.StartAt)
			assert.Equal(t, time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC), *tournament.EndAt)
		})
	}

	t.Run("warning names tournament, its dates and game date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: afterEnd}}

		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return([]games.GameStatEntity{game}, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game, nil)
		mockGames := games_repo.NewMockGamesRepo(ctrl)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(5)).Return(nil)

		recorder := recordWarnings(t)
		tournament := tournaments.TournamentModel{Id: 5, Name: "NBA 2025-2026", StartAt: &startAt, EndAt: &endAt}
		processor := NewTournamentProcessor(mockStats, mockPersistence, players_repo.NewMockPlayersRepo(ctrl), mockGames, mockAggregator, tournament)

		_, err := processor.ProcessByPeriod(startAt, endAt)

		require.NoError(t, err)
		require.Len(t, recorder.warnings, 1)
		assert.Equal(t, uint(5), recorder.warnings[0]["tournamentId"])
		assert.Equal(t, "NBA 2025-2026", recorder.warnings[0]["tournamentName"])
		assert.Equal(t, &startAt, recorder.warnings[0]["startAt"])
		assert.Equal(t, &endAt, recorder.warnings[0]["endAt"])
		assert.Equal(t, "LAL - BOS", recorder.warnings[0]["title"])
		assert.Equal(t, afterEnd, recorder.warnings[0]["scheduledAt"])
	})
}
